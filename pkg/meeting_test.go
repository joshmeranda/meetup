package meetup_test

import (
	"os"
	"path"

	"github.com/gobwas/glob"
	meetup "github.com/joshmeranda/meetup/pkg"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ManageMeeting", Ordered, func() {
	var meetupDir string
	var manager meetup.Manager
	var err error

	BeforeAll(func() {
		meetupDir, err = os.MkdirTemp("", "meetup-test")
		Expect(err).ToNot(HaveOccurred())

		manager, err = meetup.NewManager(meetup.Config{
			RootDir: meetupDir,
			Editor:  []string{"touch"},
			DefaultMetadata: meetup.Metadata{
				GroupBy: meetup.GroupByDomain,
			},
		})

		Expect(err).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		os.RemoveAll(meetupDir)
	})

	When("there are no meetings", func() {
		It("can list meetings", func() {
			meetings, err := manager.ListMeetings(meetup.MeetingQuery{
				Date:   glob.MustCompile("*"),
				Name:   glob.MustCompile("*"),
				Domain: glob.MustCompile("*"),
			})
			expected := []meetup.Meeting{}

			Expect(err).ToNot(HaveOccurred())
			Expect(meetings).To(ConsistOf(expected))
		})

		It("cannot remove non-existent meetings", func() {
			err = manager.RemoveMeeting(meetup.Meeting{
				Name:   "i-dont-exist",
				Domain: "no.exist",
				Date:   "2021-01-01",
			})
			Expect(err).To(HaveOccurred())
		})
	})

	It("can open meetings", func() {
		for _, meeting := range testMeetings {
			Expect(manager.OpenMeeting(meeting)).ToNot(HaveOccurred())
		}

		Expect(path.Join(meetupDir, "default", "2021-01-01", "sample")).Should(BeAnExistingFile())
		Expect(path.Join(meetupDir, "single", "2021-01-01", "sample")).Should(BeAnExistingFile())
		Expect(path.Join(meetupDir, "single", "double", "2021-01-01", "sample")).Should(BeAnExistingFile())
	})

	It("can list meetings", func() {
		meetings, err := manager.ListMeetings(meetup.MeetingQuery{
			Date:   glob.MustCompile("*"),
			Name:   glob.MustCompile("*"),
			Domain: glob.MustCompile("*double"),
		})
		expected := []meetup.Meeting{
			{
				Name:   "sample",
				Domain: "single.double",
				Date:   "2021-01-01",
			},
		}

		Expect(err).ToNot(HaveOccurred())
		Expect(meetings).To(ConsistOf(expected))
	})

	It("can re-organize meetup dirs", func() {
		err = manager.UpdateMeetingGroupBy(meetup.GroupByDate)
		Expect(err).ToNot(HaveOccurred())

		data, err := os.ReadFile(path.Join(meetupDir, meetup.MetadataFilename))
		Expect(err).ToNot(HaveOccurred())
		Expect(string(data)).To(Equal("group_by: date\n"))

		Expect(path.Join(meetupDir, "2021-01-01", "default", "sample")).Should(BeAnExistingFile())
		Expect(path.Join(meetupDir, "2021-01-01", "single", "sample")).Should(BeAnExistingFile())
		Expect(path.Join(meetupDir, "2021-01-01", "single", "double", "sample")).Should(BeAnExistingFile())
	})

	It("can remove meetings", func() {
		for _, meeting := range testMeetings {
			Expect(manager.RemoveMeeting(meeting)).ToNot(HaveOccurred())
		}
	})
})
