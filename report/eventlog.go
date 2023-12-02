package report

import (
	"fmt"
	"time"
)

// Event is the struct for the event log
type Event struct {
	Time        time.Time `json:"time"`
	Description string    `json:"description"`
}

// Downtime is the struct for the downtime log
type Downtime struct {
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
}

// String returns a string representation of the event
func (d Downtime) String() string {
	return fmt.Sprintf("%s - %s (%s)", d.StartTime, d.EndTime, d.Duration)
}

// DowntimeLog is the struct for the downtime log
type DowntimeLog struct {
	Downtimes []Downtime `json:"downtimes"`
}

// Report is the struct for the report
type Report struct {
	StartEvent  Event       `json:"start_event"`
	EndEvent    Event       `json:"end_event"`
	DowntimeLog DowntimeLog `json:"downtime_log"`
}

// String returns a string representation of the event
func (r Report) String() string {
	return fmt.Sprintf("Start Event: %s\nEnd Event: %s\nDowntime Log: %v", r.StartEvent, r.EndEvent, r.DowntimeLog)
}

// GetDowntimeLog takes an event log and returns a downtime log
func GetDowntimeLog(eventLog []Event) DowntimeLog {
	var downtimeLog DowntimeLog
	var isUp bool = true
	var downtime Downtime
	for _, event := range eventLog {

		if isUp && event.Description == "Failed" {
			isUp = false
			downtime.StartTime = event.Time
			downtime.EndTime = time.Time{}
			downtime.Duration = 0
		}
		if !isUp && event.Description == "Succeeded" {
			isUp = true
			downtime.EndTime = event.Time
			downtime.Duration = downtime.EndTime.Sub(downtime.StartTime)
			downtimeLog.Downtimes = append(downtimeLog.Downtimes, downtime)
			downtime = Downtime{}
		}
	}
	return downtimeLog
}

// GetReport takes an event log and returns a report
func GetReport(eventLog []Event) Report {
	var report Report
	report.StartEvent = eventLog[0]
	report.EndEvent = eventLog[len(eventLog)-1]
	report.DowntimeLog = GetDowntimeLog(eventLog)
	return report
}
