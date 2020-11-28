package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/maheshrayas/cdkctl/pkg"
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
		stackPrefixHelp = `prefix the stack name`
	)
	app := kingpin.New(os.Args[0], help)
	deploy := app.Command("deploy", "Deploy cdk stack")

	args := deploy.Flag("args", runtimeArgsHelp).String()
	destroy := app.Command("destroy", "Destroy cdk stack")

	toolKit := app.Flag("tool-kit", toolKitNameHelp).Required().String()
	stackPrefix := app.Flag("stacks-prefix", stackPrefixHelp).String()
	stackFile := app.Flag("stacks-file", stackFileHelp).Required().String()

	if _, err := app.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	c, err := pkg.Initialize(ctx, stackFile, toolKit, args, stackPrefix)
	if err != nil {
		log.Fatal(err)
	}
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// deploy cdk
	case deploy.FullCommand():
		c.Deploy()

	// destroy cdk
	case destroy.FullCommand():
		c.Destroy()
	}

}
