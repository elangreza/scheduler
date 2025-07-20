package internal

import (
	"reflect"
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {

	mockedTimeNow := "2025-07-20T10:38:23+07:00"
	mockedTimeNowParsed, _ := time.Parse(time.RFC3339, mockedTimeNow)
	mockedTimeNowAfterOneHour := "2025-07-20T11:38:23+07:00"
	mockedTimeNowAfterOneDayParsed, _ := time.Parse(time.RFC3339, mockedTimeNowAfterOneHour)
	type args struct {
		name         string
		description  string
		startTime    string
		endTime      string
		repeatHourly string
		repeatDaily  []int
	}
	tests := []struct {
		name    string
		args    args
		want    *Task
		wantErr bool
	}{
		{
			name: "empty start time",
			args: args{
				name:         "a",
				description:  "a",
				startTime:    "",
				endTime:      "",
				repeatHourly: "",
				repeatDaily:  []int{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed to parsed start time",
			args: args{
				name:         "a",
				description:  "a",
				startTime:    "a",
				endTime:      "",
				repeatHourly: "",
				repeatDaily:  []int{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed to parsed end time",
			args: args{
				name:         "a",
				description:  "a",
				startTime:    mockedTimeNow,
				endTime:      "a",
				repeatHourly: "",
				repeatDaily:  []int{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "name is empty",
			args: args{
				name:         "",
				description:  "a",
				startTime:    mockedTimeNow,
				endTime:      "",
				repeatHourly: "",
				repeatDaily:  []int{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "start date is overlapping the end date",
			args: args{
				name:         "a",
				description:  "a",
				startTime:    mockedTimeNowAfterOneHour,
				endTime:      mockedTimeNow,
				repeatHourly: "",
				repeatDaily:  []int{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed to parsed repeat hourly",
			args: args{
				name:         "a",
				description:  "a",
				startTime:    mockedTimeNow,
				endTime:      "",
				repeatHourly: "xx",
				repeatDaily:  []int{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "repeat hourly is zero",
			args: args{
				name:         "a",
				description:  "a",
				startTime:    mockedTimeNow,
				endTime:      "",
				repeatHourly: "-1s",
				repeatDaily:  []int{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "repeat hourly is overlapping the end time",
			args: args{
				name:         "a",
				description:  "a",
				startTime:    mockedTimeNow,
				endTime:      mockedTimeNowAfterOneHour,
				repeatHourly: "2h",
				repeatDaily:  []int{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid repeatDaily",
			args: args{
				name:         "a",
				description:  "a",
				startTime:    mockedTimeNow,
				endTime:      "",
				repeatHourly: "2h",
				repeatDaily:  []int{-1},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid repeatDaily",
			args: args{
				name:         "a",
				description:  "a",
				startTime:    mockedTimeNow,
				endTime:      "",
				repeatHourly: "2h",
				repeatDaily:  []int{7},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success with hoRepeatHourly and RepeatDaily",
			args: args{
				name:         "a",
				description:  "a",
				startTime:    mockedTimeNow,
				endTime:      mockedTimeNowAfterOneHour,
				repeatHourly: "20m",
				repeatDaily:  []int{1, 5, 3},
			},
			want: &Task{
				ID:             0,
				Name:           "a",
				Description:    "a",
				StartTime:      mockedTimeNowParsed,
				EndTime:        mockedTimeNowAfterOneDayParsed,
				RepeatHourly:   "20m",
				RepeatDaily:    []int{1, 3, 5},
				isRoutine:      true,
				repeatInterval: 20 * time.Minute,
				nextRunAt:      time.Time{},
				CreatedAt:      time.Time{},
				UpdatedAt:      time.Time{},
			},
			wantErr: false,
		},
		{
			name: "success no routine",
			args: args{
				name:         "a",
				description:  "a",
				startTime:    mockedTimeNow,
				endTime:      "",
				repeatHourly: "",
				repeatDaily:  []int{},
			},
			want: &Task{
				ID:             0,
				Name:           "a",
				Description:    "a",
				StartTime:      mockedTimeNowParsed,
				EndTime:        time.Time{},
				RepeatHourly:   "",
				RepeatDaily:    []int{},
				isRoutine:      false,
				repeatInterval: 0,
				nextRunAt:      time.Time{},
				CreatedAt:      time.Time{},
				UpdatedAt:      time.Time{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTask(tt.args.name, tt.args.description, tt.args.startTime, tt.args.endTime, tt.args.repeatHourly, tt.args.repeatDaily)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTask() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestTask_GetNextRunAt(t *testing.T) {

	mockedTimeNow := "2025-07-20T10:38:23+07:00"
	mockedTimeNowParsed, _ := time.Parse(time.RFC3339, mockedTimeNow)
	mockedTimeNowAfterOneHour := "2025-07-20T11:38:23+07:00"
	// mockedTimeNowAfterOneDayParsed, _ := time.Parse(time.RFC3339, mockedTimeNowAfterOneHour)
	type fields struct {
		name         string
		description  string
		startTime    string
		endTime      string
		repeatHourly string
		repeatDaily  []int
	}
	type args struct {
		lastRun time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		{
			name: "not a routine",
			fields: fields{
				name:        "a",
				description: "a",
				startTime:   mockedTimeNow,
				endTime:     mockedTimeNowAfterOneHour,
			},
			args: args{
				lastRun: time.Time{},
			},
			want: time.Time{},
		},
		{
			name: "last run is empty and within endtime",
			fields: fields{
				name:         "a",
				description:  "a",
				startTime:    mockedTimeNow,
				endTime:      mockedTimeNowAfterOneHour,
				repeatHourly: "20m",
				repeatDaily:  []int{1, 5, 3},
			},
			args: args{
				lastRun: time.Time{},
			},
			want: mockedTimeNowParsed.Add(20 * time.Minute),
		},
		{
			name: "last run is not empty and within endtime",
			fields: fields{
				name:         "a",
				description:  "a",
				startTime:    mockedTimeNow,
				endTime:      mockedTimeNowAfterOneHour,
				repeatHourly: "10m",
				repeatDaily:  []int{1, 5, 3},
			},
			args: args{
				lastRun: mockedTimeNowParsed.Add(60 * time.Minute),
			},
			want: mockedTimeNowParsed.AddDate(0, 0, 1),
		},
		{
			name: "last run is not empty and within endtime",
			fields: fields{
				name:         "a",
				description:  "a",
				startTime:    mockedTimeNow,
				endTime:      mockedTimeNowAfterOneHour,
				repeatHourly: "10m",
				repeatDaily:  []int{5, 3},
			},
			args: args{
				lastRun: mockedTimeNowParsed.Add(60 * time.Minute),
			},
			want: mockedTimeNowParsed.AddDate(0, 0, 3),
		},
		{
			name: "last run is not empty. within endtime, and repeatDaily is available",
			fields: fields{
				name:         "a",
				description:  "a",
				startTime:    mockedTimeNow,
				endTime:      mockedTimeNowAfterOneHour,
				repeatHourly: "10m",
				repeatDaily:  []int{0},
			},
			args: args{
				lastRun: mockedTimeNowParsed.Add(60 * time.Minute),
			},
			want: mockedTimeNowParsed.AddDate(0, 0, 7),
		},
		{
			name: "last run is not empty. within endtime, and repeatDaily is empty",
			fields: fields{
				name:         "a",
				description:  "a",
				startTime:    mockedTimeNow,
				endTime:      mockedTimeNowAfterOneHour,
				repeatHourly: "10m",
				repeatDaily:  []int{},
			},
			args: args{
				lastRun: mockedTimeNowParsed.Add(60 * time.Minute),
			},
			want: time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewTask(tt.fields.name, tt.fields.description, tt.fields.startTime, tt.fields.endTime, tt.fields.repeatHourly, tt.fields.repeatDaily)
			if err != nil {
				t.Error(err)
				return
			}
			if got := s.GetNextRunAt(tt.args.lastRun); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Task.GetNextRunAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateRunTimeSequence(t *testing.T) {
	type args struct {
		lasSeq time.Time
		seqs   []int
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "a",
			args: args{
				lasSeq: time.Date(2025, 7, 2, 0, 0, 0, 0, time.UTC),
				seqs:   []int{3, 4},
			},
			want: time.Date(2025, 7, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "b",
			args: args{
				lasSeq: time.Date(2025, 7, 2, 0, 0, 0, 0, time.UTC),
				seqs:   []int{3, 5},
			},
			want: time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "c",
			args: args{
				lasSeq: time.Date(2025, 7, 2, 0, 0, 0, 0, time.UTC),
				seqs:   []int{3, 6},
			},
			want: time.Date(2025, 7, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "d",
			args: args{
				lasSeq: time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC),
				seqs:   []int{3, 5},
			},
			want: time.Date(2025, 7, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "e",
			args: args{
				lasSeq: time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC),
				seqs:   []int{5},
			},
			want: time.Date(2025, 7, 11, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "f",
			args: args{
				lasSeq: time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC),
				seqs:   []int{4},
			},
			want: time.Date(2025, 7, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "f",
			args: args{
				lasSeq: time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC),
				seqs:   []int{32},
			},
			want: time.Time{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateRunTimeSequence(tt.args.lasSeq, tt.args.seqs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Seq() = %v, want %v", got, tt.want)
			}
		})
	}
}
