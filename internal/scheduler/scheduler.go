package scheduler

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("colomn 'repeat' is empty")
	}

	currentDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}

	s := strings.Split(repeat, " ")

	var newDate time.Time
	switch s[0] {
	case "d": // ++
		if len(s) != 2 {
			return "", fmt.Errorf("invalid repeat format: %s", repeat)
		}
		days, err := ParseDays(s[1])
		if err != nil {
			return "", fmt.Errorf("invalid day range: %s", s[1])
		}
		newDate = NextNewDay(now, currentDate, days)
	case "y": // ++
		newDate = NextValidYear(now, currentDate)
	case "w": // ++
		if len(s) != 2 {
			return "", fmt.Errorf("incorrect repeat format: %s", repeat)
		}
		weekDays, err := ParseWeekDays(s[1])
		if err != nil {
			return "", err
		}
		newDate = NextNearestWeekDay(now, currentDate, weekDays)
	case "m": // ++
		if len(s) < 2 || len(s) > 3 {
			return "", fmt.Errorf("incorrect repeat format: %s", repeat)
		}
		days, err := ParseDaysInMonth(s[1])
		if err != nil {
			return "", err
		}
		//months := make([]int, 0, 2)
		/*allMonths := true
		if len(s) == 3 {
			allMonths = false
			months, err = ParsevalidMonths(s[2])
			if err != nil {
				return "", err
			}
		}*/
		if len(s) == 2 {
			newDate = NextValidDayInAllMonths(now, currentDate, days)
		} else {
			months, err := ParsevalidMonths(s[2])
			if err != nil {
				return "", err
			}
			newDate = NextValidDayInMonth(now, currentDate, days, months)
		}

	default:
		return "", fmt.Errorf("invalid repeat format: %s", repeat)
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

func ParseWeekDays(weekDaysStr string) ([]int, error) {
	weekDays := strings.Split(weekDaysStr, ",")
	wDays := make([]int, len(weekDays))
	for idx, weekDay := range weekDays {
		wd, err := strconv.Atoi(weekDay)
		if err != nil {
			return nil, fmt.Errorf("invalid day format: %w", err)
		}
		if wd < 1 || wd > 7 {
			return nil, fmt.Errorf("day must be between 1 and 7")
		}
		wDays[idx] = wd
	}
	return wDays, nil
}

func ParsevalidMonths(monthsStr string) ([]int, error) {
	months := strings.Split(monthsStr, ",")
	result := make([]int, len(months))
	for idx, month := range months {
		m, err := strconv.Atoi(month)
		if err != nil {
			return nil, fmt.Errorf("invalid day format: %w", err)
		}
		if m < 1 || m > 12 {
			return nil, fmt.Errorf("month must be between 1 and 12")
		}
		result[idx] = m
	}
	sort.Ints(result)
	return result, nil
}

func ParseDaysInMonth(daysStr string) ([]int, error) {
	monthDays := strings.Split(daysStr, ",")
	mDays := make([]int, len(monthDays))
	for idx, monthDay := range monthDays {
		md, err := strconv.Atoi(monthDay)
		if err != nil {
			return nil, fmt.Errorf("invalid day format: %w", err)
		}
		if md < -2 || md > 31 || md == 0 {
			return nil, fmt.Errorf("day must be between 1 and 31 or -1, -2")
		}
		mDays[idx] = md
	}
	mDays = customSort(mDays)
	return mDays, nil
}

func NextNewDay(now, date time.Time, days int) time.Time {
	newDate := date.AddDate(0, 0, days)
	if !newDate.After(now) {
		dif := now.Sub(newDate).Hours() / 24

		interval := int(math.Ceil(float64(dif) / float64(days)))
		newDate = newDate.AddDate(0, 0, interval*days)
	}
	return newDate
}

func NextValidYear(now, date time.Time) time.Time {
	newDate := date.AddDate(1, 0, 0)
	if !newDate.After(now) {
		year := now.Year() - newDate.Year()
		newDate = newDate.AddDate(year, 0, 0)
	}
	return newDate
}

func NextNearestWeekDay(now time.Time, date time.Time, weekDays []int) time.Time {
	if len(weekDays) == 0 {
		return date
	}
	distances := make([]int, len(weekDays))

	dif := now.Sub(date).Hours() / 24
	curdate := date
	if dif > 7 {
		curdate = now
	}
	currentWeekDay := int(curdate.Weekday())

	for idx, day := range weekDays {
		daysUntilTarget := (day - currentWeekDay + 7) % 7
		if daysUntilTarget == 0 {
			daysUntilTarget = 7
		}
		distances[idx] = daysUntilTarget
	}
	sort.Ints(distances)

	var newDate time.Time
	for _, dist := range distances {
		newDate = curdate.AddDate(0, 0, dist)
		if newDate.After(now) {
			break
		}
	}

	return newDate
}

func NextValidDayInMonth(now time.Time, date time.Time, days, months []int) time.Time {
	var newDate time.Time
	curMonth := date.Month()
	curDay := date.Day()
	curYear := date.Year()
	var valid bool
	copyMonths := make([]int, 0, len(months))

	for _, month := range months {
		if month >= int(curMonth) {
			if month > int(curMonth) {
				curDay = 0
			}
			if res, ok := checkDay(curYear, time.Month(month), curDay, days); ok {
				newDate = res
				valid = true
				break
			}
			copyMonths = append(copyMonths, month)
		} else {
			copyMonths = append(copyMonths, month)
		}
	}

	if !valid || newDate.Before(now) {
		valid = false
		curMonth := now.Month()
		curDay := now.Day()
		curYear := now.Year()

		for _, month := range months {
			if month == int(curMonth) {
				if res, ok := checkDay(curYear, time.Month(month), curDay, days); ok {
					newDate = res
					valid = true
					break
				}
			} else if month > int(curMonth) {
				if res, ok := checkDay(curYear, time.Month(month), 0, days); ok {
					newDate = res
					valid = true
					break
				}
			}
		}
	}

	if !valid {
		if copyMonths[0]+1 > 12 {
			curMonth = time.Month(1)
			curYear++
		} else {
			curMonth = time.Month(copyMonths[0] + 1)
		}
		if res, ok := checkDay(now.Year()+1, curMonth, 0, days); ok {
			newDate = res
		}
	}

	return newDate
}

func NextValidDayInAllMonths(now time.Time, date time.Time, days []int) time.Time {
	var newDate time.Time
	curMonth := date.Month()
	curDay := date.Day()
	curYear := date.Year()
	var valid bool

	if res, ok := checkDay(curYear, curMonth, curDay, days); ok {
		newDate = res
		valid = true
	} else {
		if int(curMonth)+1 > 12 {
			curMonth = time.Month(1)
			curYear++
		} else {
			curMonth = time.Month(int(curMonth) + 1)
		}
		if res, ok := checkDay(curYear, curMonth, 0, days); ok {
			newDate = res
			valid = true
		}
	}

	if !valid || newDate.Before(now) {
		valid = false
		curMonth := now.Month()
		curDay := now.Day()
		curYear := now.Year()
		if res, ok := checkDay(curYear, curMonth, curDay, days); ok {
			newDate = res
			valid = true
		} else {
			if int(curMonth)+1 > 12 {
				curMonth = time.Month(1)
				curYear++
			} else {
				curMonth = time.Month(int(curMonth) + 1)
			}
			if res, ok := checkDay(curYear, curMonth, 0, days); ok {
				newDate = res
				valid = true
			}
		}
	}

	if !valid {
		if int(curMonth)+1 > 12 {
			curMonth = time.Month(1)
			curYear++
		} else {
			curMonth = time.Month(int(curMonth) + 1)
		}
		if res, ok := checkDay(now.Year()+1, curMonth, 0, days); ok {
			newDate = res
		}
	}

	return newDate
}

func customSort(arr []int) []int {
	// Массив для остальных чисел
	var regular []int
	var minusOneFound, minusTwoFound bool

	// Проход по массиву и отбор значений
	for _, num := range arr {
		if num == -1 {
			minusOneFound = true
		} else if num == -2 {
			minusTwoFound = true
		} else {
			regular = append(regular, num)
		}
	}

	// Сортировка обычных чисел
	sort.Ints(regular)

	// Добавляем -2 и -1 в нужном порядке
	if minusTwoFound {
		regular = append(regular, -2)
	}
	if minusOneFound {
		regular = append(regular, -1)
	}

	return regular
}

func daysInMonth(year int, month time.Month) int {
	nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	return nextMonth.AddDate(0, 0, -1).Day()
}

func checkDay(year int, month time.Month, curDay int, days []int) (time.Time, bool) {
	maxDay := daysInMonth(year, time.Month(month))
	for _, day := range days {
		checkDay := day
		if day == -1 {
			checkDay = maxDay
		} else if day == -2 {
			checkDay = maxDay - 1
		} else if day > maxDay {
			continue
		}
		if checkDay > curDay {
			return time.Date(year, month, checkDay, 0, 0, 0, 0, time.UTC), true
		}
	}
	return time.Time{}, false
}
