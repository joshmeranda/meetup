package meetup

import (
	"bufio"
	"os"
	"strings"
	"sync"

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
	Meeting     MeetingQuery
	Complete    *bool
	Description glob.Glob
}

func (t TaskQuery) Match(task Task) bool {
	return t.Meeting.Match(task.Meeting) &&
		(t.Complete == nil || *t.Complete == task.Complete) &&
		t.Description.Match(task.Description)
}

func (m *Manager) searchMeeting(meeting Meeting, query TaskQuery) ([]Task, error) {
	tasks := []Task{}

	meetingPath := m.pathForMeeting(m.metadata.GroupBy, meeting, false)

	meetingFile, err := os.Open(meetingPath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(meetingFile)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		var task Task
		switch {
		case strings.HasPrefix(line, DefaultTaskPrefix):
			task = Task{
				Meeting:     meeting,
				Complete:    false,
				Description: strings.TrimPrefix(line, DefaultTaskPrefix),
			}
		case strings.HasPrefix(line, DefaultTaskCompletedPrefix):
			task = Task{
				Meeting:     meeting,
				Complete:    true,
				Description: strings.TrimPrefix(line, DefaultTaskCompletedPrefix),
			}
		default:
			continue
		}

		if !query.Match(task) {
			continue
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (m *Manager) Tasks(query TaskQuery) ([]Task, error) {
	meetings, err := m.ListMeetings(query.Meeting)
	if err != nil {
		return nil, err
	}

	tasks := []Task{}

	jq := NewJobQueue(5)
	errChan := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(len(meetings))

	for _, meeting := range meetings {
		select {
		case err := <-errChan:
			return nil, err
		default:
		}

		meeting := meeting

		jq.Run(func() {
			foundTasks, err := m.searchMeeting(meeting, query)
			if err != nil {
				errChan <- err
				return
			}

			tasks = append(tasks, foundTasks...)
			wg.Done()
		})
	}

	wg.Wait()

	return tasks, nil
}
