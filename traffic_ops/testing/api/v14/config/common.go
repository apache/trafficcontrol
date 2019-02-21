package config

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const MilitaryTimeFmt = "15:04"

func Validate24HrTimeRange(rng string) error {
	rangeFormat := regexp.MustCompile(`(\S+)-(\S+)`)
	match := rangeFormat.FindStringSubmatch(rng)
	if match == nil {
		return fmt.Errorf("string %v is not a range", rng)
	}

	_, err := time.Parse(MilitaryTimeFmt, match[1])
	if err != nil {
		return fmt.Errorf("time range must be a 24Hr format")
	}

	_, err = time.Parse(MilitaryTimeFmt, match[2])
	if err != nil {
		return fmt.Errorf("time range must be a 24Hr format")
	}

	return nil
}

// 1) order matters
// 2) maximum values?
func ValidateDHMSTimeFormat(time string) error {

	if time == "" {
		return fmt.Errorf("time string cannot be empty")
	}
	timeFormat := regexp.MustCompile(`^(\d+d)?(\d+h)?(\d+m)?(\d+s)?$`)
	match := timeFormat.FindStringSubmatch(time)
	if match == nil {
		return fmt.Errorf("time format must match sequences of digits followed by units, where time units are: d, h, m, or s\n")
	}

	for i := 1; i < len(match); i++ {
		last := len(match[i]) - 1
		if last == -1 {
			continue
		}
		if _, err := strconv.Atoi(match[i][:last]); err != nil {
			return err
		}
	}

	return nil
}
