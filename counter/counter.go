package counter

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// register function to Counter
//
// IntermediateChan: make sure no data sent to closed channel
// JobsChan: actual processing data
type Counter struct {
	IntermediateChan chan Data
	JobsChan         chan Data
}

// Data: holds index & duration
type Data struct {
	Index    int
	Duration int
}

func NewCounter() Counter {
	return Counter{
		IntermediateChan: make(chan Data),
		JobsChan:         make(chan Data),
	}
}

// Start:
// start adding process to counter
func (c Counter) Start(simulationTime []int) {
	eventIndex := 0
	for i := range simulationTime {
		c.CallbackFunc(eventIndex, simulationTime[i])
		eventIndex++
		time.Sleep(time.Millisecond * 100)
	}
	fmt.Println("Done processing...")
}

// CallbackFunc: register each event to intermediate channel
func (c Counter) CallbackFunc(event, duration int) {
	c.IntermediateChan <- Data{
		Index:    event,
		Duration: duration,
	}
}

// StartProcess:
// loop for every data in intermediateChan (move to jobsChan) or
// until ctx.Done received
func (c Counter) StartProcess(ctx context.Context) {
	for {
		select {
		case job := <-c.IntermediateChan:
			c.JobsChan <- job
		case <-ctx.Done():
			fmt.Println("Counter received cancellation signal, closing jobsChan!")
			close(c.JobsChan)
			fmt.Println("Counter closed")
			return // exit this function so we don't consume anything more from from the intermittentChan
		}
	}
}

// WorkerJobs:
// starts a worker function that will range on the jobsChan until that channel closes
func (c Counter) WorkerJobs(wg *sync.WaitGroup, index int) {
	defer wg.Done()

	fmt.Printf("Worker %d starting\n", index)
	for each := range c.JobsChan {
		// simulate work taking between 1-3 seconds
		fmt.Printf("Worker %d started job %d\n", index, each.Index)
		sleepFor := time.Second * time.Duration(each.Duration)
		time.Sleep(sleepFor)
		fmt.Printf("Worker %d finished processing job %d in %ds\n", index, each.Index, sleepFor)
	}
	fmt.Printf("Worker %d interrupted\n", index)
}
