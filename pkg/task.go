package meetup

import (
	"bufio"
	"os"
	"strings"

	"github.com/gobwas/glob"
)

const (
	DefaultTaskPrefix          = "- [ ] "
	DefaultTaskCompletedPrefix = "- [x] "
)

type Task struct {
	Meeting     Meeting
	Complete    bool
	Description string
}

type TaskQuery struct {
	Meeting     MeetingWildcard
	Complete    *bool
	Description glob.Glob
}

func (m *Manager) Tasks(query TaskQuery) ([]Task, error) {
	meetings, err := m.ListMeetings(query.Meeting)
	if err != nil {
		return nil, err
	}

	tasks := []Task{}

	for _, meeting := range meetings {
		meetingPath := m.pathForMeeting(m.metadata.GroupBy, meeting, false)

		meetingFile, err := os.Open(meetingPath)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(meetingFile)

		// todo: we probably want to break this into a work group or something
		for scanner.Scan() {
			var task Task
			switch {
			case strings.HasPrefix(scanner.Text(), DefaultTaskPrefix):
				task = Task{
					Meeting:     meeting,
					Complete:    false,
					Description: strings.TrimPrefix(scanner.Text(), DefaultTaskPrefix),
				}
			case strings.HasPrefix(scanner.Text(), DefaultTaskCompletedPrefix):
				task = Task{
					Meeting:     meeting,
					Complete:    true,
					Description: strings.TrimPrefix(scanner.Text(), DefaultTaskCompletedPrefix),
				}
			default:
				continue
			}

			if query.Complete != nil && *query.Complete != task.Complete || !query.Description.Match(task.Description) {
				continue
			}

			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}
