package filesystem

import (
	"github.com/akleinloog/lazy-rest/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestWithRealFileSystem(t *testing.T) {

	configuration := config.New()
	fs := New(&configuration)

	content := "Hello"
	location := "tests/case-001"

	WriteReadAndRemoveFile(t, fs, location, content)

	err := os.RemoveAll("./data")
	assert.NoError(t, err, "Error occurred while cleaning up test directory")
}

func TestWithInMemoryFileSystem(t *testing.T) {

	viper.Set("in-memory", true)
	configuration := config.New()
	fs := New(&configuration)

	content := "Hello"
	location := "tests/case-001"

	WriteReadAndRemoveFile(t, fs, location, content)

	viper.Set("in-memory", nil)
}

func WriteReadAndRemoveFile(t *testing.T, fs Fs, location string, content string) {

	err := fs.WriteFile(location, []byte(content))

	if assert.NoError(t, err, "Error occurred while writing to file") {

		exists, err := fs.Exists(location)
		if assert.NoError(t, err, "Error occurred while checking if file exists") {
			assert.True(t, exists, "File does not exist")
		}

		isDir, err := fs.IsDir(location)
		if assert.NoError(t, err, "Error occurred while checking if file exists") {
			assert.False(t, isDir, "Expected file, but was a directory")
		}
	}

	bytes, err := fs.ReadFile(location)

	if assert.NoError(t, err, "Error occurred while reading from file") {

		assert.Equal(t, content, string(bytes), "File content not as expected")
	}

	err = fs.Remove(location)

	if assert.NoError(t, err, "Error occurred while removing file") {

		exists, err := fs.Exists(location)
		if assert.NoError(t, err, "Error occurred while checking if file still exists") {
			assert.False(t, exists, "File still exist")
		}

		directory := path.Dir(location)
		exists, err = fs.Exists(directory)
		if assert.NoError(t, err, "Error occurred while checking if directory still exists") {
			assert.True(t, exists, "Directory no longer exist")
		}
	}
}
