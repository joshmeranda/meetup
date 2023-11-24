package meetup

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

type GroupStrategy string

const (
	GroupByDomain GroupStrategy = "domain"
	GroupByDate   GroupStrategy = "date"
)

type Config struct {
	RootDir       string        `yaml:"root_dir"`
	DefaultDomain string        `yaml:"default_domain"`
	GroupBy       GroupStrategy `yaml:"group_by"`
	Editor        []string      `yaml:"editor"`
}

func DefaultConfig() (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("could not find user home dir: %w", err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	return Config{
		RootDir:       path.Join(homeDir, ".meetup"),
		DefaultDomain: "default",
		GroupBy:       GroupByDomain,
		Editor:        []string{editor},
	}, nil
}

type Manager struct {
	Config

	baseCmd *exec.Cmd
}

func NewManager(config Config) (Manager, error) {
	if len(config.Editor) == 0 {
		return Manager{}, fmt.Errorf("editor cannot be empty")
	}

	path, args := config.Editor[0], config.Editor[1:]
	path, err := exec.LookPath(path)
	if err != nil {
		return Manager{}, fmt.Errorf("could not find editor: %w", err)
	}

	cmd := exec.Command(path, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	return Manager{
		Config: config,

		baseCmd: cmd,
	}, nil
}
