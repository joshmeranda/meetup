package meetup

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
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

	return Manager{
		Config: config,

		baseCmd: exec.Command(path, args...),
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

func (m Manager) fillMeeting(meeting Meeting) Meeting {
	if meeting.Domain == "" {
		meeting.Domain = m.Config.DefaultDomain
	}

	return meeting
}

// OpenMeeting opens a meeting in the editor, and creates it if it doesn't not exist.
func (m Manager) OpenMeeting(meeting Meeting) error {
	meeting = m.fillMeeting(meeting)
	meetingPath := m.pathForMeeting(meeting)
	meetingDir := path.Dir(meetingPath)

	if err := os.MkdirAll(meetingDir, 0755); err != nil {
		return fmt.Errorf("could not create meeting directory: %w", err)
	}

	cmd := *m.baseCmd
	cmd.Args = append(cmd.Args, meetingPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not open editor: %w", err)
	}

	return nil
}

func (m Manager) ListMeetings(mw MeetingWildcard) ([]Meeting, error) {
	meetings := []Meeting{}

	filepath.WalkDir(m.RootDir, func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
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
	meeting = m.fillMeeting(meeting)
	path := m.pathForMeeting(meeting)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("could not delete meeting: %w", err)
	}

	return nil
}
