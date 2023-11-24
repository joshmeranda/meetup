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
	var manager meetup.Manager
	var err error

	meetupDir := "meetup-test"

	BeforeAll(func() {
		manager, err = meetup.NewManager(meetup.Config{
			RootDir:       meetupDir,
			DefaultDomain: "default",
			GroupBy:       meetup.GroupByDomain,
			Editor:        []string{"touch"},
		})

		Expect(err).ToNot(HaveOccurred())
	})

	AfterAll(func() {
		os.RemoveAll(meetupDir)
	})

	It("can open meetings", func() {
		err = manager.OpenMeeting(meetup.Meeting{
			Name:   "sample",
			Domain: "",
			Date:   "2021-01-01",
		})
		Expect(err).ToNot(HaveOccurred())

		err = manager.OpenMeeting(meetup.Meeting{
			Name:   "sample",
			Domain: "single.double",
			Date:   "2021-01-01",
		})
		Expect(err).ToNot(HaveOccurred())

		err = manager.OpenMeeting(meetup.Meeting{
			Name:   "sample",
			Domain: "single",
			Date:   "2021-01-01",
		})
		Expect(err).ToNot(HaveOccurred())

		Expect(path.Join(meetupDir, "default", "2021-01-01", "sample")).Should(BeAnExistingFile())
		Expect(path.Join(meetupDir, "single", "2021-01-01", "sample")).Should(BeAnExistingFile())
		Expect(path.Join(meetupDir, "single", "double", "2021-01-01", "sample")).Should(BeAnExistingFile())
	})

	It("can list meetings", func() {
		meetings, err := manager.ListMeetings(meetup.MeetingWildcard{
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

	It("can remove meetings", func() {
		err = manager.RemoveMeeting(meetup.Meeting{
			Name:   "sample",
			Domain: "",
			Date:   "2021-01-01",
		})
		Expect(err).ToNot(HaveOccurred())

		err = manager.RemoveMeeting(meetup.Meeting{
			Name:   "sample",
			Domain: "single",
			Date:   "2021-01-01",
		})
		Expect(err).ToNot(HaveOccurred())

		err = manager.RemoveMeeting(meetup.Meeting{
			Name:   "sample",
			Domain: "single.double",
			Date:   "2021-01-01",
		})
		Expect(err).ToNot(HaveOccurred())
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
