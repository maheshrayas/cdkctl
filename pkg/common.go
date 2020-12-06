package pkg

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func Consume(ch <-chan string, wg *sync.WaitGroup, stackname, operation string) {
	defer wg.Done()
	for {
		select {
		case msg := <-ch:
			log.Info("Result of stack ", msg)
			return
		default:
			log.Info(operation + " of stack " + stackname + " in progress .. ")
			time.Sleep(60 * time.Second)
		}
	}
}

func Remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
