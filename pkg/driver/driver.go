package driver

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

type DriverBackend string

const (
	DriverBackendSimple DriverBackend = "simple"
)

type SimpleDriverConfig struct {
	Command []string `yaml:"command"`
}

type DriverConfig struct {
	DriverBackend DriverBackend       `yaml:"backend"`
	SimpleConfig  *SimpleDriverConfig `yaml:"simple"`
}

type Driver interface {
	Open(files ...string) error
}

func NewDriver(config DriverConfig) (Driver, error) {
	switch config.DriverBackend {
	case DriverBackendSimple:
		// todo: if config is nil, use EDITOR env to create (right now this will panic is nil)
		return NewSimpleDriver(*config.SimpleConfig)
	default:
		return nil, fmt.Errorf("unknown driver backend: %s", config.DriverBackend)
	}
}

// SimpleDriver is a basic driver that runs a command with the given file as input.
type SimpleDriver struct {
	SimpleDriverConfig

	cmd *exec.Cmd
}

func NewSimpleDriver(config SimpleDriverConfig) (Driver, error) {
	foundPath, err := exec.LookPath(config.Command[0])
	if err != nil {
		return nil, fmt.Errorf("could not create driver from nonexistant file: %w", err)
	}

	cmd := exec.Command(foundPath, config.Command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return SimpleDriver{
		SimpleDriverConfig: config,
		cmd:                cmd,
	}, nil
}

func (d SimpleDriver) Open(files ...string) error {
	newCmd := *d.cmd
	newCmd.Args = append(newCmd.Args, files...)

	for _, f := range files {
		dir := path.Dir(f)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("could not create meeting directory: %w", err)
		}
	}

	if err := newCmd.Run(); err != nil {
		return fmt.Errorf("could not run command: %w", err)
	}

	return nil
}
