package applog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GetAppDataDir returns the platform-specific directory for application data
func GetAppDataDir(appName string) (string, error) {
	var basePath string

	switch runtime.GOOS {
	case "darwin":
		// macOS: ~/Library/Application Support/[AppName]/logs
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		basePath = filepath.Join(home, "Library", "Application Support", appName, "logs")
	case "windows":
		// Windows: %APPDATA%\[AppName]\logs
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
		basePath = filepath.Join(appData, appName, "logs")
	case "linux":
		// Linux: ~/.config/[AppName]/logs
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		basePath = filepath.Join(home, ".config", appName, "logs")
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return "", err
	}

	return basePath, nil
}
