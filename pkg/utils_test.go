package meetup_test

import (
	"fmt"
	"testing"

	meetup "github.com/joshmeranda/meetup/pkg"
	"github.com/stretchr/testify/assert"
)

func TestMeetingFromPath(t *testing.T) {
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
		t.Run(testCase.Name, func(t *testing.T) {
			meeting, err := meetup.MeetingFromPath(testCase.GroupBy, testCase.Path)

			assert.Equal(t, testCase.Error, err)
			assert.Equal(t, testCase.Meeting.Name, meeting.Name)
			assert.Equal(t, testCase.Meeting.Domain, meeting.Domain)
			assert.Equal(t, testCase.Meeting.Date, meeting.Date)
		})
	}
}
