package httphelper

import (
	"fmt"
	"regexp"
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
)

const (
	defaultBrand         = `" Not;A Brand";v="99"`
	defaultEngine        = "Chromium"
	defaultEngineVersion = "92"
)

var (
	browserRegex  = regexp.MustCompile(`(?i)(chrome|firefox|safari|edg|opera)/(\d+\.\d+|\d+\.\d+\.\d+)`)
	platformRegex = regexp.MustCompile(`(?i)(Windows|Mac OS X|Linux)`)
)

func GetUserAgent() string {
	return browser.Computer()
}

// GetSecCHUA constructs the Sec-CH-UA header value based on the User-Agent.
func GetSecCHUA(userAgent string) string {
	engine := defaultEngine               // Default engine
	engineVersion := defaultEngineVersion // Default engine version

	// Match the User-Agent to detect browser name and version
	matches := browserRegex.FindAllStringSubmatch(userAgent, -1)
	if len(matches) > 0 {
		browser := capitalizeFirstLetter(matches[0][1]) // Capitalize browser name (e.g., Chrome, Firefox)
		version := matches[0][2]

		// Set the appropriate engine and engine version based on the browser
		switch browser {
		case "Firefox":
			engine = "Gecko"
			engineVersion = version
		case "Safari":
			engine = "WebKit"
			engineVersion = version
		case "Chrome", "Edg", "Opera":
			engine = "Chromium"
			engineVersion = version
		}

		return fmt.Sprintf(`"%s";v="%s", %s, "%s";v="%s"`, browser, version, defaultBrand, engine, engineVersion)
	}

	// Return default Sec-CH-UA value if no browser is detected
	return fmt.Sprintf(`%s, "%s";v="%s"`, defaultBrand, defaultEngine, defaultEngineVersion)
}

// GetSecCHUAPlatform extracts platform from user agent string.
func GetSecCHUAPlatform(userAgent string) string {
	match := platformRegex.FindString(userAgent)

	switch match {
	case "Windows":
		return "Windows"
	case "Mac OS X":
		return "macOS"
	case "Linux":
		return "Linux"
	case "Android":
		return "Android"
	case "iPhone", "iPad":
		return "iOS"
	default:
		return "Unknown"
	}
}

// capitalizeFirstLetter capitalizes only the first letter of a string.
func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return strings.ToUpper(s)
	}

	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}
