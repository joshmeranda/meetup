package meetup_test

import (
	"testing"

	meetup "github.com/joshmeranda/meetup/pkg"
	"github.com/joshmeranda/meetup/pkg/driver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManagerGroupByDomain(t *testing.T) {
	actual := []string{}

	manager, err := meetup.NewManager(meetup.Config{
		RootDir:       "/",
		DefaultDomain: "default",
		GroupBy:       meetup.GroupByDomain,
		Driver: driver.DriverConfig{
			DriverBackend: driver.DriverBackendCallback,
			CallbackConfig: &driver.CallbackDriverConfig{
				Fn: func(files ...string) error {
					actual = append(actual, files...)
					return nil
				},
			},
		},
	})

	manager.AddMeeting(meetup.Meeting{
		Name:   "sample",
		Domain: "",
		Date:   "2021-01-01",
	})
	manager.AddMeeting(meetup.Meeting{
		Name:   "sample",
		Domain: "single",
		Date:   "2021-01-01",
	})
	manager.AddMeeting(meetup.Meeting{
		Name:   "sample",
		Domain: "single.double",
		Date:   "2021-01-01",
	})

	expected := []string{
		"/default/2021-01-01/sample",
		"/single/2021-01-01/sample",
		"/single/double/2021-01-01/sample",
	}

	require.NoError(t, err)
	assert.ElementsMatch(t, expected, actual)
}

func TestManagerGroupByDate(t *testing.T) {
	actual := []string{}

	manager, err := meetup.NewManager(meetup.Config{
		RootDir:       "/",
		DefaultDomain: "default",
		GroupBy:       meetup.GroupByDate,
		Driver: driver.DriverConfig{
			DriverBackend: driver.DriverBackendCallback,
			CallbackConfig: &driver.CallbackDriverConfig{
				Fn: func(files ...string) error {
					actual = append(actual, files...)
					return nil
				},
			},
		},
	})

	manager.AddMeeting(meetup.Meeting{
		Name:   "sample",
		Domain: "",
		Date:   "2021-01-01",
	})
	manager.AddMeeting(meetup.Meeting{
		Name:   "sample",
		Domain: "single",
		Date:   "2021-01-01",
	})
	manager.AddMeeting(meetup.Meeting{
		Name:   "sample",
		Domain: "single.double",
		Date:   "2021-01-02",
	})

	expected := []string{
		"/2021-01-01/default/sample",
		"/2021-01-01/single/sample",
		"/2021-01-02/single/double/sample",
	}

	require.NoError(t, err)
	assert.ElementsMatch(t, expected, actual)
}
