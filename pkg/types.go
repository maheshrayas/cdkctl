package pkg

import (
	"context"
)

type Deployer struct {
	ctx         context.Context
	config      *Config
	toolkit     *string
	environment *string
	args        []string
	batch       *string
}

type Config struct {
	Stacks  []string `json:"stacks"`
	verbose bool     `json:"verbose"`
	trace   bool     `json:"trace"`
}
