package raftkv

import (
	"labgob"
	"labrpc"
	"log"
	"raft"
	"sync"
	"time"
)

type Command string

const(
	PUT Command = "Put"
	GET Command = "Get"
	APPEND Command = "Append"
)

const Debug = 0

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

type Op struct {
	// Your definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.
	Command   Command // "get" | "put" | "append"
	ClientId  int
	RequestId int
	Key       string
	Value     string
}

type Result struct{
	Command     Command
	OK          bool
	ClientId    int
	RequestId   int
	WrongLeader bool
	Err         Err
	Key         string
	Value       string
}

type KVServer struct {
	mu      sync.Mutex
	me      int
	rf      *raft.Raft
	applyCh chan raft.ApplyMsg

	maxraftstate int // snapshot if log grows this big

	// Your definitions here.
	data     map[string]string   // key-value data
	ack      map[int]int     // client's latest request id (for deduplication)
	resultCh map[int]chan Result // log index to result of applying that entry
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	// Your code here.
	entry := Op{
		Command:   GET,
		ClientId:  args.ClientId,
		RequestId: args.RequestId,
		Key:       args.Key,
	}

	result := kv.appendEntryToLog(entry)
	if !result.OK{
		reply.WrongLeader = true
		return
	}
	reply.WrongLeader = false
	reply.Err = result.Err
	reply.Value = result.Value
}

// PutAppend call through RPC by Client
func (kv *KVServer) PutAppend(args *PutAppendArgs, reply *PutAppendReply) {

	defer func() {
		DPrintf("PutAppend, server %d, received args %+v, reply %+v", kv.me, args, reply)
	}()
	// Your code here.
	entry := Op{
		ClientId:  args.ClientId,
		RequestId: args.RequestId,
		Key:       args.Key,
		Value:     args.Value,
		Command:   args.Op,
	}

	result := kv.appendEntryToLog(entry)
	if !result.OK {
		reply.WrongLeader = true
		return
	}

	reply.WrongLeader = false
	reply.Err = result.Err

}

//
// the tester calls Kill() when a KVServer instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (kv *KVServer) Kill() {
	kv.rf.Kill()
	// Your code here, if desired.
}

//
// servers[] contains the ports of the set of
// servers that will cooperate via Raft to
// form the fault-tolerant key/value service.
// me is the index of the current server in servers[].
// the k/v server should store snapshots through the underlying Raft
// implementation, which should call persister.SaveStateAndSnapshot() to
// atomically save the Raft state along with the snapshot.
// the k/v server should snapshot when Raft's saved state exceeds maxraftstate bytes,
// in order to allow Raft to garbage-collect its log. if maxraftstate is -1,
// you don't need to snapshot.
// StartKVServer() must return quickly, so it should start goroutines
// for any long-running work.
//
func StartKVServer(servers []*labrpc.ClientEnd, me int, persister *raft.Persister, maxraftstate int) *KVServer {
	// call labgob.Register on structures you want
	// Go's RPC library to marshall/unmarshall.
	labgob.Register(Op{})

	kv := new(KVServer)
	kv.me = me
	kv.maxraftstate = maxraftstate

	// You may need initialization code here.

	kv.applyCh = make(chan raft.ApplyMsg)
	kv.rf = raft.Make(servers, me, persister, kv.applyCh)

	// You may need initialization code here.
	kv.data = make(map[string]string)
	kv.ack = make(map[int]int)
	kv.resultCh = make(map[int]chan Result) // commit index to Result

	go kv.Run()
	return kv
}

// appendEntryToLog is called when received client rpc,
// start raft consensus algorithm
// and check the reply through resultCh(update by another go routine)
func (kv *KVServer) appendEntryToLog(entry Op) Result{
	// begin raft consensus algorithm
	index, _, isLeader := kv.rf.Start(entry)

	if !isLeader || index<1{
		return Result{WrongLeader: true, OK: false}
	}
	kv.mu.Lock()
	DPrintf("appendEntryToLog: server %d, index %d, op %+v", kv.me, index, entry)
	if _, ok := kv.resultCh[index]; !ok{
		kv.resultCh[index] = make(chan Result, 1)
	}
	kv.mu.Unlock()
	select {
	case result := <-kv.resultCh[index]:
		DPrintf("poll: server %d, index %d, result %+v", kv.me, index, result)
		if isMatch(entry, result) {
			DPrintf("server %d match", kv.me)
			return result
		}
		return Result{OK: false}
	case <-time.After(240 * time.Millisecond):
		return Result{OK: false}
	}


}


func isMatch(entry Op, result Result) bool{
	// return by the same
	return entry.ClientId == result.ClientId && entry.RequestId == result.RequestId
}


func (kv *KVServer) Run() {
	// keeps fetching the response, that raft peer return
	for{
		// read back the msg from the receive channel
		msg := <-kv.applyCh
		//DPrintf("Run, server %d, received %+v",kv.me, msg)
		kv.mu.Lock()
		if msg.UseSnapshot{
			//todo: implement it
		}else{
			op := msg.Command.(Op)
			result := kv.applyOp(op)
			_, ok := kv.resultCh[msg.CommandIndex]
			if !ok{
				kv.resultCh[msg.CommandIndex] = make(chan Result, 1)
			}
			kv.resultCh[msg.CommandIndex]<-result
			//todo: compress log

		}
		kv.mu.Unlock()
	}


}

// applyOp when receive raft peer response, try apply operation to db and return result
func (kv *KVServer) applyOp(op Op) Result{
	result := Result{
		Command:     op.Command,
		OK:          true,
		WrongLeader: false,
		ClientId:    op.ClientId,
		RequestId:   op.RequestId,
		Key:         op.Key,
		Value:       op.Value,
	}

	switch op.Command {
	case PUT:
		if !kv.isDuplicated(op){
			kv.data[op.Key] = op.Value
		}
		DPrintf("server %d, duplicate %v, db %+v", kv.me, kv.isDuplicated(op), kv.data)
	case APPEND:
		if !kv.isDuplicated(op){
			kv.data[op.Key] += op.Value
		}
		DPrintf("server %d, db %+v", kv.me, kv.data)
	case GET:
		if value, ok := kv.data[op.Key]; ok{
			result.Value = value
		}else{
			result.Err = ErrNoKey
		}

	}
	kv.ack[op.ClientId] = op.RequestId
	return result
}

func (kv *KVServer) isDuplicated(op Op) bool{
	if lastSeqId, ok := kv.ack[op.ClientId]; ok{
		return lastSeqId >= op.RequestId
	}
	return false
}