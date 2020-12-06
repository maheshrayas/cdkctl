package pkg

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func (r *Deployer) Destroy() {
	var status bool = true
	var deployIndependentStacks []string
	x := make(map[string][]string)
	for _, y := range r.config.Stacks {
		if len(y.Dependson) == 0 {
			x[y.ID] = y.Name
			deployIndependentStacks = append(deployIndependentStacks, y.Name...)
		}
	}
	r.destroy(deployIndependentStacks, &status)
	defer r.checkStatus(&status)
}

func (r *Deployer) destroy(y []string, status *bool) {
	var wg sync.WaitGroup
	message := make(chan string, len(y))
	for _, stackname := range y {
		go r.runCdkDestroy(stackname, message, status)
		wg.Add(1)
		go Consume(message, &wg, *r.prefix+"-"+stackname, "Deleting")
		// wait for 2 sec so that we do  not hit aws cloudformation api limit
		time.Sleep(2 * time.Second)
	}
	wg.Wait()
}

func (r *Deployer) runCdkDestroy(stackName string, message chan<- string, status *bool) {
	var deployStack string
	if *r.prefix != "" {
		deployStack = *r.prefix + "-" + stackName
	} else {
		deployStack = stackName
	}
	log.Info("Destroying ", deployStack)
	var cdkRun = []string{"cdk", "destroy", deployStack, "--force", "--toolkit-stack-name", *r.toolkit}
	cmd := exec.Command(cdkRun[0], cdkRun[1:]...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Info(fmt.Sprint(err) + ": " + string(stdoutStderr))
		*status = false
		r.failedStack = append(r.failedStack, deployStack)
	}
	message <- string(stdoutStderr)
	log.Info("Finished destroying : ", stackName)
}
