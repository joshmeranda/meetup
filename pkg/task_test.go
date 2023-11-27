package meetup_test

import (
	"os"
	"path"

	"github.com/gobwas/glob"
	meetup "github.com/joshmeranda/meetup/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testTemplate = `# {{.Name}} - {{.Date}} - {{.Domain}}

## Tasks

- [ ] do something for {{ .Domain }}-{{ .Name }}
- [x] make schedule for {{ .Domain }}-{{ .Name }}
`

func ToPtr[T any](t T) *T {
	return &t
}

var _ = Describe("Task", Ordered, func() {
	var meetupDir string
	var manager meetup.Manager
	var err error

	BeforeEach(func() {
		meetupDir, err = os.MkdirTemp("", "meetup-test")
		Expect(err).ToNot(HaveOccurred())

		manager, err = meetup.NewManager(meetup.Config{
			RootDir: meetupDir,
			Editor:  []string{"touch"},
		})
		Expect(err).ToNot(HaveOccurred())

		templatePath := path.Join(meetupDir, "template.md")
		Expect(os.WriteFile(templatePath, []byte(testTemplate), 0600)).ToNot(HaveOccurred())

		Expect(manager.AddTemplate(templatePath)).ToNot(HaveOccurred())
		Expect(os.Remove(templatePath)).ToNot(HaveOccurred())

		for _, meeting := range testMeetings {
			meeting.Template = "template.md"
			Expect(manager.OpenMeeting(meeting)).ToNot(HaveOccurred())
		}
	})

	AfterEach(func() {
		os.RemoveAll(meetupDir)
	})

	It("can list all tasks", func() {
		tasks, err := manager.Tasks(meetup.TaskQuery{
			Meeting: meetup.MeetingWildcard{
				Date:   glob.MustCompile("*"),
				Name:   glob.MustCompile("*"),
				Domain: glob.MustCompile("*"),
			},
			Complete:    nil,
			Description: glob.MustCompile("*"),
		})
		Expect(err).ToNot(HaveOccurred())

		expected := []meetup.Task{
			{
				Meeting: meetup.Meeting{
					Name:   "sample",
					Date:   "2021-01-01",
					Domain: "default",
				},
				Complete:    false,
				Description: "do something for default-sample",
			},
			{
				Meeting: meetup.Meeting{
					Name:   "sample",
					Date:   "2021-01-01",
					Domain: "default",
				},
				Complete:    true,
				Description: "make schedule for default-sample",
			},

			{
				Meeting: meetup.Meeting{
					Name:   "sample",
					Date:   "2021-01-01",
					Domain: "single",
				},
				Complete:    false,
				Description: "do something for single-sample",
			},
			{
				Meeting: meetup.Meeting{
					Name:   "sample",
					Date:   "2021-01-01",
					Domain: "single",
				},
				Complete:    true,
				Description: "make schedule for single-sample",
			},

			{
				Meeting: meetup.Meeting{
					Name:   "sample",
					Date:   "2021-01-01",
					Domain: "single.double",
				},
				Complete:    false,
				Description: "do something for single.double-sample",
			},
			{
				Meeting: meetup.Meeting{
					Name:   "sample",
					Date:   "2021-01-01",
					Domain: "single.double",
				},
				Complete:    true,
				Description: "make schedule for single.double-sample",
			},
		}

		Expect(tasks).To(ConsistOf(expected))
	})

	When("filtering tasks", func() {
		It("can list completed tasks", func() {
			tasks, err := manager.Tasks(meetup.TaskQuery{
				Meeting: meetup.MeetingWildcard{
					Date:   glob.MustCompile("*"),
					Name:   glob.MustCompile("*"),
					Domain: glob.MustCompile("*"),
				},
				Complete:    ToPtr(true),
				Description: glob.MustCompile("*"),
			})
			Expect(err).ToNot(HaveOccurred())

			expected := []meetup.Task{
				{
					Meeting: meetup.Meeting{
						Name:   "sample",
						Date:   "2021-01-01",
						Domain: "default",
					},
					Complete:    true,
					Description: "make schedule for default-sample",
				},
				{
					Meeting: meetup.Meeting{
						Name:   "sample",
						Date:   "2021-01-01",
						Domain: "single",
					},
					Complete:    true,
					Description: "make schedule for single-sample",
				},
				{
					Meeting: meetup.Meeting{
						Name:   "sample",
						Date:   "2021-01-01",
						Domain: "single.double",
					},
					Complete:    true,
					Description: "make schedule for single.double-sample",
				},
			}

			Expect(tasks).To(ConsistOf(expected))
		})

		It("can list uncompleted tasks", func() {
			tasks, err := manager.Tasks(meetup.TaskQuery{
				Meeting: meetup.MeetingWildcard{
					Date:   glob.MustCompile("*"),
					Name:   glob.MustCompile("*"),
					Domain: glob.MustCompile("*"),
				},
				Complete:    ToPtr(false),
				Description: glob.MustCompile("*"),
			})
			Expect(err).ToNot(HaveOccurred())

			expected := []meetup.Task{
				{
					Meeting: meetup.Meeting{
						Name:   "sample",
						Date:   "2021-01-01",
						Domain: "default",
					},
					Complete:    false,
					Description: "do something for default-sample",
				},
				{
					Meeting: meetup.Meeting{
						Name:   "sample",
						Date:   "2021-01-01",
						Domain: "single",
					},
					Complete:    false,
					Description: "do something for single-sample",
				},
				{
					Meeting: meetup.Meeting{
						Name:   "sample",
						Date:   "2021-01-01",
						Domain: "single.double",
					},
					Complete:    false,
					Description: "do something for single.double-sample",
				},
			}

			Expect(tasks).To(ConsistOf(expected))
		})

		It("can list tasks by description", func() {
			tasks, err := manager.Tasks(meetup.TaskQuery{
				Meeting: meetup.MeetingWildcard{
					Date:   glob.MustCompile("*"),
					Name:   glob.MustCompile("*"),
					Domain: glob.MustCompile("*"),
				},
				Complete:    nil,
				Description: glob.MustCompile("do*"),
			})
			Expect(err).ToNot(HaveOccurred())

			expected := []meetup.Task{
				{
					Meeting: meetup.Meeting{
						Name:   "sample",
						Date:   "2021-01-01",
						Domain: "default",
					},
					Complete:    false,
					Description: "do something for default-sample",
				},
				{
					Meeting: meetup.Meeting{
						Name:   "sample",
						Date:   "2021-01-01",
						Domain: "single",
					},
					Complete:    false,
					Description: "do something for single-sample",
				},
				{
					Meeting: meetup.Meeting{
						Name:   "sample",
						Date:   "2021-01-01",
						Domain: "single.double",
					},
					Complete:    false,
					Description: "do something for single.double-sample",
				},
			}

			Expect(tasks).To(ConsistOf(expected))
		})
	})
})
