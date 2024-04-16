package digit

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var digitalSizes = regexp.MustCompile("(?i)(KB|MB|GB|TB)$")

func ParseSize(sizeStr string) (int64, error) {
	digitalSuffix := digitalSizes.FindString(sizeStr)
	if digitalSuffix == "" {
		return 0, fmt.Errorf("Invalid string format: %s", sizeStr)
	}

	numStr := strings.TrimSuffix(sizeStr, digitalSuffix)

	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Cannot parse number: %s", numStr)
	}

	digitalSuffixUpper := strings.ToUpper(digitalSuffix)

	var size int64

	switch digitalSuffixUpper {
	case "KB":
		size = num * 1024
	case "MB":
		size = num * 1024 * 1024
	case "GB":
		size = num * 1024 * 1024 * 1024
	case "TB":
		size = num * 1024 * 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("Incorrect size unit: %s", digitalSuffix)
	}

	return size, nil
}
