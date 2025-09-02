package worker

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"wcgo/service"
)

func isFileExcluded(path string, excludeFilePatterns []string) bool {
	for _, excludePattern := range excludeFilePatterns {
		if strings.Contains(path, excludePattern) {
			return true
		}
	}
	return false
}

func CounterWorker(jobs chan string, words []string, excludeFilePatterns []string, c *service.ConcurrentCounter, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for path := range jobs {
		if isFileExcluded(path, excludeFilePatterns) {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Failed to load %v\n", path)
			fmt.Println(err)
			continue
		}
		content := string(data)
		for _, word := range words {
			count := strings.Count(content, word)
			c.Increment(word, uint64(count))
		}
	}
}
