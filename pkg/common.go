package pkg

import (
	"sync"

	log "github.com/sirupsen/logrus"
)

func Consume(ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info("Result of stack ", <-ch)
}

func Remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
