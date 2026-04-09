package main

import (
	"fmt"
	"sync"
)

type Job func() (int, int, error)

type Result struct {
	JobID    int
	WorkerID int
	Value    int
	Err      error
}

func worker(WorkerID int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		jobID, result, err := job()
		results <- Result{JobID: jobID, WorkerID: WorkerID, Value: result, Err: err}
	}
}

func main() {
	var wg sync.WaitGroup

	jobs := make([]Job, 0, 10)
	for i := 0; i < 10; i++ {
		j := i
		jobs = append(jobs, func() (int, int, error) {
			return j, j * 2, nil
		})
	}

	jobsChan := make(chan Job)
	resultsChan := make(chan Result)

	numWorkers := 3

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobsChan, resultsChan, &wg)
	}

	go func() {
		for _, job := range jobs {
			jobsChan <- job
		}
		close(jobsChan)
	}()

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for result := range resultsChan {
		fmt.Printf("Job ID: %d, Worker ID: %d, Result: %d, Error: %v\n", result.JobID, result.WorkerID, result.Value, result.Err)
	}
}
