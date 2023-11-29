package meetup

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"

	"github.com/gobwas/glob"
	"github.com/otiai10/copy"
)

type Meeting struct {
	Name     string
	Date     string
	Domain   string
	Template string
}

// GetPath retusn the path to the meeting with meetupDir as the root.
func (m Meeting) GetPath(meetupDir string, gs GroupStrategy) string {
	domainComponents := path.Join(strings.Split(m.Domain, ".")...)

	switch gs {
	case GroupByDomain:
		return path.Join(meetupDir, domainComponents, m.Date, m.Name)
	case GroupByDate:
		return path.Join(meetupDir, m.Date, domainComponents, m.Name)
	default:
		panic(fmt.Sprintf("unknown group_by: %s", gs))
	}
}

func (m Meeting) String() string {
	return fmt.Sprintf("%s %s %s", m.Date, m.Domain, m.Name)
}

type MeetingQuery struct {
	Name   glob.Glob
	Domain glob.Glob
	Date   glob.Glob
}

func (mw MeetingQuery) Match(m Meeting) bool {
	return mw.Name.Match(m.Name) &&
		mw.Domain.Match(m.Domain) &&
		mw.Date.Match(m.Date)
}

func (m *Manager) createMeetingFile(meeting Meeting) (string, error) {
	// meetingPath := m.pathForMeeting(m.metadata.GroupBy, meeting, false)
	meetingPath := meeting.GetPath(m.RootDir, m.metadata.GroupBy)

	meetingDir := path.Dir(meetingPath)

	if err := os.MkdirAll(meetingDir, 0755); err != nil {
		return "", fmt.Errorf("could not create meeting directory: %w", err)
	}

	outFile, err := os.Create(meetingPath)
	if os.IsExist(err) {
		return meetingPath, nil
	} else if err != nil {
		return "", fmt.Errorf("could not create meeting file: %w", err)
	}

	defer outFile.Close()

	if meeting.Template != "" {
		templatePath := path.Join(m.Config.RootDir, TemplateDirName, meeting.Template)

		template, err := template.ParseFiles(templatePath)
		if err != nil {
			return "", fmt.Errorf("could not parse template: %w", err)
		}

		if err := template.Execute(outFile, meeting); err != nil {
			return "", fmt.Errorf("could not execute template: %w", err)
		}
	}

	return meetingPath, nil
}

// OpenMeeting opens a meeting in the editor, and creates it if it doesn't not exist.
func (m *Manager) OpenMeeting(meeting Meeting) error {
	meetingPath, err := m.createMeetingFile(meeting)
	if err != nil {
		return fmt.Errorf("could not create meeting file: %w", err)
	}

	cmd := *m.baseCmd
	cmd.Args = append(cmd.Args, meetingPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not open editor: %w", err)
	}

	return nil
}

func (m *Manager) ListMeetings(mw MeetingQuery) ([]Meeting, error) {
	meetings := []Meeting{}

	err := filepath.WalkDir(m.RootDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() && entry.Name() == TemplateDirName {
			return filepath.SkipDir
		}

		if !entry.IsDir() {
			meeting, err := MeetingFromPath(m.metadata.GroupBy, strings.TrimPrefix(path, m.RootDir))
			if err != nil {
				return err
			}

			if mw.Match(meeting) {
				meetings = append(meetings, meeting)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not list meetings: %w", err)
	}

	return meetings, nil
}

func (m *Manager) RemoveMeeting(meeting Meeting) error {
	meetingPath := meeting.GetPath(m.RootDir, m.metadata.GroupBy)

	if err := os.Remove(meetingPath); err != nil {
		return fmt.Errorf("could not delete meeting: %w", err)
	}

	for domainPath := path.Dir(meetingPath); domainPath != m.RootDir; domainPath = path.Dir(domainPath) {
		err := os.Remove(domainPath)
		if err != nil {
			pathErr := err.(*os.PathError)
			if pathErr.Err == syscall.ENOTEMPTY {
				break
			}

			return fmt.Errorf("could not delete meeting: %w", err)
		}
	}

	return nil
}

func (m *Manager) UpdateMeetingGroupBy(newGs GroupStrategy) error {
	if m.metadata.GroupBy == newGs {
		return nil
	}

	oldGs := m.metadata.GroupBy

	meetings, err := m.ListMeetings(MeetingQuery{
		Name:   glob.MustCompile("*"),
		Domain: glob.MustCompile("*"),
		Date:   glob.MustCompile("*"),
	})
	if err != nil {
		return fmt.Errorf("could not list meetings: %w", err)
	}

	for _, meeting := range meetings {
		oldMeetingPath := meeting.GetPath(m.RootDir, oldGs)
		newMeetingPath := meeting.GetPath(m.RootDir, newGs)

		if err := os.MkdirAll(path.Dir(newMeetingPath), 0755); err != nil {
			return fmt.Errorf("could not create meeting directory: %w", err)
		}

		if err := copy.Copy(oldMeetingPath, newMeetingPath); err != nil {
			return fmt.Errorf("could not copy meeting: %w", err)
		}

		if err := m.RemoveMeeting(meeting); err != nil {
			return fmt.Errorf("could not remove meeting: %w", err)
		}
	}

	m.metadata.GroupBy = newGs
	if err := m.SyncMetadata(); err != nil {
		return fmt.Errorf("could not sync metadata: %w", err)
	}

	return nil
}
