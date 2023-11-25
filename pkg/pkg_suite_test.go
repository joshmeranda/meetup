package meetup_test

import (
	"path"
	"testing"

	meetup "github.com/joshmeranda/meetup/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	exampleDir = path.Join("..", "examples")

	testMeetings = []meetup.Meeting{
		{
			Name:   "sample",
			Domain: "default",
			Date:   "2021-01-01",
		}, {
			Name:   "sample",
			Domain: "single.double",
			Date:   "2021-01-01",
		}, {
			Name:   "sample",
			Domain: "single",
			Date:   "2021-01-01",
		},
	}
)

func TestPkg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Manager Suite")
}
