package meetup

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/joshmeranda/meetup/pkg/driver"
)

type GroupStrategy string

const (
	GroupByDomain GroupStrategy = "domain"
	GroupByDate   GroupStrategy = "date"
)

type Config struct {
	RootDir       string              `yaml:"root_dir"`
	DefaultDomain string              `yaml:"default_domain"`
	GroupBy       GroupStrategy       `yaml:"group_by"`
	Driver        driver.DriverConfig `yaml:"driver"`
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
		Driver: driver.DriverConfig{
			DriverBackend: driver.DriverBackendSimple,
			SimpleConfig: &driver.SimpleDriverConfig{
				Command: []string{editor},
			},
		},
	}, nil
}

type Meeting struct {
	Name   string
	Date   string
	Domain string
}

type Manager struct {
	Config
	driver driver.Driver
}

func NewManager(config Config) (Manager, error) {
	driver, err := driver.NewDriver(config.Driver)
	if err != nil {
		return Manager{}, fmt.Errorf("could not create driver: %w", err)
	}

	return Manager{
		Config: config,
		driver: driver,
	}, nil
}

func (m Manager) pathForMeeting(meeting Meeting) string {
	domainComponents := path.Join(strings.Split(meeting.Domain, ".")...)

	switch m.Config.GroupBy {
	case GroupByDomain:
		return path.Join(m.Config.RootDir, domainComponents, meeting.Date, meeting.Name)
	case GroupByDate:
		return path.Join(m.Config.RootDir, meeting.Date, meeting.Domain, meeting.Name)
	default:
		panic(fmt.Sprintf("unknown group_by: %s", m.Config.GroupBy))
	}
}

func (m Manager) AddMeeting(meeting Meeting) error {
	path := m.pathForMeeting(meeting)
	return m.driver.Open(path)
}
