package main

import (
	"flag"
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"wcgo/constants"
	"wcgo/service"
	"wcgo/util"
	"wcgo/worker"
)

func startWork(files []string, words []string, numWorkers float64) {
	counter := service.NewConcurrentCounter()
	jobs := make(chan string, constants.JobBuffer)
	wg := sync.WaitGroup{}

	for i := 1; i <= int(numWorkers); i++ {
		go worker.CounterWorker(jobs, words, counter, &wg)
	}

	for _, file := range files {
		fmt.Println(file)
		jobs <- file
	}
	close(jobs)

	wg.Wait()
	fmt.Print(counter.Store)
}

func main() {
	var (
		dirPtr     = flag.String("dir", constants.NoDirectory, "Directory path")
		workersPtr = flag.String("workers", strconv.Itoa(int(constants.MinWorkers)), "Number of workers")
		wordsPtr   = flag.String("words", "", "Set of words to be tracked (csv)")
	)

	flag.Parse()

	dir := filepath.Clean(*dirPtr)
	if !util.IsValidDirectory(dir) {
		return
	}

	numWorkers := constants.MinWorkers
	if num, err := strconv.Atoi(*workersPtr); err == nil {
		requiredWorkers := float64(num)
		if requiredWorkers > 0 {
			numWorkers = math.Min(constants.MaxWorkers, requiredWorkers)
		}
	}

	files := util.GetFiles(dir)
	words := strings.Split(*wordsPtr, ",")
	startWork(files, words, numWorkers)
}
