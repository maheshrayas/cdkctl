package pkg

import (
	"context"
)

type Deployer struct {
	ctx         context.Context
	config      *Config
	toolkit     *string
	args        []string
	batch       *string
	dependent   map[string]string
	prefix      *string
	failedStack []string
}

type Processing struct {
	stacks []string
	status map[string]string
}
type Stackgroup struct {
	stackgp map[string][]string
}

type Config struct {
	Stacks []struct {
		ID        string        `json:"id"`
		Name      []string      `json:"name"`
		Dependson []interface{} `json:"dependson"`
		Complete  string
	} `json:"stacks"`
}
