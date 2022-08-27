package deployment

import "fmt"

const (
	STATUS_WAITING           = "waiting"
	STATUS_DEPLOYING         = "deploying"
	STATUS_AVAILABLE         = "available"
	STATUS_RETIRE_REQUESTING = "retire_requesting"
	STATUS_RETIRED           = "retired"
	STATUS_ERROR             = "error"
)

var (
	LivingStatuses = []string{STATUS_WAITING, STATUS_DEPLOYING, STATUS_AVAILABLE, STATUS_RETIRE_REQUESTING}
)

func ValidateStatus(status string) error {
	if status == STATUS_WAITING {
		return nil
	}
	if status == STATUS_DEPLOYING {
		return nil
	}
	if status == STATUS_AVAILABLE {
		return nil
	}
	if status == STATUS_RETIRE_REQUESTING {
		return nil
	}
	if status == STATUS_RETIRED {
		return nil
	}
	if status == STATUS_ERROR {
		return nil
	}
	return fmt.Errorf("unknown staus: %s", status)
}

const (
	TYPE_ONE      = "one"
	TYPE_EACH     = "each"
	TYPE_NODEPLOY = ""
)
