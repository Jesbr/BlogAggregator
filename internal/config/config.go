package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"fmt"
	"context"
	"time"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/Jesbr/BlogAggregator/internal/database"
	"database/sql"
)

const configFileName = ".gatorconfig.json"

// structs
type Config struct {
	DBURL			string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

type State struct {
	Config *Config
	DB  *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	handlers map[string]func(*State, Command) error
}

// commands methods
func (c *Commands) Register(name string, handler func(*State, Command) error) {
	c.handlers[name] = handler
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, exists := c.handlers[cmd.Name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}

// functions
func NewCommands() *Commands {
	return &Commands{
		handlers: make(map[string]func(*State, Command) error),
	}
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(home, configFileName)
	return fullPath, nil
}

func Read() (Config, error) {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return write(*cfg)
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username required")
	}
	username := cmd.Args[0]

	// check if user in database
	_, err := s.DB.GetUser(context.Background(), username)
		if err != nil {
			if err == sql.ErrNoRows {
			return fmt.Errorf("user %s does not exist", username)
		}
		return err
	}

	// set to active user
	err = s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("User set to %s\n", username)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username required")
	}

	name := cmd.Args[0]
	now := time.Now()

	user, err := s.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	})

	// duplicate check
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return fmt.Errorf("user already exists")
		}
		return err
	}

	// set current user in config
	err = s.Config.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Printf("User created: %s\n", name)
	fmt.Printf("DEBUG user: %+v\n", user)

	return nil
}