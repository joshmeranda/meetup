package meetup

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"gopkg.in/yaml.v3"
)

type GroupStrategy string

const (
	GroupByDomain GroupStrategy = "domain"
	GroupByDate   GroupStrategy = "date"

	MetadataFilename = ".metadata.yaml"
)

type Metadata struct {
	GroupBy         GroupStrategy     `yaml:"group_by"`
	DomainTemplates map[string]string `yaml:"domain_templates"`
}

func DefaultMetadata() Metadata {
	return Metadata{
		GroupBy: GroupByDomain,
	}
}

type Config struct {
	RootDir         string   `yaml:"root_dir"`
	Editor          []string `yaml:"editor"`
	DefaultMetadata Metadata `yaml:"default_metadata"`
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
		RootDir:         path.Join(homeDir, ".meetup"),
		Editor:          []string{editor},
		DefaultMetadata: DefaultMetadata(),
	}, nil
}

type Manager struct {
	Config

	baseCmd  *exec.Cmd
	metadata Metadata
}

func NewManager(config Config) (Manager, error) {
	data, err := os.ReadFile(path.Join(config.RootDir, MetadataFilename))
	if err != nil && !os.IsNotExist(err) {
		return Manager{}, fmt.Errorf("could not read metadata file: %w", err)
	}

	if data == nil {
		data = make([]byte, 0)
	}

	metadata := config.DefaultMetadata
	if err := yaml.Unmarshal(data, &metadata); err != nil {
		return Manager{}, fmt.Errorf("could not load metadata: %w", err)
	}

	if len(config.Editor) == 0 {
		return Manager{}, fmt.Errorf("editor cannot be empty")
	}

	path, args := config.Editor[0], config.Editor[1:]

	cmd := exec.Command(path, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	return Manager{
		Config: config,

		baseCmd:  cmd,
		metadata: metadata,
	}, nil
}

func (m *Manager) SyncMetadata() error {
	data, err := yaml.Marshal(m.metadata)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error marshalling metadata: %w", err)
	}

	metadataFile := path.Join(m.RootDir, MetadataFilename)

	if err := os.WriteFile(metadataFile, data, 0644); err != nil {
		return fmt.Errorf("error writing metadata: %w", err)
	}

	return nil
}
