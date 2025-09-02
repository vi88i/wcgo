package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"wcgo/constants"
	"wcgo/service"
	"wcgo/util"
	"wcgo/worker"
)

func startWork(files []string, words []string, excludeFilePatterns []string, numWorkers float64) {
	counter := service.NewConcurrentCounter()
	jobs := make(chan string, constants.JobBuffer)
	wg := sync.WaitGroup{}

	for i := 1; i <= int(numWorkers); i++ {
		go worker.CounterWorker(jobs, words, excludeFilePatterns, counter, &wg)
	}

	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	wg.Wait()
	jsonData, err := json.MarshalIndent(counter.Store, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}

func main() {
	var (
		dirPtr         = flag.String("dir", constants.NoDirectory, "Directory path")
		workersPtr     = flag.String("workers", strconv.Itoa(int(constants.MinWorkers)), "Number of workers")
		wordsPtr       = flag.String("words", "", "Set of words to be tracked (csv)")
		excludeFilePtr = flag.String("exclude", constants.DefaultExcludedFilesPattern, "Comma-separated list of keywords - files containing these keywords in their names will be skipped (e.g., 'test,vendor,node_modules')")
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

	startTime := time.Now()
	var (
		files               = util.GetFiles(dir)
		words               = strings.Split(*wordsPtr, ",")
		excludeFilePatterns = strings.Split(*excludeFilePtr, ",")
	)

	startWork(files, words, excludeFilePatterns, numWorkers)

	elapsedTime := time.Since(startTime)
	fmt.Printf("Elapsed Time: %v\n", elapsedTime)
}
