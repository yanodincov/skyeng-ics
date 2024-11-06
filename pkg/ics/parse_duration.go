package ics

import (
	"fmt"
	"strings"
	"time"
)

const hourPerDay = 24

func ConvertDurationToICS(duration time.Duration) string {
	negative := duration < 0
	if negative {
		duration = -duration
	}

	if duration < time.Second {
		return "PT0S"
	}

	var builder strings.Builder

	builder.WriteString("P")

	days := duration / (hourPerDay * time.Hour)
	duration %= hourPerDay * time.Hour
	hours := duration / time.Hour
	duration %= time.Hour
	minutes := duration / time.Minute
	duration %= time.Minute
	seconds := duration / time.Second

	if days > 0 {
		fmt.Fprintf(&builder, "%dD", days)
	}

	if hours > 0 || minutes > 0 || seconds > 0 {
		builder.WriteString("T")

		if hours > 0 {
			fmt.Fprintf(&builder, "%dH", hours)
		}

		if minutes > 0 {
			fmt.Fprintf(&builder, "%dM", minutes)
		}

		if seconds > 0 {
			fmt.Fprintf(&builder, "%dS", seconds)
		}
	}

	result := builder.String()
	if negative {
		result = "-" + result
	}

	return result
}
