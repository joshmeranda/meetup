package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gobwas/glob"
	meetup "github.com/joshmeranda/meetup/pkg"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

const (
	DateFormat = "2006-01-02"
)

var Version string

func LoadManagerConfig() (meetup.Config, error) {
	config, err := meetup.DefaultConfig()
	if err != nil {
		return meetup.Config{}, fmt.Errorf("could not create default config: %w", err)
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return meetup.Config{}, fmt.Errorf("could not find user config dir: %w", err)
	}

	configPath := path.Join(configDir, "meetup", "config.yaml")

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}

		return meetup.Config{}, fmt.Errorf("could not read config file: %w", err)
	}

	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return meetup.Config{}, fmt.Errorf("could not parse config file: %w", err)
	}

	return config, nil
}

func GetManager() (meetup.Manager, error) {
	config, err := LoadManagerConfig()
	if err != nil {
		return meetup.Manager{}, fmt.Errorf("could not create manager: %w", err)
	}

	return meetup.NewManager(config)
}

func Open(ctx *cli.Context) error {
	manager, err := GetManager()
	if err != nil {
		return err
	}

	var name, domain string

	switch ctx.NArg() {
	case 2:
		domain = ctx.Args().Get(1)
		name = ctx.Args().First()
	case 1:
		name = ctx.Args().First()
	case 0:
		return fmt.Errorf("missing required meeting name")
	default:
		return fmt.Errorf("too many arguments")
	}

	manager.AddMeeting(meetup.Meeting{
		Name:   name,
		Domain: domain,
		Date:   ctx.Timestamp("date").Format(DateFormat),
	})

	return nil
}

func List(ctx *cli.Context) error {
	manager, err := GetManager()
	if err != nil {
		return err
	}

	meetings, err := manager.ListMeetings(meetup.MeetingWildcard{
		Name:   glob.MustCompile(ctx.String("name")),
		Date:   glob.MustCompile(ctx.String("date")),
		Domain: glob.MustCompile(ctx.String("domain")),
	})

	if err != nil {
		return err
	}

	for _, meeting := range meetings {
		fmt.Println(meeting)
	}

	return nil
}

func Remove(ctx *cli.Context) error {
	manager, err := GetManager()
	if err != nil {
		return err
	}

	err = manager.RemoveMeeting(meetup.Meeting{
		Name:   ctx.String("name"),
		Date:   ctx.String("date"),
		Domain: ctx.String("domain"),
	})
	if err != nil {
		return err
	}

	return nil
}

func Run(args []string) error {
	app := cli.App{
		Name:    "meetup",
		Version: Version,
		Usage:   "meetup is a tool for managing meeting notes",
		Commands: []*cli.Command{
			{
				Name:      "open",
				Usage:     "open an existing or create a new meeting",
				UsageText: "meetup new <name> [domain]",
				Action:    Open,
				Flags: []cli.Flag{
					&cli.TimestampFlag{
						Name:   "date",
						Layout: DateFormat,
						Usage:  "date of the meeting",
						Value:  cli.NewTimestamp(time.Now()),
					},
				},
			},
			{
				Name:      "list",
				Aliases:   []string{"ls"},
				Usage:     "list existing meeting",
				UsageText: "meetup list",
				Action:    List,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "date",
						Usage: "date of the meeting as a wildcard",
						Value: "*",
					},
					&cli.StringFlag{
						Name:  "name",
						Usage: "the name of the meeting as a wildcard",
						Value: "*",
					},
					&cli.StringFlag{
						Name:  "domain",
						Usage: "the domain of the meeting as a wildcard",
						Value: "*",
					},
				},
			},
			{
				Name:      "remove",
				Aliases:   []string{"rm"},
				Usage:     "remove an existing meeting",
				UsageText: "meetup remove <name> [domain]",
				Action:    Remove,
				Flags: []cli.Flag{
					&cli.TimestampFlag{
						Name:   "date",
						Layout: DateFormat,
						Usage:  "date of the meeting",
						Value:  cli.NewTimestamp(time.Now()),
					},
				},
			},
		},
	}

	return app.Run(args)
}

func main() {
	if err := Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
