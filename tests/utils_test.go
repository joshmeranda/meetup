package meetup_test

import (
	"fmt"

	meetup "github.com/joshmeranda/meetup/pkg"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MeetingFromPath", func() {
	type TestCase struct {
		Name    string
		Path    string
		GroupBy meetup.GroupStrategy
		Meeting meetup.Meeting
		Error   error
	}

	testCases := []TestCase{
		{
			Name:    "FlatDomainByDomain",
			Path:    "/default/2021-01-01/sample",
			GroupBy: meetup.GroupByDomain,
			Meeting: meetup.Meeting{
				Name:   "sample",
				Domain: "default",
				Date:   "2021-01-01",
			},
		},
		{
			Name:    "NestedDomainByDomain",
			Path:    "/non/default/2021-01-01/sample",
			GroupBy: meetup.GroupByDomain,
			Meeting: meetup.Meeting{
				Name:   "sample",
				Domain: "non.default",
				Date:   "2021-01-01",
			},
		},
		{
			Name:    "FlatDomainByDate",
			Path:    "/2021-01-01/non/default/sample",
			GroupBy: meetup.GroupByDate,
			Meeting: meetup.Meeting{
				Name:   "sample",
				Domain: "non.default",
				Date:   "2021-01-01",
			},
		},
		{
			Name:    "NestedDomainByDate",
			Path:    "/2021-01-01/default/sample",
			GroupBy: meetup.GroupByDate,
			Meeting: meetup.Meeting{
				Name:   "sample",
				Domain: "default",
				Date:   "2021-01-01",
			},
		},
		{
			Name:    "TooShortPath",
			Path:    "/2021-01-01/domain",
			GroupBy: meetup.GroupByDate,
			Meeting: meetup.Meeting{},
			Error:   fmt.Errorf("path does not have enough components '/2021-01-01/domain'"),
		},
	}

	for _, testCase := range testCases {
		When(testCase.Name, func() {
			It("parses as expected", func() {
				meeting, err := meetup.MeetingFromPath(testCase.GroupBy, testCase.Path)
				Expect(err).To(Equal(testCase.Error))
				Expect(meeting).To(Equal(testCase.Meeting))
			})
		})
	}
})
