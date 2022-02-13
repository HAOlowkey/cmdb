package task

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/infraboard/mcube/http/request"
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

func NewSyncFailedRecord(taskId, instanceId, instanceName, message string) *Record {
	return &Record{
		TaskId:     taskId,
		InstanceId: instanceId,
		Name:       instanceName,
		Message:    message,
		IsSuccess:  false,
	}
}

func NewSyncSucceedRecord(taskId, instanceId, instanceName string) *Record {
	return &Record{
		TaskId:     taskId,
		InstanceId: instanceId,
		Name:       instanceName,
		IsSuccess:  true,
	}
}

func (t *Task) AddDetail(d *Record) {
	if d.IsSuccess {
		t.TotalSucceed++
	} else {
		t.TotalFailed++
	}
}

func NewTaskSet() *TaskSet {
	return &TaskSet{
		Items: []*Task{},
	}
}

func (r *TaskSet) Add(item *Task) {
	r.Items = append(r.Items, item)
}

func NewDefaultTask() *Task {
	return &Task{}
}

func NewRecordSet() *RecordSet {
	return &RecordSet{
		Items: []*Record{},
	}
}

func (r *RecordSet) Add(item *Record) {
	r.Items = append(r.Items, item)
}

func NewDefaultTaskRecord() *Record {
	return &Record{}
}

func NewCreateTaskRequst() *CreateTaskRequst {
	return &CreateTaskRequst{}
}

func NewQueryTaskRequestFromHTTP(r *http.Request) *QueryTaskRequest {
	qs := r.URL.Query()

	kw := qs.Get("keywords")

	return &QueryTaskRequest{
		Page:     request.NewPageRequestFromHTTP(r),
		Keywords: kw,
	}
}

func NewQueryTaskRecordRequest(id string) *QueryTaskRecordRequest {
	return &QueryTaskRecordRequest{
		TaskId: id,
	}
}

func NewDescribeTaskRequestWithId(id string) *DescribeTaskRequest {
	return &DescribeTaskRequest{
		Id: id,
	}
}
