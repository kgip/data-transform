package to

type TaskState struct {
	State    int    `json:"state"`
	Status   int    `json:"status"`
	Count    int    `json:"count"`
	Progress int    `json:"progress"`
	Cause    string `json:"cause"`
}
