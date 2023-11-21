package main

import (
	"fmt"
	"os"
	"path"
	"time"

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

func New(ctx *cli.Context) error {
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

func Run(args []string) error {
	app := cli.App{
		Name:    "meetup",
		Version: Version,
		Usage:   "meetup is a tool for managing meeting notes",
		Commands: []*cli.Command{
			{
				Name:      "new",
				Usage:     "create a new meeting notepad",
				UsageText: "meetup new <name> [domain]",
				Action:    New,
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
