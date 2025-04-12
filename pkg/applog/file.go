package applog

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

// processLogFileNameTemplate replaces template placeholders in the log file name
func processLogFileNameTemplate(template string) string {
	now := time.Now()

	// Replace timestamp placeholder with formatted current time
	// Formats: {{timestamp}}, {{ timestamp }}, {{timestamp:format}}, {{ timestamp:format }}
	tsRegex := regexp.MustCompile(`\{\{\s*timestamp(?::([^}]+))?\s*\}\}`)
	template = tsRegex.ReplaceAllStringFunc(template, func(match string) string {
		submatch := tsRegex.FindStringSubmatch(match)
		format := "20060102-150405" // Default format: YYYYMMDD-HHMMSS
		if len(submatch) > 1 && submatch[1] != "" {
			format = submatch[1]
		}
		return now.Format(format)
	})

	// Replace pid placeholder with current process ID
	// Formats: {{pid}}, {{ pid }}
	pidRegex := regexp.MustCompile(`\{\{\s*pid\s*\}\}`)
	template = pidRegex.ReplaceAllString(template, strconv.Itoa(os.Getpid()))

	// Replace hostname placeholder with machine hostname
	// Formats: {{hostname}}, {{ hostname }}
	hostnameRegex := regexp.MustCompile(`\{\{\s*hostname\s*\}\}`)
	if hostname, err := os.Hostname(); err == nil {
		template = hostnameRegex.ReplaceAllString(template, hostname)
	} else {
		template = hostnameRegex.ReplaceAllString(template, "unknown-host")
	}

	// Replace username placeholder with current user
	// Formats: {{username}}, {{ username }}
	usernameRegex := regexp.MustCompile(`\{\{\s*username\s*\}\}`)
	if currentUser, err := getCurrentUsername(); err == nil {
		template = usernameRegex.ReplaceAllString(template, currentUser)
	} else {
		template = usernameRegex.ReplaceAllString(template, "unknown-user")
	}

	return template
}

// getCurrentUsername tries to get the current username
func getCurrentUsername() (string, error) {
	// First try environment variables which are usually set
	for _, envVar := range []string{"USER", "USERNAME"} {
		if username := os.Getenv(envVar); username != "" {
			return username, nil
		}
	}

	// Fallback to os/user if available
	userInfo, err := userCurrent()
	if err != nil {
		return "", err
	}
	return userInfo.Username, nil
}

// userCurrent is a wrapper around user.Current for easier testing/mocking
var userCurrent = func() (*userInfo, error) {
	return &userInfo{Username: "user"}, nil // This will be replaced with actual implementation
}

// userInfo provides a minimal interface needed for username
type userInfo struct {
	Username string
}

func init() {
	// Set up the actual user.Current implementation
	// This avoids having to import user in tests
	userCurrent = func() (*userInfo, error) {
		// In a real implementation, this would use os/user.Current()
		// For simplicity, we're just getting the USER/USERNAME env var
		username := os.Getenv("USER")
		if username == "" {
			username = os.Getenv("USERNAME") // Windows typically uses USERNAME
		}
		if username == "" {
			return nil, fmt.Errorf("could not determine current user")
		}
		return &userInfo{Username: username}, nil
	}
}
