package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/alessio/shellescape"
	log "github.com/sirupsen/logrus"
)

func (r *Deployer) Deploy() {
	var status bool = true
	var currentStatus = make(map[string]string)
	var deployIndependentStacks []string
	var m Stackgroup
	x := make(map[string][]string)
	for _, y := range r.config.Stacks {
		currentStatus[y.ID] = "NOTCOMPLETED"
		if len(y.Dependson) == 0 {
			x[y.ID] = y.Name
			deployIndependentStacks = append(deployIndependentStacks, y.Name...)
		}
	}
	m.stackgp = x
	r.dependent = currentStatus
	r.deployStacks(deployIndependentStacks, &status, m)
	for {
		stackToExecute := make([]string, 0)
		for _, ev := range r.config.Stacks {
			var completed int = 0
			//check if the id status is completed
			if r.dependent[ev.ID] != "COMPLETED" && len(ev.Dependson) > 0 {
				// find the ids whose dependson is completed
				for _, y := range ev.Dependson {
					if r.dependent[y.(string)] == "COMPLETED" {
						completed++
					}
				}
				if completed == len(ev.Dependson) {
					x[ev.ID] = ev.Name
					m.stackgp = x
					stackToExecute = append(stackToExecute, ev.Name...)
				}
			}
		}
		if len(stackToExecute) > 0 {
			r.deployStacks(stackToExecute, &status, m)
		} else {
			break
		}
	}
	defer r.checkStatus(&status)
}

func (r *Deployer) deployStacks(y []string, status *bool, stackgroup Stackgroup) {
	var wg sync.WaitGroup
	message := make(chan string, len(y))
	for _, stackname := range y {
		go r.runCdkDeploy(stackname, message, status, &stackgroup)
		wg.Add(1)
		go Consume(message, &wg)
		// wait for 1 sec so that we do  not hit aws cloudformation api limit
		time.Sleep(2 * time.Second)
	}
	wg.Wait()
	defer r.checkStackCompletion(&stackgroup)
}

func (r *Deployer) checkStackCompletion(stackgroup *Stackgroup) {
	for key := range stackgroup.stackgp {
		if len(stackgroup.stackgp[key]) == 0 {
			r.dependent[key] = "COMPLETED"
		}
	}
}

func Initialize(ctx context.Context, stacks, toolkit, argsFile, prefix *string) (*Deployer, error) {
	var runStack Config
	stackConfig, err := ioutil.ReadFile(*stacks)
	if err != nil {
		log.Fatal("Error while reading stacks config file")
	}
	err = json.Unmarshal(stackConfig, &runStack)
	if err != nil {
		log.Fatal("Error while parsing stacks config file")
	}

	// load runtime parameters
	var s []string
	if len(*argsFile) > 0 {
		s = getContextArgs(s, argsFile)
	}
	r := &Deployer{
		ctx:     ctx,
		config:  &runStack,
		toolkit: toolkit,
		args:    s,
		prefix:  prefix,
	}
	return r, nil
}

func (r *Deployer) runCdkDeploy(stackName string, message chan<- string, status *bool, stackgroup *Stackgroup) {
	var deployStack string
	if *r.prefix != "" {
		deployStack = *r.prefix + "-" + stackName
	} else {
		deployStack = stackName
	}
	log.Info("Deploying ", deployStack)
	var cdkRun = []string{"cdk", "deploy", deployStack, "--require-approval", "never", "--toolkit-stack-name", *r.toolkit}
	if len(r.args) > 0 {
		cdkRun = append(cdkRun, r.args...)
	}
	cmd := exec.Command(cdkRun[0], cdkRun[1:]...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Info(fmt.Sprint(err) + ": " + string(stdoutStderr))
		*status = false
		r.failedStack = append(r.failedStack, deployStack)
	} else {
		//loop thru the map and remove successful stack
		var newstacks []string
		for id, element := range stackgroup.stackgp {
			for _, stacks := range element {
				if stacks == stackName {
					newstacks = Remove(element, stackName)
				}
			}
			stackgroup.stackgp[id] = newstacks
		}
	}
	message <- string(stdoutStderr)
	log.Info("Finished deploying", stackName)
}

func getContextArgs(s []string, argsFile *string) []string {
	args, err := ioutil.ReadFile(*argsFile)
	if err != nil {
		log.Fatal("Error while reading args config file")
	}
	var argsData map[string]string
	err = json.Unmarshal(args, &argsData)
	if err != nil {
		log.Fatal("Error while parsing args file")
	}
	for key, value := range argsData {
		if strings.Index(value, "env.") > -1 {
			runes := []rune(value)
			value = os.Getenv(string(runes[4:]))
		}
		fmt.Println((value))
		s = append(s, "--context")
		s = append(s, key+"="+shellescape.Quote(value))
	}
	return s
}

func (r *Deployer) checkStatus(status *bool) {
	if !*status {
		log.Fatal("some of the stacks failed %s, please visit logs", r.failedStack)
	} else {
		log.Info("Successful execution of stacks")
	}
}
