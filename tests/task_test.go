package meetup_test

import (
	"path"

	"github.com/gobwas/glob"
	meetup "github.com/joshmeranda/meetup/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func ToPtr[T any](t T) *T {
	return &t
}

var _ = Describe("Task", Ordered, func() {
	var manager meetup.Manager
	var err error

	BeforeEach(func() {
		manager, err = meetup.NewManager(meetup.Config{
			RootDir: path.Join(meetupSampleDir, "group-by-domain"),
			Editor:  []string{"touch"},
			DefaultMetadata: meetup.Metadata{
				GroupBy: meetup.GroupByDomain,
			},
		})
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("can list all tasks", func() {
		tasks, err := manager.Tasks(meetup.TaskQuery{
			Meeting: meetup.MeetingQuery{
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
					Domain: "triple",
				},
				Complete:    false,
				Description: "do something for triple-sample",
			},
			{
				Meeting: meetup.Meeting{
					Name:   "sample",
					Date:   "2021-01-01",
					Domain: "triple",
				},
				Complete:    true,
				Description: "make schedule for triple-sample",
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
				Meeting: meetup.MeetingQuery{
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
						Domain: "triple",
					},
					Complete:    true,
					Description: "make schedule for triple-sample",
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
				Meeting: meetup.MeetingQuery{
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
						Domain: "triple",
					},
					Complete:    false,
					Description: "do something for triple-sample",
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
				Meeting: meetup.MeetingQuery{
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
						Domain: "triple",
					},
					Complete:    false,
					Description: "do something for triple-sample",
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
