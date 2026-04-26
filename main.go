package main

import _ "github.com/lib/pq"

import (
	"database/sql"
	"os"
	"log"
	"github.com/Jesbr/BlogAggregator/internal/config"
	"github.com/Jesbr/BlogAggregator/internal/database"
)

func main() {
	// read config
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// open connection to the database
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	// check if database is reachable
	err = db.Ping()
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	// initialize queries
	dbQueries := database.New(db)

	// create State
	s := &config.State{
		Config: &cfg,
		DB: dbQueries,
	}

	// create commands instance and register handler
	cmds := config.NewCommands()
	cmds.Register("login", config.HandlerLogin)
	cmds.Register("register", config.HandlerRegister)

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