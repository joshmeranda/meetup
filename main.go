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

	configDir, err := os.UserHomeDir()
	if err != nil {
		return meetup.Config{}, fmt.Errorf("could not find user config dir: %w", err)
	}

	configPath := path.Join(configDir, ".meetup", "config.yaml")

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

func MeetingOpen(ctx *cli.Context) error {
	if ctx.NArg() > 2 {
		return fmt.Errorf("too many arguments")
	}

	if ctx.NArg() < 2 {
		return fmt.Errorf("missing required arguments")
	}

	domain := ctx.Args().Get(0)
	name := ctx.Args().Get(1)

	manager, err := GetManager()
	if err != nil {
		return err
	}

	manager.OpenMeeting(meetup.Meeting{
		Name:     name,
		Domain:   domain,
		Date:     ctx.String("date"),
		Template: ctx.String("template"),
	})

	return nil
}

func MeetingList(ctx *cli.Context) error {
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

func MeetingRemove(ctx *cli.Context) error {
	if ctx.NArg() < 3 {
		return fmt.Errorf("missing required arguments")
	}

	if ctx.NArg() > 3 {
		return fmt.Errorf("too many arguments")
	}

	manager, err := GetManager()
	if err != nil {
		return err
	}

	err = manager.RemoveMeeting(meetup.Meeting{
		Name:   ctx.Args().Get(2),
		Domain: ctx.Args().Get(1),
		Date:   ctx.Args().Get(0),
	})
	if err != nil {
		return err
	}

	return nil
}

func TemplateAdd(ctx *cli.Context) error {
	templates := ctx.Args().Slice()
	if len(templates) == 0 {
		return fmt.Errorf("expected template paths, but found none")
	}

	manager, err := GetManager()
	if err != nil {
		return err
	}

	for _, template := range templates {
		if err := manager.AddTemplate(template); err != nil {
			return fmt.Errorf("could not add template: %w", err)
		}
	}

	return nil
}

func TemplateList(ctx *cli.Context) error {
	manager, err := GetManager()
	if err != nil {
		return err
	}

	templates, err := manager.ListTemplates()
	if err != nil {
		return err
	}

	for _, template := range templates {
		fmt.Println(template)
	}

	return nil
}

func TemplateRemove(ctx *cli.Context) error {
	templates := ctx.Args().Slice()
	if len(templates) == 0 {
		return fmt.Errorf("expected template names, but found none")
	}

	manager, err := GetManager()
	if err != nil {
		return err
	}

	if err := manager.RemoveTemplate(templates...); err != nil {
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
				Name:  "meeting",
				Usage: "manage meetings",
				Subcommands: []*cli.Command{
					{
						Name:      "open",
						Usage:     "open an existing or create a new meeting",
						UsageText: "meetup open <domain> <name>",
						Action:    MeetingOpen,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "date",
								Usage: "date of the meeting",
								Value: cli.NewTimestamp(time.Now()).Value().Format(DateFormat),
								Action: func(ctx *cli.Context, date string) error {
									if _, err := time.Parse(DateFormat, date); err != nil {
										return fmt.Errorf("invalid date format: %w", err)
									}

									return nil
								},
							},
							&cli.StringFlag{
								Name:    "template",
								Usage:   "template to use for the meeting",
								Aliases: []string{"t"},
							},
						},
					},
					{
						Name:      "list",
						Aliases:   []string{"ls"},
						Usage:     "list existing meeting",
						UsageText: "meetup list",
						Action:    MeetingList,
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
						UsageText: "meetup remove <date> <domain> <name>",
						Action:    MeetingRemove,
					},
				},
			},
			{
				Name:  "template",
				Usage: "manage meeting templates",
				Subcommands: []*cli.Command{
					{
						Name:   "add",
						Usage:  "add a new template",
						Action: TemplateAdd,
					},
					{
						Name:    "list",
						Aliases: []string{"ls"},
						Usage:   "list existing templates",
						Action:  TemplateList,
					},
					{
						Name:    "remove",
						Aliases: []string{"rm"},
						Usage:   "remove a template",
						Action:  TemplateRemove,
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
