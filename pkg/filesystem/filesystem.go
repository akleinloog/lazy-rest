package filesystem

import (
	"github.com/akleinloog/lazy-rest/config"
	"github.com/spf13/afero"
	"os"
)

var (
	_fs           afero.Fs
	configuration *config.Config
)

type Fs struct {
}

// New returns the file system.
func New(config *config.Config) Fs {
	configuration = config
	return Fs{}
}

// Exists indicates if a location exists on the file system (could be a file or a directory).
func (*Fs) Exists(location string) (bool, error) {
	return afero.Exists(fs(), location)
}

// IsDir indicates if a location is a directory or not.
func (*Fs) IsDir(location string) (bool, error) {
	return afero.IsDir(fs(), location)
}

// DirExists indicates if a directory exists or not.
func (*Fs) DirExists(location string) (bool, error) {
	return afero.DirExists(fs(), location)
}

// ReadFile returns the content of a file.
func (*Fs) ReadFile(location string) ([]byte, error) {
	return afero.ReadFile(fs(), location)
}

func (*Fs) ReadDir(location string) ([]os.FileInfo, error) {
	return afero.ReadDir(fs(), location)
}

func (*Fs) WriteFile(location string, data []byte) error {
	return afero.WriteFile(fs(), location, data, 0644)
}

func (*Fs) Remove(location string) error {
	return fs().Remove(location)
}

func fs() afero.Fs {
	if _fs == nil {
		if configuration.InMemory() {
			_fs = afero.NewMemMapFs()
		} else {
			_fs = afero.NewBasePathFs(afero.NewOsFs(), "./data")
		}
	}
	return _fs
}
