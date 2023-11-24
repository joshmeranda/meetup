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

func (m Manager) createMeetingFile(meeting Meeting) (string, error) {
	meetingPath := m.pathForMeeting(meeting)
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
func (m Manager) OpenMeeting(meeting Meeting) error {
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

func (m Manager) ListMeetings(mw MeetingWildcard) ([]Meeting, error) {
	meetings := []Meeting{}

	filepath.WalkDir(m.RootDir, func(path string, entry fs.DirEntry, err error) error {
		if entry == nil {
			return nil
		}

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
