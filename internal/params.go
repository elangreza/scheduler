package internal

// old
// type CreateTaskParams struct {
// 	Name, Description string

// 	StartTime, EndTime string

// 	RepeatHourly string
// 	RepeatDaily  []int
// }

type CreateTaskParams struct {
	Name, Description string
}

type UpdateTaskParams struct {
	Name, Description string
}
