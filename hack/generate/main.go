package main

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	meetup "github.com/joshmeranda/meetup/pkg"
)

var (
	MeetupExampleDir = path.Join("tests", "meetup-samples")
	GenerateDir      = path.Join("hack", "generate")
	DefaultTemplate  = "template.md"

	Meetings = []meetup.Meeting{
		{
			Name:     "sample",
			Domain:   "triple",
			Date:     "2021-01-01",
			Template: DefaultTemplate,
		},
		{
			Name:     "sample",
			Domain:   "single",
			Date:     "2021-01-01",
			Template: DefaultTemplate,
		},
		{
			Name:     "sample",
			Domain:   "single.double",
			Date:     "2021-01-01",
			Template: DefaultTemplate,
		},
	}
)

func createMeetupDir(gs meetup.GroupStrategy) error {
	config := meetup.Config{
		RootDir: path.Join(MeetupExampleDir, fmt.Sprintf("group-by-%s", gs)),
		Editor:  []string{"true"}, // noop
		DefaultMetadata: meetup.Metadata{
			GroupBy: gs,
			DomainTemplates: map[string]string{
				"meetup.test": "simple.md",
			},
		},
	}

	manager, err := meetup.NewManager(config)
	if err != nil {
		return fmt.Errorf("could not create new manager: %w", err)
	}

	if err := manager.AddTemplate(path.Join(GenerateDir, DefaultTemplate)); err != nil {
		return fmt.Errorf("could not add template: %w", err)
	}

	for _, meeting := range Meetings {
		if err := manager.OpenMeeting(meeting); err != nil {
			return fmt.Errorf("could not create meeting: %w", err)
		}
	}

	if err := manager.SyncMetadata(); err != nil {
		return fmt.Errorf("could not sync metadata: %w", err)
	}

	return nil
}

func main() {
	slog.Info("cleaning up existing meetup directories")

	if err := os.RemoveAll(MeetupExampleDir); err != nil {
		slog.Error("could not remove existing meetup directories")
		os.Exit(1)
	}

	for _, gs := range []meetup.GroupStrategy{meetup.GroupByDomain, meetup.GroupByDate} {
		if err := createMeetupDir(gs); err != nil {
			slog.Error("could not create meetup dir",
				"group_strategy", gs,
				"error", err,
			)
			os.Exit(1)
		}

		slog.Info("created meetup dir",
			"group_strategy", gs,
		)
	}
}
