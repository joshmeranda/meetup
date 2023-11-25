package meetup

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobwas/glob"
)

type Meeting struct {
	Name     string
	Date     string
	Domain   string
	Template string
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

func (m *Manager) pathForMeeting(gs GroupStrategy, meeting Meeting, isBackup bool) string {
	var topLevel string
	if isBackup {
		topLevel = path.Join(m.Config.RootDir, ".backup")
	} else {
		topLevel = m.Config.RootDir
	}

	domainComponents := path.Join(strings.Split(meeting.Domain, ".")...)

	switch gs {
	case GroupByDomain:
		return path.Join(topLevel, domainComponents, meeting.Date, meeting.Name)
	case GroupByDate:
		return path.Join(topLevel, meeting.Date, domainComponents, meeting.Name)
	default:
		panic(fmt.Sprintf("unknown group_by: %s", m.metadata.GroupBy))
	}
}

func (m *Manager) fillMeeting(meeting Meeting) Meeting {
	if meeting.Domain == "" {
		meeting.Domain = m.Config.DefaultDomain
	}

	return meeting
}

func (m *Manager) createMeetingFile(meeting Meeting) (string, error) {
	meetingPath := m.pathForMeeting(m.metadata.GroupBy, meeting, false)
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
		templatePath := path.Join(m.Config.RootDir, TemplateDir, meeting.Template)

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
	meeting = m.fillMeeting(meeting)

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

func (m *Manager) ListMeetings(mw MeetingWildcard) ([]Meeting, error) {
	meetings := []Meeting{}

	filepath.WalkDir(m.RootDir, func(path string, entry fs.DirEntry, err error) error {
		if entry == nil {
			return nil
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

	return meetings, nil
}

func (m *Manager) RemoveMeeting(meeting Meeting) error {
	meeting = m.fillMeeting(meeting)
	path := m.pathForMeeting(m.metadata.GroupBy, meeting, false)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("could not delete meeting: %w", err)
	}

	return nil
}

func (m *Manager) UpdateMeetingGroupBy(newGs GroupStrategy) error {
	if m.metadata.GroupBy == newGs {
		return nil
	}

	oldGs := m.metadata.GroupBy

	meetings, err := m.ListMeetings(MeetingWildcard{
		Name:   glob.MustCompile("*"),
		Domain: glob.MustCompile("*"),
		Date:   glob.MustCompile("*"),
	})
	if err != nil {
		return fmt.Errorf("could not list meetings: %w", err)
	}

	rootBackup := m.RootDir + ".backup"
	if err := os.Rename(m.RootDir, rootBackup); err != nil {
		return fmt.Errorf("could not backup root directory: %w", err)
	}

	if err := os.MkdirAll(m.RootDir, 0755); err != nil {
		return fmt.Errorf("could not create root directory: %w", err)
	}

	m.metadata.GroupBy = newGs
	if err := m.SyncMetadata(); err != nil {
		return fmt.Errorf("could not sync metadata: %w", err)
	}

	if err := os.Rename(rootBackup, path.Join(m.RootDir, ".backup")); err != nil {
		return fmt.Errorf("could not move meeting backup, backup can be found at %s: %w", rootBackup, err)
	}

	for _, meeting := range meetings {
		backupMeetingPath := m.pathForMeeting(oldGs, meeting, true)
		meetingPath := m.pathForMeeting(m.metadata.GroupBy, meeting, false)

		if err := os.MkdirAll(path.Dir(meetingPath), 0755); err != nil {
			return fmt.Errorf("could not create meeting directory: %w", err)
		}

		if err := os.Rename(backupMeetingPath, meetingPath); err != nil {
			return fmt.Errorf("could not move meeting: %w", err)
		}
	}

	return nil
}
