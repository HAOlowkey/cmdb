package task

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/xid"
)

const (
	AppName = "task"
)

var (
	validate = validator.New()
)

func (req *CreateTaskRequst) Validate() error {
	return validate.Struct(req)
}

func NewTaskFromReq(req *CreateTaskRequst) (*Task, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &Task{
		Id:           xid.New().String(),
		Type:         req.Type,
		DryRun:       req.DryRun,
		SecretId:     req.SecretId,
		ResourceType: req.ResourceType,
		Region:       req.Region,
		Timeout:      int32(req.Timeout),
	}, nil
}

func (t *Task) UpdateSecretDesc(secretDesc string) {
	t.SecretDescription = secretDesc
}

func (t *Task) Run() {
	t.StartAt = time.Now().UnixMilli()
	t.Status = Status_RUNNING
}

func (t *Task) Failed(message string) {
	t.EndAt = time.Now().UnixMilli()
	t.Status = Status_FAILED
	t.Message = message
}

func (t *Task) Completed() {
	t.EndAt = time.Now().UnixMilli()
	if t.Status != Status_FAILED {
		if t.TotalFailed == 0 {
			t.Status = Status_SUCCESS
		} else {
			t.Status = Status_WARNING
		}
	}
}
