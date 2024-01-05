package meetup_test

import (
	"os"
	"path"

	meetup "github.com/joshmeranda/meetup/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ManageTemplates", Ordered, func() {
	var manager meetup.Manager
	var err error
	var meetupDir, templateSrc, templateDst string

	BeforeAll(func() {
		meetupDir = "meetup-test"
		templateSrc = path.Join(exampleDir, "templates", "simple.md")
		templateDst = path.Join(meetupDir, meetup.TemplateDirName, "simple.md")

		manager, err = meetup.NewManager(meetup.Config{
			RootDir: meetupDir,
			Editor:  []string{"touch"},
			DefaultMetadata: meetup.Metadata{
				GroupBy: meetup.GroupByDomain,
			},
		})
	})

	AfterAll(func() {
		os.RemoveAll(meetupDir)
	})

	It("can add templates", func() {
		err = manager.AddTemplate(templateSrc)
		Expect(err).ToNot(HaveOccurred())
		Expect(templateDst).Should(BeAnExistingFile())
	})

	It("can list templates", func() {
		templates, err := manager.ListTemplates()
		Expect(err).ToNot(HaveOccurred())
		Expect(templates).To(ContainElement("simple.md"))
	})

	It("can open a meeting with a template", func() {
		err = manager.OpenMeeting(meetup.Meeting{
			Name:     "example-meeting",
			Date:     "2021-01-01",
			Domain:   "meetup.template.test",
			Template: "simple.md",
		})
		Expect(err).ToNot(HaveOccurred())

		expected := "2021-01-01 meetup.template.test Example-M<eeting"
		meetingFile := path.Join(meetupDir, "meetup", "template", "test", "2021-01-01", "example-meeting")
		data, err := os.ReadFile(meetingFile)

		Expect(err).ToNot(HaveOccurred())
		Expect(string(data)).To(Equal(expected))
	})

	It("can remove templates", func() {
		err = manager.RemoveTemplate("simple.md")
		Expect(err).ToNot(HaveOccurred())
		Expect(path.Join(meetupDir, meetup.TemplateDirName, "simple.md")).ShouldNot(BeAnExistingFile())
	})
})
