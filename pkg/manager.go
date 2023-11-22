package meetup

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
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

func (m Meeting) String() string {
	return fmt.Sprintf("%s %s %s", m.Date, m.Domain, m.Name)
}

type MeetingWildcard struct {
	Name   glob.Glob
	Domain glob.Glob
	Date   glob.Glob
}

func (mw MeetingWildcard) Match(m Meeting) bool {
	return mw.Name.Match(m.Name) &&
		mw.Domain.Match(m.Domain) &&
		mw.Date.Match(m.Date)
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
		return path.Join(m.Config.RootDir, meeting.Date, domainComponents, meeting.Name)
	default:
		panic(fmt.Sprintf("unknown group_by: %s", m.Config.GroupBy))
	}
}

func (m Manager) AddMeeting(meeting Meeting) error {
	if meeting.Domain == "" {
		meeting.Domain = m.Config.DefaultDomain
	}

	path := m.pathForMeeting(meeting)
	return m.driver.Open(path)
}

func (m Manager) ListMeetings(mw MeetingWildcard) ([]Meeting, error) {
	meetings := []Meeting{}

	filepath.WalkDir(m.RootDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			meeting, err := MeetingFromPath(m.GroupBy, strings.TrimPrefix(path, m.RootDir))
			if err != nil {
				return err
			}

			if mw.Match(meeting) {
				meetings = append(meetings, meeting)
			}
		}

		return nil
	})

	return meetings, nil
}

func (m Manager) RemoveMeeting(meeting Meeting) error {
	path := m.pathForMeeting(meeting)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("could not delete meeting: %w", err)
	}

	return nil
}
