// Copyright 2020 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"context"
	"sync"

	"github.com/apex/log"
)

// Task is a generic interface for a runnable operation
type Task interface {
	Run(ctx context.Context) (Result, error)
	String() string
}

// Result is a generic interface for a result returned from a worker
type Result interface{}

// This function creates a waitgroup and a taskQueue for the workerPool.
// It will create "n" workers which will listen for Tasks on the taskQueue.
// It returns the taskQueue and a wait func.
// The clients should add new tasks to this taskQueue
// and call the call the wait func when done.
// Do not separately close the taskQueue, make use of the wait func.
func WorkerPool(ctx context.Context, n int) (chan<- Task, chan<- Result, func()) {
	var wg sync.WaitGroup
	numTasks := 1024
	taskQueue := make(chan Task, numTasks)
	results := make(chan Result, numTasks)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go worker(ctx, &wg, taskQueue, results)
	}

	wait := func() {
		close(taskQueue)
		close(results)
		wg.Wait()
	}

	return taskQueue, results, wait

}

func worker(ctx context.Context, wg *sync.WaitGroup, taskQueue <-chan Task, results chan<- Result) {
	defer wg.Done()
	for task := range taskQueue {
		select {
		case <-ctx.Done():
			return
		default:
			result, err := task.Run(ctx)
			if err != nil {
				log.WithError(err).Debugf("Task failed: %s", task)
			} else if result != nil {
				results <- result
			}
		}
	}
}
