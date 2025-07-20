package internal

type CreateTaskParams struct {
	Name, Description string

	StartTime, EndTime string

	RepeatHourly string
	RepeatDaily  []int
}
