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

type Downtime struct {
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
}

func (d Downtime) String() string {
	return fmt.Sprintf("%s - %s (%s)", d.StartTime, d.EndTime, d.Duration)
}

type DowntimeLog struct {
	Downtimes []Downtime `json:"downtimes"`
}

type Report struct {
	StartEvent  Event       `json:"start_event"`
	EndEvent    Event       `json:"end_event"`
	DowntimeLog DowntimeLog `json:"downtime_log"`
}

func (r Report) String() string {
	return fmt.Sprintf("Start Event: %s\nEnd Event: %s\nDowntime Log: %v", r.StartEvent, r.EndEvent, r.DowntimeLog)
}
func GetDowntimeLog(eventLog []Event) DowntimeLog {
	var downtimeLog DowntimeLog
	var downtime Downtime
	for _, event := range eventLog {

		if event.Description == "Succeeded" || event.Description == "Failed" {
			downtime.EndTime = event.Time
			downtime.Duration = downtime.EndTime.Sub(downtime.StartTime)
			downtimeLog.Downtimes = append(downtimeLog.Downtimes, downtime)
			downtime = Downtime{}
		}
	}
	return downtimeLog
}
func GetReport(eventLog []Event) Report {
	var report Report
	report.StartEvent = eventLog[0]
	report.EndEvent = eventLog[len(eventLog)-1]
	report.DowntimeLog = GetDowntimeLog(eventLog)
	return report
}
