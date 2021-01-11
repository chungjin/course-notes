package mapreduce

import (
	"fmt"
	"sync"
)

//
// schedule() starts and waits for all tasks in the given phase (mapPhase
// or reducePhase). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	var ntasks int
	var n_other int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mapFiles)
		n_other = nReduce
	case reducePhase:
		ntasks = nReduce
		n_other = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, n_other)

	// All ntasks tasks have to be scheduled on workers. Once all tasks
	// have completed successfully, schedule() should return.
	//
	// Your code here (Part III, Part IV).
	//

	var wg sync.WaitGroup
	for i := 0; i < ntasks; i++ { //how to create a pool for goroutine
		debug("wait %d:%d %v task find the available channel\n", i, ntasks, phase)
		wg.Add(1)
		go assignTask(jobName, mapFiles[i], phase, i, n_other, &wg, registerChan)
	}
	wg.Wait() //wait for all go routine done
}

func assignTask(jobName string, file string, phase jobPhase, taskNum int, n_other int, wg *sync.WaitGroup, registerChan chan string) {
	for {
		args := &DoTaskArgs{
			JobName:       jobName,
			File:          file,
			Phase:         phase,
			TaskNumber:    taskNum,
			NumOtherPhase: n_other,
		}

		newWorker := <-registerChan
		ok := call(newWorker, "Worker.DoTask", args, nil)
		if ok {
			wg.Done()
			registerChan <- newWorker
			break
		} else {
			//registerChan <- newWorker
			fmt.Printf("worker %s not reply", newWorker)
		}
	}
}
