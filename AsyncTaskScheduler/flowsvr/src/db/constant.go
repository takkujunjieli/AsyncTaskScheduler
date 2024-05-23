package db

type TaskEnum int

const (
	TASK_STATUS_PENDING    TaskEnum = 1
	TASK_STATUS_PROCESSING TaskEnum = 2
	TASK_STATUS_SUCC       TaskEnum = 3
	TASK_STATUS_FAILED     TaskEnum = 4
)

const (
	//
	MAX_PRIORITY = 3600 * 24 * 30 * 12
)

// IsValidStatus
func IsValidStatus(status TaskEnum) bool {
	if status == TASK_STATUS_PENDING {
		return true
	}
	if status == TASK_STATUS_PROCESSING {
		return true
	}
	if status == TASK_STATUS_SUCC {
		return true
	}
	if status == TASK_STATUS_FAILED {
		return true
	}
	return false
}
