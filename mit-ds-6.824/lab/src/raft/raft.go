package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"bytes"
	"labgob"
	"labrpc"
	"math/rand"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// import "bytes"
// import "labgob"

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in Lab 3 you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh; at that point you can add fields to
// ApplyMsg, but set CommandValid to false for these other uses.
//

type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
	UseSnapshot bool
}

type LogRecord struct {
	Term    int
	Command interface{}
}

type Role string

const (
	Leader    Role = "Leader"
	Candidate Role = "Candidate"
	Follower  Role = "Follower"
)

const (
	HeatbeartInterval = time.Duration(100) * time.Millisecond
	ElectionTimeout   = time.Duration(300) * time.Millisecond
)

//The state for the machine, all fileds are public, for snapshot and restore
type SnapshotState struct {
	VotedFor    int
	Term        int
	Log         []LogRecord
	CommitIndex int
	LastApplied int
}

func init(){
	//log.SetLevel(log.DebugLevel)
	log.SetLevel(log.InfoLevel)
}

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	ticker *time.Ticker

	//persistent state on all servers
	role        Role
	applyCh     chan ApplyMsg //to reply to client
	currentTerm int
	votedFor    int         //candidatedId that received vote in current Term
	log         []LogRecord //log entries, from 1

	//Volatile state on all servers
	commitIndex int //index of highest log entry known to be committed, from 0
	lastApplied int //index of highest log entry applied to state machine

	//Volatile state on leaders
	nextIndex        []int //for each server, index of the next log entry send to that server
	matchIndex       []int //for each server, index of highest log entry known to be replicated on server
	AppendEntryArgs  chan AppendEntriesArgs
	AppendEntryReply chan AppendEntriesReply
	RequestVoteArgs  chan RequestVoteArgs
	RequestVoteReply chan *RequestVoteReply
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (2A).
	rf.mu.Lock()
	defer rf.mu.Unlock()

	term = rf.currentTerm
	isleader = rf.role == Leader
	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	//todo(2c) persist the state? when to call?
	var buf bytes.Buffer
	e := labgob.NewEncoder(&buf)
	snapshot := &SnapshotState{
		VotedFor:    rf.votedFor,
		Term:        rf.currentTerm,
		Log:         rf.log,
		CommitIndex: rf.commitIndex,
		LastApplied: rf.lastApplied,
	}
	e.Encode(snapshot)
	data := buf.Bytes()
	rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }

	r := bytes.NewBuffer(data)
	d := labgob.NewDecoder(r)
	var snapshot SnapshotState
	if d.Decode(&snapshot) != nil {
		//todo(2c): how to deal with error
		log.Fatalf("unable to decode snapshot")
	} else {
		rf.votedFor = snapshot.VotedFor
		rf.currentTerm = snapshot.Term
		rf.log = snapshot.Log
		rf.commitIndex = snapshot.CommitIndex
		rf.lastApplied = snapshot.LastApplied
	}
}

//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Term        int  //always piggyback term id
	VoteGranted bool //true indicate approve, false means not
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	rf.RequestVoteArgs <- *args
	res := <-rf.RequestVoteReply
	reply.Term = res.Term
	reply.VoteGranted = res.VoteGranted
}

func (rf *Raft) RequestVoteHandler(args *RequestVoteArgs) *RequestVoteReply {
	logger := log.WithFields(log.Fields{
		"id": rf.me,
		"term": rf.currentTerm,
		"func": "requestVote",
	})

	logger.Debugf("term %d %d receive requestVote from %d, arg %+v, last log %+v, index %d votedfor %d", rf.currentTerm, rf.me, args.CandidateId, args, rf.log[len(rf.log)-1], len(rf.log)-1, rf.votedFor)

	reply := &RequestVoteReply{}
	reply.Term = rf.currentTerm
	reply.VoteGranted = false
	// Your code here (2A, 2B).

	defer func() {
		if reply.VoteGranted {
			logger.Debugf("%d grantVote to %d", rf.me, args.CandidateId)
		}
	}()

	if args.Term < rf.currentTerm {
		return reply
	}

	//cannot grant serveral
	if args.Term == rf.currentTerm {
		if !(rf.votedFor == -1 || rf.votedFor == args.CandidateId) {
			return reply
		}
	}

	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.role = Follower
	}

	if args.LastLogTerm > rf.log[len(rf.log)-1].Term {
		reply.VoteGranted = true
		rf.votedFor = args.CandidateId
		rf.persist()
		return reply
	}

	if args.LastLogTerm == rf.log[len(rf.log)-1].Term {
		if args.LastLogIndex >= len(rf.log)-1 {
			reply.VoteGranted = true
			rf.votedFor = args.CandidateId
			rf.persist()
			return reply
		}
	}

	reply.VoteGranted = false
	return reply
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
	return ok
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	// Your code here (2B).
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if rf.role != Leader {
		return -1, -1, false
	}

	record := LogRecord{
		Term:    rf.currentTerm,
		Command: command,
	}
	rf.log = append(rf.log, record)
	rf.persist()
	log.WithFields(log.Fields{
		"func": "Start",
		"id": rf.me,
	}).Debugf("receive %+v", record)

	return len(rf.log) - 1, rf.currentTerm, true
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.log = []LogRecord{{Term: 0}} //start from 1
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).
	rf.commitIndex = 0 //the log index to commit
	rf.lastApplied = 0 //the log index have applied to state machine
	rf.votedFor = -1
	rf.ticker = getRamdomTickerTime() //timeout timer
	rf.role = Follower
	rf.applyCh = applyCh
	rf.AppendEntryArgs = make(chan AppendEntriesArgs, 10)   //the channel to receive appendEntry
	rf.AppendEntryReply = make(chan AppendEntriesReply, 10) //the channel to receive appendEntry
	rf.RequestVoteArgs = make(chan RequestVoteArgs, 10)     //the channel to receive requestVote
	rf.RequestVoteReply = make(chan *RequestVoteReply)

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())
	log.WithFields(log.Fields{
		"id": rf.me,
		"term": rf.currentTerm,
		"func": "Make",
	}).Debugf("restart")

	go rf.acceptRPC()

	return rf
}

func (rf *Raft) acceptRPC() {
	for {
		log.WithFields(log.Fields{
			"id": rf.me,
			"term": rf.currentTerm,
			"func": "acceptRPC",
		}).Debugf("change role %s", rf.role)
		switch rf.role {
		case Leader:
			rf.leaderRoutine()
		case Candidate:
			rf.candidateRoutine()
		case Follower:
			rf.followerRoutine()
		}
	}

}

func (rf *Raft) leaderRoutine() {
	//init volatile

	logger := log.WithFields(log.Fields{
		"id": rf.me,
		"term": rf.currentTerm,
		"func": "leaderRoutine",
	})
	rf.mu.Lock()
	rf.nextIndex = make([]int, len(rf.peers))
	for i := range rf.nextIndex {
		rf.nextIndex[i] = len(rf.log) //initial to leader last log index + 1
	}
	rf.matchIndex = make([]int, len(rf.peers)) //initial to zero
	rf.mu.Unlock()

	heartBeatTicker := time.NewTicker(HeatbeartInterval)
	appendEntriesReplyCh := make(chan AppendEntriesReply, 10) //deal with rpc reply

leader:
	for {
		select {
		//receive request vote
		case args := <-rf.RequestVoteArgs:
			rf.mu.Lock()
			rf.RequestVoteReply <- rf.RequestVoteHandler(&args)
			if rf.role == Follower {
				rf.mu.Unlock()
				break leader
			}
			rf.mu.Unlock()

		//receive AppendEntryCh, means discover leader with higher term, become follower
		case request := <-rf.AppendEntryArgs:
			rf.mu.Lock()
			//debug(fmt.Sprintf("leader %d receive appendEntry, currentTerm %d, request Term %d", rf.me, rf.currentTerm, request.Term))
			rf.AppendEntryReply <- rf.AppendEntriesHandler(request)
			if rf.role == Follower {
				rf.mu.Unlock()
				break leader
			}
			rf.mu.Unlock()

		//receive appendEntriesReply
		//if successful, update nextIndex and matchIndex accordingly
		case reply := <-appendEntriesReplyCh:
			rf.mu.Lock()
			if reply.Term > rf.currentTerm {
				rf.role = Follower
				rf.mu.Unlock()
				break leader
			}
			if !reply.IsHeartBeat {
				logger.Debugf("receive response from %d reply, with %+v", reply.ServerId, reply)
				if reply.Success {
					//update the next index and match index
					rf.nextIndex[reply.ServerId] = reply.NextInt
					rf.matchIndex[reply.ServerId] = reply.NextInt - 1
					if rf.updateCommitIndex() {
						rf.applyToStateMachine()
					}
				} else {
					//FIXME: optimize to reduce the number of reject
					if rf.nextIndex[reply.ServerId] > 1 {
						rf.nextIndex[reply.ServerId]--
					}
				}

			} else if reply.NextInt != 0 {
				rf.nextIndex[reply.ServerId] = reply.NextInt
			}
			rf.mu.Unlock()

		//send heartbeat to each peer
		case <-heartBeatTicker.C:
			rf.mu.Lock()
			logger.Debugf("leader last log %+v, index %d", rf.log[len(rf.log)-1], len(rf.log)-1)
			args := AppendEntriesArgs{
				Term:         rf.currentTerm,
				LeaderId:     rf.me,
				LeaderCommit: rf.commitIndex,
			}
			for i := range rf.peers {
				if i == rf.me {
					continue
				}
				args.Entries = make([]LogRecord, 0)
				args.PrevLogIndex = rf.nextIndex[i] - 1
				logger.Debugf("peer %d index %d len %d", i, args.PrevLogIndex, len(rf.log))
				args.PrevLogTerm = rf.log[args.PrevLogIndex].Term
				if rf.nextIndex[i] < len(rf.log) { //append log
					args.Entries = make([]LogRecord, len(rf.log[args.PrevLogIndex+1:]))
					copy(args.Entries, rf.log[args.PrevLogIndex+1:])
				}
				go func(index int, args AppendEntriesArgs) {
					reply := AppendEntriesReply{}
					if ok := rf.sendAppendEntries(index, &args, &reply); ok {
						appendEntriesReplyCh <- reply
					}
				}(i, args)
			}
			rf.mu.Unlock()
		default:
			time.Sleep(HeatbeartInterval)
		}
	}
}

func (rf *Raft) candidateRoutine() {
	rf.mu.Lock()
	rf.currentTerm++
	rf.votedFor = rf.me
	rf.ticker = getRamdomTickerTime()
	//rf.ticker = time.NewTicker(ElectionTimeout)
	count := 1
	replyCh := make(chan RequestVoteReply, 10)
	success := make(chan bool, 10)
	//broadcast request vote
	args := RequestVoteArgs{
		Term:         rf.currentTerm,
		CandidateId:  rf.me,
		LastLogIndex: len(rf.log) - 1,
		LastLogTerm:  rf.log[len(rf.log)-1].Term,
	}
	rf.mu.Unlock()
	for i := range rf.peers {
		if i == rf.me {
			continue
		}
		go func(index int) {
			reply := RequestVoteReply{}
			if ok := rf.sendRequestVote(index, &args, &reply); ok {
				//process the reply
				//debug(fmt.Sprintf("reply %v", reply))
				replyCh <- reply
			}
		}(i)
	}

candidate:
	for {
		select {

		//receive appendEntry from new leader, becomes follower
		case args := <-rf.AppendEntryArgs:
			//debug(fmt.Sprintf("%d receive appendEntry", rf.me))
			rf.AppendEntryReply <- rf.AppendEntriesHandler(args)
			rf.mu.Lock()
			if args.Term >= rf.currentTerm {
				rf.role = Follower
				rf.mu.Unlock()
				break candidate
			}
			rf.mu.Unlock()
		//get the majority vote, becomes leader
		case <-success:
			rf.mu.Lock()
			rf.role = Leader
			rf.mu.Unlock()
			break candidate
		//receive voteRequest reply
		case re := <-replyCh:
			rf.mu.Lock()
			if re.Term > rf.currentTerm {

				rf.role = Follower
				rf.mu.Unlock()
				break candidate
			}
			if re.VoteGranted {
				count++
			}
			if count > len(rf.peers)/2 { //get the majority
				success <- true
			}
			rf.mu.Unlock()

		//receive request vote
		case args := <-rf.RequestVoteArgs:
			rf.mu.Lock()
			res := rf.RequestVoteHandler(&args)
			rf.RequestVoteReply <- res

			if rf.role == Follower {
				rf.mu.Unlock()
				break candidate
			}
			rf.mu.Unlock()
		//timeout fire, new election
		case <-rf.ticker.C:
			break candidate
		}

	}

}

func (rf *Raft) followerRoutine() {
follower:
	for {

		select {
		case args := <-rf.AppendEntryArgs: //response to applyentriesRPC
			//reset timer
			rf.mu.Lock()
			rf.AppendEntryReply <- rf.AppendEntriesHandler(args)
			rf.ticker = getRamdomTickerTime()
			rf.mu.Unlock()

		//reset timer
		case args := <-rf.RequestVoteArgs:
			rf.mu.Lock()
			rf.RequestVoteReply <- rf.RequestVoteHandler(&args)
			rf.ticker = getRamdomTickerTime()
			rf.mu.Unlock()

		//timeout fire, compete for leader
		case <-rf.ticker.C:
			//(fmt.Sprintf("%d follower, ticker fire, compete for election", rf.me))
			rf.mu.Lock()
			rf.role = Candidate
			rf.mu.Unlock()
			break follower //jump out infinite for loop, change to another role, break only break from inner most "for", "switch", "select"
		}

	}
}

//
// example AppendEntries RPC arguments structure.
// field names must start with capital letters!
//
type AppendEntriesArgs struct {
	// Your data here (2A).
	Term         int
	LeaderId     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogRecord
	LeaderCommit int //leader's commitIndex
}

//
// example AppendEntries RPC reply structure.
// field names must start with capital letters!
//
type AppendEntriesReply struct {
	// Your data here (2A).
	NextInt     int
	IsHeartBeat bool
	ServerId    int
	Term        int
	Success     bool
}

// Will update the current term to the most updated
// example AppendEntries RPC handler. Will receive heartbeat or command
//
func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	// Your code here (2A).
	rf.AppendEntryArgs <- *args
	res := <-rf.AppendEntryReply
	reply.Term = res.Term
	reply.ServerId = res.ServerId
	reply.Success = res.Success
	reply.IsHeartBeat = res.IsHeartBeat
	reply.NextInt = res.NextInt
}

func (rf *Raft) AppendEntriesHandler(args AppendEntriesArgs) AppendEntriesReply {
	//if len(args.Entries) != 0 {
	logger := log.WithFields(log.Fields{
		"id": rf.me,
		"term": rf.currentTerm,
		"func": "AppendEntriesHandler",
	})
	logger.Debugf("receive append arg %+v, last log %+v, index %d", args, rf.log[len(rf.log)-1], len(rf.log)-1)
	//}

	reply := AppendEntriesReply{
		IsHeartBeat: len(args.Entries) == 0,
		NextInt: 0,
		Success: false,
		ServerId: rf.me,
		Term: rf.currentTerm,
	}

	if args.Term < rf.currentTerm {
		return reply
	}

	if args.Term >= rf.currentTerm {
		rf.role = Follower
		rf.currentTerm = args.Term
	}

	if len(rf.log) <= args.PrevLogIndex {
		reply.NextInt = len(rf.log)
		return reply
	}

	if rf.log[args.PrevLogIndex].Term != args.PrevLogTerm {
		return reply
	}

	if args.LeaderCommit < rf.commitIndex {
		return reply
	}

	if len(rf.log) == args.PrevLogIndex+1 {
		rf.log = append(rf.log, args.Entries...)
	} else {
		index := args.PrevLogIndex + 1
		for _, entry := range args.Entries {
			if len(rf.log) == index { //apply directly
				rf.log = append(rf.log, entry)
			} else {
				original := rf.log[index]
				if original.Term != entry.Term {
					//delete all following
					rf.log = rf.log[:index]
					rf.log = append(rf.log, entry)
				}
			}
			index += 1
		}
	}

	if args.LeaderCommit > rf.commitIndex {
		//debug(fmt.Sprintf("%d role %s update commit from %d to %d", rf.me, rf.role, rf.commitIndex, args.LeaderCommit))
		if args.LeaderCommit > len(rf.log)-1 {
			rf.commitIndex = len(rf.log) - 1
		} else {
			rf.commitIndex = args.LeaderCommit
		}
		rf.applyToStateMachine()
	}

	if len(args.Entries) != 0 {
		//reply.NextIndex = args.Entries[len(args.Entries)-1].Index + 1
		reply.Success = true
		reply.NextInt = len(rf.log)
	}
	log.Debugf("Receiver last log %+v, index %d", rf.log[len(rf.log)-1], len(rf.log)-1)

	return reply
}

// execute to state machine, if is leader, reply to client
func (rf *Raft) applyToStateMachine() {
	var stateChanged bool = false
	for rf.lastApplied < rf.commitIndex {
		stateChanged = true
		rf.lastApplied++
		// reply to service
		msg := ApplyMsg{
			CommandValid: true,
			Command:      rf.log[rf.lastApplied].Command,
			CommandIndex: rf.lastApplied,
		}
		log.WithFields(log.Fields{
			"id": rf.me,
			"term": rf.currentTerm,
			"func": "applyToStateMachine",
		}).Debugf("role %s, apply %+v, commit index %d, last log %+v, index %d", rf.role, msg, rf.commitIndex, rf.log[len(rf.log)-1], len(rf.log)-1)
		rf.applyCh <- msg
	}
	if stateChanged {
		rf.persist()
	}
}

//return ticker with 150~300 millisecond interval
func getRamdomTickerTime() *time.Ticker {
	random := time.Duration(rand.Int31n(150) + 150)
	return time.NewTicker(random * time.Millisecond)
}

func (rf *Raft) updateCommitIndex() bool {
	major := len(rf.peers) / 2
	for i := len(rf.log) - 1; i > rf.commitIndex; i-- {
		if rf.log[i].Term != rf.currentTerm {
			break //commit the log in current term
		}
		count := 1
		for peer := range rf.peers {
			if peer == rf.me {
				continue
			}
			if rf.matchIndex[peer] >= i {
				count++
			}
		}
		if count > major {
			log.WithFields(log.Fields{
				"id": rf.me,
				"role": rf.role,
				"term": rf.currentTerm,
				"func": "updateCommitIndex",
			}).Debugf("update commit, commit %d, matchIndex %+v", i, rf.matchIndex)
			rf.commitIndex = i
			return true
		}
	}
	return false
}