package worker

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"wcgo/service"
)

func CounterWorker(jobs chan string, words []string, c *service.ConcurrentCounter, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for path := range jobs {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Failed to load %v\n", path)
			fmt.Println(err)
			continue
		}
		fields := strings.FieldsSeq(string(data))
		for v := range fields {
			if slices.Contains(words, v) {
				c.Increment(v)
			}
		}
	}
}
