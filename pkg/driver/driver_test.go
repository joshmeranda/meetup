package driver_test

import (
	"os"
	"testing"

	"github.com/joshmeranda/meetup/pkg/driver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleDriver(t *testing.T) {
	defer func() { require.NoError(t, os.RemoveAll("test")) }()

	simpleDriver, err := driver.NewDriver(driver.DriverConfig{
		DriverBackend: driver.DriverBackendSimple,
		SimpleConfig: &driver.SimpleDriverConfig{
			Command: []string{"touch", "test/file-0"},
		},
	})
	require.NoError(t, err)

	simpleDriver.Open("test/file-1", "test/file-2")

	assert.FileExists(t, "test/file-0")
	assert.FileExists(t, "test/file-1")
	assert.FileExists(t, "test/file-2")
}
