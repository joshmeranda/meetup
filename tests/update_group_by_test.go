package meetup_test

import (
	"os"
	"path"

	meetup "github.com/joshmeranda/meetup/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/otiai10/copy"
	"gopkg.in/yaml.v3"
)

var _ = Describe("UpdateMeetingGroupBy", func() {
	var config *meetup.Config

	BeforeEach(func() {
		os.TempDir()
		meetupDir, err := os.MkdirTemp("", "meetup-test")
		Expect(err).ToNot(HaveOccurred())

		config = &meetup.Config{
			RootDir:         meetupDir,
			Editor:          []string{"touch"},
			DefaultMetadata: meetup.Metadata{},
		}
	})

	AfterEach(func() {
		os.RemoveAll(config.RootDir)
	})

	It("can update group strategy to "+string(meetup.GroupByDate), func() {
		err := copy.Copy(path.Join(meetupSampleDir, "group-by-domain"), config.RootDir, copy.Options{
			OnDirExists: func(src string, dest string) copy.DirExistsAction {
				return copy.Replace
			},
		})
		Expect(err).ToNot(HaveOccurred())

		manager, err := meetup.NewManager(*config)
		Expect(err).ToNot(HaveOccurred())

		err = manager.UpdateMeetingGroupBy(meetup.GroupByDate)
		Expect(err).ToNot(HaveOccurred())

		Expect(path.Join(manager.RootDir, "single")).ToNot(BeADirectory())
		Expect(path.Join(manager.RootDir, "triple")).ToNot(BeADirectory())

		Expect(path.Join(manager.RootDir, "2021-01-01", "single", "sample")).To(BeAnExistingFile())
		Expect(path.Join(manager.RootDir, "2021-01-01", "single", "double", "sample")).To(BeAnExistingFile())
		Expect(path.Join(manager.RootDir, "2021-01-01", "triple", "sample")).To(BeAnExistingFile())

		Expect(path.Join(manager.RootDir, meetup.TemplateDirName)).To(BeADirectory())
		Expect(path.Join(manager.RootDir, meetup.MetadataFilename)).To(BeAnExistingFile())

		data, err := os.ReadFile(path.Join(config.RootDir, meetup.MetadataFilename))
		Expect(err).ToNot(HaveOccurred())

		metadata := meetup.Metadata{}
		err = yaml.Unmarshal(data, &metadata)
		Expect(err).ToNot(HaveOccurred())

		Expect(metadata.GroupBy).To(Equal(meetup.GroupByDate))
	})

	It("can update group strategy to "+string(meetup.GroupByDomain), func() {
		err := copy.Copy(path.Join(meetupSampleDir, "group-by-date"), config.RootDir, copy.Options{
			OnDirExists: func(src string, dest string) copy.DirExistsAction {
				return copy.Replace
			},
		})
		Expect(err).ToNot(HaveOccurred())

		manager, err := meetup.NewManager(*config)
		Expect(err).ToNot(HaveOccurred())

		err = manager.UpdateMeetingGroupBy(meetup.GroupByDomain)
		Expect(err).ToNot(HaveOccurred())

		Expect(path.Join(manager.RootDir, "2021-01-01")).ToNot(BeADirectory())

		Expect(path.Join(manager.RootDir, "single", "2021-01-01", "sample")).To(BeAnExistingFile())
		Expect(path.Join(manager.RootDir, "single", "double", "2021-01-01", "sample")).To(BeAnExistingFile())
		Expect(path.Join(manager.RootDir, "triple", "2021-01-01", "sample")).To(BeAnExistingFile())

		data, err := os.ReadFile(path.Join(config.RootDir, meetup.MetadataFilename))
		Expect(err).ToNot(HaveOccurred())

		metadata := meetup.Metadata{}
		err = yaml.Unmarshal(data, &metadata)
		Expect(err).ToNot(HaveOccurred())

		Expect(metadata.GroupBy).To(Equal(meetup.GroupByDomain))
	})
})
