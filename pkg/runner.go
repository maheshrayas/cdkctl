package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/alessio/shellescape"
	log "github.com/sirupsen/logrus"
)

func (r *Deployer) runDependentStacks() {
	//push all the stackname to array

}

var x = 0

func (r *Deployer) Run() {

	var status bool = true
	// var batchsize int
	// var err error
	var tempStatus = make(map[string]string)
	var deployIndependentStacks []string
	for _, y := range r.config.Stacks {
		tempStatus[y.ID] = "NOTCOMPLETED"
		if len(y.Dependson) == 0 {
			deployIndependentStacks = append(deployIndependentStacks, y.Name...)
		}
	}
	r.dependent = tempStatus
	r.executeStacks(deployIndependentStacks, &status)
	for {
		stackToExecute := make([]string, 0)
		for _, ev := range r.config.Stacks {
			var completed int = 0
			//check if the id status is completed
			if r.dependent[ev.ID] != "COMPLETED" && len(ev.Dependson) > 0 {
				// find the ids whose dependson is completed

				for x, _ := range ev.Dependson {
					i := strconv.Itoa(x + 1)
					fmt.Println(r.dependent[i])
					if r.dependent[i] == "COMPLETED" {
						completed = +1
					}
				}
				if completed == len(ev.Dependson) {
					//push the stack names to array
					stackToExecute = append(stackToExecute, ev.Name...)
				}

			}

		}
		if len(stackToExecute) <= 0 {
			break
		}

		r.executeStacks(stackToExecute, &status)

	}

	// 	defer checkStatus(&status)
}

func (r *Deployer) executeStacks(y []string, status *bool) {

	if x == 0 {
		r.dependent["1"] = "COMPLETED"
		r.dependent["4"] = "COMPLETED"
	}
	if x == 1 {
		fmt.Println(x)
		r.dependent["2"] = "COMPLETED"
	}
	if x == 3 {
		r.dependent["3"] = "COMPLETED"
	}
	x++

	// fmt.Println(x)

	// if len(y) > 0 {
	// 	var wg sync.WaitGroup
	// 	message := make(chan string, len(y))
	// 	for _, stackname := range y {
	// 		go r.runCdkDeploy(stackname, message, status)
	// 		wg.Add(1)
	// 		go consume(message, &wg)
	// 	}
	// 	wg.Wait()
	// }
}

func consume(ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	println("Result of stack ", <-ch)
}

func Initialize(ctx context.Context, stacks, toolkit, environment, argsFile, batch *string) (*Deployer, error) {
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
		ctx:         ctx,
		config:      &runStack,
		toolkit:     toolkit,
		environment: environment,
		args:        s,
		batch:       batch,
	}
	return r, nil
}

func (r *Deployer) runCdkDeploy(stackName string, message chan<- string, status *bool) {
	log.Info("Deploying ", stackName)
	var cdkRun = []string{"cdk", "deploy", stackName, "--require-approval", "never", "--toolkit-stack-name", *r.toolkit}
	if len(r.args) > 0 {
		cdkRun = append(cdkRun, r.args...)
	}
	cmd := exec.Command(cdkRun[0], cdkRun[1:]...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Info(fmt.Sprint(err) + ": " + string(stdoutStderr))
		*status = false
	}
	message <- string(stdoutStderr)
	log.Info("Finished processing %s", stackName)
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

func checkStatus(status *bool) {
	if !*status {
		log.Fatal("Some of the stacks failed, please visit logs")
	} else {
		log.Info("Successful execution of stacks")
	}
}
