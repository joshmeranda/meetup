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
		templateSrc = path.Join(exampleDir, "template.md")
		templateDst = path.Join(meetupDir, meetup.TemplateDir, "template.md")

		manager, err = meetup.NewManager(meetup.Config{
			RootDir:       meetupDir,
			DefaultDomain: "default",
			Editor:        []string{"touch"},
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
		Expect(templates).To(ContainElement("template.md"))
	})

	It("can open a meeting with a template", func() {
		err = manager.OpenMeeting(meetup.Meeting{
			Name:     "example-meeting",
			Date:     "2021-01-01",
			Domain:   "meetup.template.test",
			Template: "template.md",
		})
		Expect(err).ToNot(HaveOccurred())

		expected := "2021-01-01 meetup.template.test example-meeting"
		meetingFile := path.Join(meetupDir, "meetup", "template", "test", "2021-01-01", "example-meeting")
		data, err := os.ReadFile(meetingFile)

		Expect(err).ToNot(HaveOccurred())
		Expect(string(data)).To(Equal(expected))
	})

	It("can remove templates", func() {
		err = manager.RemoveTemplate("template.md")
		Expect(err).ToNot(HaveOccurred())
		Expect(path.Join(meetupDir, meetup.TemplateDir, "template.md")).ShouldNot(BeAnExistingFile())
	})
})
