package scheduler

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

func parseWeekDays(weekDaysStr string) ([]int, error) {
	weekDays := strings.Split(weekDaysStr, ",")
	wDays := make([]int, len(weekDays))
	for idx, weekDay := range weekDays {
		wd, err := strconv.Atoi(weekDay)
		if err != nil {
			return nil, fmt.Errorf("invalid day format: %v", err)
		}
		if wd < 1 || wd > 7 {
			return nil, fmt.Errorf("day must be between 1 and 7")
		}
		wDays[idx] = wd
	}
	return wDays, nil
}

// нужен рефакт, есть повторяющиеся участки
func NextNearestWeekDay(now time.Time, date time.Time, weekDays []int) time.Time {
	if len(weekDays) == 0 {
		return date
	}
	distances := make([]int, len(weekDays))

	currentWeekDay := int(date.Weekday())

	for idx, day := range weekDays {
		daysUntilTarget := (day - currentWeekDay + 7) % 7
		if daysUntilTarget == 0 {
			daysUntilTarget = 7
		}
		distances[idx] = daysUntilTarget
	}
	sort.Ints(distances)

	minDistance := distances[0]
	maxDistance := distances[len(distances)-1]

	maxDate := date.AddDate(0, 0, maxDistance)
	if !maxDate.After(now) {
		currentWeekDay := int(now.Weekday())

		for idx, day := range weekDays {
			daysUntilTarget := (day - currentWeekDay + 7) % 7
			if daysUntilTarget == 0 {
				daysUntilTarget = 7
			}
			distances[idx] = daysUntilTarget
		}
		sort.Ints(distances)

		daysUntilTarget := (distances[0] - currentWeekDay + 7) % 7
		return now.AddDate(0, 0, daysUntilTarget)
	}

	// если мин дистанс меньше текущей даты, то нужно обойти циклом
	newDate := date.AddDate(0, 0, minDistance)
	if !newDate.After(now) {
		for idx := 1; idx < len(distances); idx++ {
			newDate := date.AddDate(0, 0, distances[idx])
			if newDate.After(now) {
				return newDate
			}
		}
	}

	return newDate
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("colomn 'repeat' is empty")
	}

	currentDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("time cannot pasre: %w", err)
	}
	fmt.Println(currentDate)

	s := strings.Split(repeat, " ")

	var newDate time.Time
	switch s[0] {
	case "d":
		if len(s) != 2 {
			return "", fmt.Errorf("incorrect repeat format: %s", repeat)
		}
		days, err := ParseDays(s[1])
		if err != nil {
			return "", err
		}
		newDate = NextNewDay(now, currentDate, days)
	case "y":
		newDate = currentDate.AddDate(1, 0, 0)
		if !newDate.After(now) {
			year := now.Year() - newDate.Year()
			newDate = newDate.AddDate(year, 0, 0)
		}
	case "w":
		/*if len(s) != 2 {
			return "", fmt.Errorf("incorrect repeat format: %s", repeat)
		}
		weekDays, err := parseWeekDays(s[1])
		if err != nil {
			return "", err
		}
		newDate = NextNearestWeekDay(now, currentDate, weekDays)*/
		fallthrough
	case "m":
		/*if len(s) < 2 || len(s) > 3 {
			return "", fmt.Errorf("incorrect repeat format: %s", repeat)
		}
		result := strings.Split(s[1], ",")
		days := make([]int, len(result))
		for idx, day := range result {
			d, err := strconv.Atoi(day)
			if err != nil {
				return "", fmt.Errorf("invalid day format: %v", err)
			}
			if d < -2 || d > 31 {
				return "", fmt.Errorf("day must be between -2 and 31")
			}
			days[idx] = d
		}
		sort.Ints(days)
		if len(s) == 3 {
			result := strings.Split(s[2], ",")
			months := make([]int, len(result))
			for idx, month := range result {
				m, err := strconv.Atoi(month)
				if err != nil {
					return "", fmt.Errorf("invalid day format: %v", err)
				}
				if m < 1 || m > 12 {
					return "", fmt.Errorf("day must be between 1 and 12")
				}
				months[idx] = m
			}
		}

		currentDay := currentDate.Day()
		currentMonth := currentDate.Month()
		currentYear := currentDate.Year()
		var dateValid bool
		for _, day := range days {
			if day > currentDay {
				dif := day - currentDay
				newDate = currentDate.AddDate(0, 0, dif)
				dateValid = true
				//time.Date(currentYear, currentMonth, day, 0, 0, 0, 0, currentDate.Location())
			}
		}
		if !dateValid {
			nextMonth := currentMonth + 1
			for nextMonth > 12 {
				nextMonth = 1
				currentYear++
			}
			//newDate = currentDate.AddDate(currentYear-currentDate.Year(), 0, days[0])
			newDate = time.Date(currentYear, nextMonth, days[0], 0, 0, 0, 0, currentDate.Location())
		}
		dateValid = false
		if !newDate.After(now) {
			for _, day := range days {
				if day >= currentDay {
					dif := day - currentDay
					newDate = now.AddDate(0, 0, dif)
					dateValid = true
					//time.Date(currentYear, currentMonth, day, 0, 0, 0, 0, currentDate.Location())
				}
			}
			if !dateValid {
				currentMonth = now.Month()
				currentYear = now.Year()
				nextMonth := currentMonth + 1
				for nextMonth > 12 {
					nextMonth = 1
					currentYear++
				}
				newDate = time.Date(currentYear, nextMonth, days[0], 0, 0, 0, 0, currentDate.Location())
			}
		}*/
		fallthrough
	default:
		return "", fmt.Errorf("incorrect repeat format: %s", repeat)
	}

	return newDate.Format("20060102"), nil
}

func ParseDays(daysStr string) (int, error) {
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		return 0, fmt.Errorf("invalid day format: %v", err)
	}
	if days < 1 || days > 400 {
		return 0, fmt.Errorf("day must be between 1 and 400")
	}
	return days, nil
}

func NextNewDay(now, date time.Time, days int) time.Time {
	newDate := date.AddDate(0, 0, days)
	if !newDate.After(now) {
		day := now.Day() - newDate.Day()
		interval := int(math.Ceil(float64(day) / float64(days)))
		newDate = newDate.AddDate(0, 0, interval*days)
	}
	return newDate
}
