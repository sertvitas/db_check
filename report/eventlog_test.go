package report

import (
	"testing"
	"time"
)

func TestDowntime_String(t *testing.T) {
	type fields struct {
		StartTime time.Time
		EndTime   time.Time
		Duration  time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test 1",
			fields: fields{
				StartTime: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Duration:  0,
			},
			want: "2021-01-01 00:00:00 +0000 UTC - 2021-01-01 00:00:00 +0000 UTC (0s)",
		},
		{
			name: "Test 2",
			fields: fields{
				StartTime: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Duration:  1 * time.Second,
			},
			want: "2021-01-01 00:00:00 +0000 UTC - 2021-01-01 00:00:00 +0000 UTC (1s)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Downtime{
				StartTime: tt.fields.StartTime,
				EndTime:   tt.fields.EndTime,
				Duration:  tt.fields.Duration,
			}
			if got := d.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReport_String(t *testing.T) {
	type fields struct {
		StartEvent  Event
		EndEvent    Event
		DowntimeLog DowntimeLog
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test 1",
			fields: fields{
				StartEvent: Event{
					Time:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Description: "Start",
				},
				EndEvent: Event{
					Time:        time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Description: "End",
				},
				DowntimeLog: DowntimeLog{
					Downtimes: []Downtime{
						{
							StartTime: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							EndTime:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							Duration:  0,
						},
					},
				},
			},
			want: "Start Event: 2021-01-01 00:00:00 +0000 UTC - Start\nEnd Event: 2021-01-01 00:00:00 +0000 UTC - End\nDowntime Log: {Downtimes:[2021-01-01 00:00:00 +0000 UTC - 2021-01-01 00:00:00 +0000 UTC (0s)]}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Report{
				StartEvent:  tt.fields.StartEvent,
				EndEvent:    tt.fields.EndEvent,
				DowntimeLog: tt.fields.DowntimeLog,
			}
			if got := r.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
