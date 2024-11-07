package mem

import "fmt"

const unit = 1024

func GetHumanReadableSize(bytes int) string {
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := unit, 0

	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	value := float64(bytes) / float64(div)

	return fmt.Sprintf("%.1f %cB", value, "KMGTPE"[exp])
}
