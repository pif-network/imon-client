package task

import "encoding/json"

type TaskState string

const (
	Begin TaskState = "Begin"
	Break TaskState = "Break"
	Back  TaskState = "Back"
	End   TaskState = "End"
	Idle  TaskState = "Idle"
)

func (t TaskState) String() string {
	return string(t)
}

type Task struct {
	BeginTime string    `json:"begin_time"`
	Duration  int       `json:"duration"`
	EndTime   string    `json:"end_time"`
	Name      string    `json:"name"`
	State     TaskState `json:"state"`
}

func (t *Task) UnmarshalJSON(data []byte) error {
	type Alias Task
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type Record struct {
	ID          int    `json:"id"`
	UserName    string `json:"user_name"`
	CurrentTask Task   `json:"current_task"`
	TaskHistory []Task `json:"task_history"`
}

func (t *Record) UnmarshalJSON(data []byte) error {
	type Alias Record
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

type STask struct {
	Name       string `json:"name"`
	Descrition string `json:"description"`
	CreatedAt  string `json:"created_at"`
}

var UpstreamEventType = struct {
	GetSingleRecord   string
	GetAllUserRecords string
}{
	GetSingleRecord:   "get_single_record",
	GetAllUserRecords: "get_all_record",
}
