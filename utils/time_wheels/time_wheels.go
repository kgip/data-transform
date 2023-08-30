package time_wheels

import (
	"errors"
	"go.uber.org/atomic"
	"strconv"
	"strings"
	"time"
)

var CronErr = errors.New("cron error")

type YearWheels map[int]MonthWheels

type MonthWheels map[int]DayWheels

type DayWheels map[int]HourWheels

type HourWheels map[int]MinuteWheels

type MinuteWheels map[int]SecondWheels

type SecondWheels map[int]*Task

type Task struct {
	name string
	cron string // * * * */3 12 *
	f    func()
}

type TimeWheels struct {
	wheels       map[int]YearWheels
	active       *atomic.Bool
	taskNames    map[string]struct{}
	removedTasks map[string]struct{}
}

func (t *TimeWheels) AddTask(f func(), name, cron string) {
	task := &Task{name: name, cron: cron, f: f}
}

func (t *TimeWheels) RemoveTask(name string) {
	t.removedTasks[name] = struct{}{}
}

func (t *TimeWheels) Run() {
	t.active = &atomic.Bool{}
	t.active.Store(true)
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			if t.active.Load() {

			} else {
				ticker.Stop()
				break
			}
		}
	}()
}

func (t *TimeWheels) Stop() {
	t.active.Store(false)
}

type TimeUnit int

const (
	Year = TimeUnit(iota)
	Month
	Day
	Hour
	Minute
	Second
)

// ResolveCron 解析cron表达式，返回下一次执行时间
func ResolveCron(cron string) (map[TimeUnit]int, error) {
	timeExpr := strings.Split(cron, " ")
	if len(timeExpr) < 5 {
		return nil, CronErr
	}
	result := map[TimeUnit]int{}
	now := time.Now()
	//* 12 */12
	if strings.Contains(timeExpr[4], "/") {
		split := strings.Split(timeExpr[4], "/")
		if len(split) != 2 || split[0] != "*" {
			return nil, CronErr
		}
		month, err := strconv.ParseInt(split[1], 10, 64)
		if err != nil {
			return nil, CronErr
		}
		date := now.AddDate(0, int(month), 0)
		result[Year] = date.Year()
		result[Month] = int(date.Month())
		result[Day] = date.Day()
		result[Hour] = date.Hour()
		result[Minute] = date.Minute()
		result[Second] = date.Second()
		return result, nil
	} else if timeExpr[4] != "*" {
		month, err := strconv.ParseInt(timeExpr[4], 10, 64)
		if err != nil || month <= 0 || month > 12 {
			return nil, CronErr
		}
		var year = now.Year()
		if int(month) < int(now.Month()) {
			year++
		}
		result[Year] = year
		result[Month] = int(month)
	}
	if strings.Contains(timeExpr[3], "/") {
		if _, ok := result[Year]; ok {
			return nil, CronErr
		}
		split := strings.Split(timeExpr[3], "/")
		if len(split) != 2 || split[0] != "*" {
			return nil, CronErr
		}
		day, err := strconv.ParseInt(split[1], 10, 64)
		if err != nil {
			return nil, CronErr
		}
		date := now.AddDate(0, 0, int(day))
		result[Year] = date.Year()
		result[Month] = int(date.Month())
		result[Day] = date.Day()
		result[Hour] = date.Hour()
		result[Minute] = date.Minute()
		result[Second] = date.Second()
		return result, nil
	} else if timeExpr[3] != "*" {
		day, err := strconv.ParseInt(timeExpr[4], 10, 64)
		if err != nil || day <= 0 {
			return nil, CronErr
		}
		if day > 28 {
			if day == 29 {

			} else if day == 30 {

			} else if day == 31 {

			} else {
				return nil, CronErr
			}
		}
		if month, ok := result[Month]; ok {

		}

	}
	return result, nil
}
