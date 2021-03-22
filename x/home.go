package x

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/ory/x/logrusx"
)

// homeDirUnsafe is a low-level function that returns
// the user's home directory from environment
// variables. Careful: if it cannot be determined, an
// empty string is returned. If not accounting for
// that case, use HomeDir() instead; otherwise you
// may end up using the root of the file system.
func homeDirUnsafe() string {
	home := os.Getenv("HOME")
	if home == "" && runtime.GOOS == "windows" {
		drive := os.Getenv("HOMEDRIVE")
		path := os.Getenv("HOMEPATH")
		home = drive + path
		if drive == "" || path == "" {
			home = os.Getenv("USERPROFILE")
		}
	}
	if home == "" && runtime.GOOS == "plan9" {
		home = os.Getenv("home")
	}
	return home
}

// AppConfigDir returns the directory where to store user's config.
//
// If XDG_CONFIG_HOME is set, it returns: $XDG_CONFIG_HOME/ory.
// Otherwise, os.UserConfigDir() is used; if successful, it appends
// "ory" (Windows & Mac) or "ory" (every other OS) to the path.
// If it returns an error, the fallback path "./ory" is returned.
//
// The config directory is not guaranteed to be different from
// AppDataDir().
//
// Unlike os.UserConfigDir(), this function prefers the
// XDG_CONFIG_HOME env var on all platforms, not just Unix.
//
// Ref: https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
func AppConfigDir(l *logrusx.Logger) string {
	if basedir := os.Getenv("XDG_CONFIG_HOME"); basedir != "" {
		return filepath.Join(basedir, "ory")
	}
	basedir, err := os.UserConfigDir()
	if err != nil {
		l.WithError(err).Warn("unable to determine directory for user configuration; falling back to current directory")
		return "./ory"
	}
	subdir := "ory"
	switch runtime.GOOS {
	case "windows", "darwin":
		subdir = "ory"
	}
	return filepath.Join(basedir, subdir)
}

// AppDataDir returns a directory path that is suitable for storing
// application data on disk. It uses the environment for finding the
// best place to store data, and appends a "ory" or "ory" (depending
// on OS and environment) subdirectory.
//
// For a base directory path:
// If XDG_DATA_HOME is set, it returns: $XDG_DATA_HOME/ory; otherwise,
// on Windows it returns: %AppData%/ory,
// on Mac: $HOME/Library/Application Support/ory,
// on Plan9: $home/lib/ory,
// on Android: $HOME/ory,
// and on everything else: $HOME/.local/share/ory.
//
// If a data directory cannot be determined, it returns "./ory"
// (this is not ideal, and the environment should be fixed).
//
// The data directory is not guaranteed to be different from AppConfigDir().
//
// Ref: https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
func AppDataDir() string {
	if basedir := os.Getenv("XDG_DATA_HOME"); basedir != "" {
		return filepath.Join(basedir, "ory")
	}
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("AppData")
		if appData != "" {
			return filepath.Join(appData, "ory")
		}
	case "darwin":
		home := homeDirUnsafe()
		if home != "" {
			return filepath.Join(home, "Library", "Application Support", "ory")
		}
	case "plan9":
		home := homeDirUnsafe()
		if home != "" {
			return filepath.Join(home, "lib", "ory")
		}
	case "android":
		home := homeDirUnsafe()
		if home != "" {
			return filepath.Join(home, "ory")
		}
	default:
		home := homeDirUnsafe()
		if home != "" {
			return filepath.Join(home, ".local", "share", "ory")
		}
	}
	return "./ory"
}
