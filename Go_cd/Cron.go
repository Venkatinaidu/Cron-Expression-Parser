package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter the cron expression with command to execute:")

	if scanner.Scan() {
		cronExpr := scanner.Text()

		parts := strings.Fields(cronExpr)
		if len(parts) < 6 {
			fmt.Println("Invalid cron expression")
			return
		}

		minute := parts[0]
		hour := parts[1]
		dayOfMonth := parts[2]
		month := parts[3]
		dayOfWeek := parts[4]
		command := strings.Join(parts[5:], " ")

		minuteValues, err := parseCronPart(minute, 0, 59)
		if err != nil {
			fmt.Println("Error parsing minute:", err)
			return
		}
		hourValues, err := parseCronPart(hour, 0, 23)
		if err != nil {
			fmt.Println("Error parsing hour:", err)
			return
		}
		dayOfMonthValues, err := parseCronPart(dayOfMonth, 1, 31)
		if err != nil {
			fmt.Println("Error parsing day of month:", err)
			return
		}
		monthValues, err := parseCronPart(month, 1, 12)
		if err != nil {
			fmt.Println("Error parsing month:", err)
			return
		}
		dayOfWeekValues, err := parseCronPart(dayOfWeek, 0, 6)
		if err != nil {
			fmt.Println("Error parsing day of week:", err)
			return
		}

		minuteStr := intArrayToString(minuteValues)
		hourStr := intArrayToString(hourValues)
		dayOfMonthStr := intArrayToString(dayOfMonthValues)
		monthStr := intArrayToString(monthValues)
		dayOfWeekStr := intArrayToString(dayOfWeekValues)

		fmt.Println("Minute:", minuteStr)
		fmt.Println("Hour:", hourStr)
		fmt.Println("Day of Month:", dayOfMonthStr)
		fmt.Println("Month:", monthStr)
		fmt.Println("Day of Week:", dayOfWeekStr)
		fmt.Println("Command to Execute:", command)
	} else {
		fmt.Println("Error reading input:", scanner.Err())
	}
}

func parseCronPart(part string, min, max int) ([]int, error) {
	var values []int
	if part == "*" {
		for i := min; i <= max; i++ {
			values = append(values, i)
		}
		return values, nil
	}

	if strings.Contains(part, "/") {
		parts := strings.Split(part, "/")
		step, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		rangePart := parts[0]
		if rangePart == "*" {
			for i := min; i <= max; i += step {
				values = append(values, i)
			}
		} else {
			rangeValues, err := parseRange(rangePart, min, max)
			if err != nil {
				return nil, err
			}
			for _, val := range rangeValues {
				if (val-min)%step == 0 {
					values = append(values, val)
				}
			}
		}
		return values, nil
	}

	parts := strings.Split(part, ",")
	for _, part := range parts {
		if strings.Contains(part, "-") {
			rangeValues, err := parseRange(part, min, max)
			if err != nil {
				return nil, err
			}
			values = append(values, rangeValues...)
		} else {
			value, err := strconv.Atoi(part)
			if err != nil {
				return nil, err
			}
			if value < min || value > max {
				return nil, fmt.Errorf("value %d out of range [%d, %d]", value, min, max)
			}
			values = append(values, value)
		}
	}
	return values, nil
}

func parseRange(part string, min, max int) ([]int, error) {
	var values []int
	rangeParts := strings.Split(part, "-")
	start, err := strconv.Atoi(rangeParts[0])
	if err != nil {
		return nil, err
	}
	end, err := strconv.Atoi(rangeParts[1])
	if err != nil {
		return nil, err
	}
	if start < min || end > max || start > end {
		return nil, fmt.Errorf("range %d-%d out of bounds [%d, %d]", start, end, min, max)
	}
	for i := start; i <= end; i++ {
		values = append(values, i)
	}
	return values, nil
}

func intArrayToString(values []int) string {
	var sb strings.Builder
	for i, val := range values {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(strconv.Itoa(val))
	}
	return sb.String()
}
