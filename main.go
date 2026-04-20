package main

import (
	"os"
	"log"
	"github.com/Jesbr/BlogAggregator/internal/config"
)

func main() {
	// read config
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// create State
	s := &config.State{
		Config: &cfg,
	}

	// create commands instance and register handler
	cmds := config.NewCommands()
	cmds.Register("login", config.HandlerLogin)

	// parse CLI arguments
	args := os.Args
	if len(args) < 2 {
		log.Fatal("no command provided")
	}

	cmdName := args[1]
	cmdArgs := []string{}
	if len(args) > 2 {
		cmdArgs = args[2:]
	}

	// build command
	cmd := config.Command{
		Name: cmdName,
		Args: cmdArgs,
	}

	// run command
	err = cmds.Run(s, cmd)
	if err != nil {
		log.Fatalf("error running command: %v", err)
	}
}