package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/maheshrayas/cdkdeploy/pkg"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	const (
		help            = `blah, `
		stackFileHelp   = `A file containing list of cloudformation stacks to be deployed `
		toolKitNameHelp = ` toolkit name`
		environmentHelp = `Deploy stack`
		runtimeArgsHelp = `cdk.json`
		batchHelp       = `how many stacks must be executed concurrently`
	)
	app := kingpin.New(os.Args[0], help)
	// stack := app.Flag("stack", stackFileHelp).Required().String()
	stackFile := app.Flag("stacks-file", stackFileHelp).Required().String()
	toolKit := app.Flag("tool-kit", toolKitNameHelp).Required().String()
	environment := app.Flag("environment", environmentHelp).String()
	args := app.Flag("args", runtimeArgsHelp).String()
	batch := app.Flag("batch", batchHelp).String()
	if _, err := app.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	c, err := pkg.Initialize(ctx, stackFile, toolKit, environment, args, batch)
	if err != nil {
		log.Fatal(err)
	}
	c.Run()
}
