package impl

import (
	"context"
	"time"

	"github.com/infraboard/mcube/exception"

	"github.com/HAOlowkey/cmdb/apps/resource"
	"github.com/HAOlowkey/cmdb/apps/secret"
	"github.com/HAOlowkey/cmdb/apps/task"
	"github.com/HAOlowkey/cmdb/conf"
)

type SyncTaskCallback func(*task.Task)

// 通过回调更新任务状态
func (s *service) SyncTaskCallback(t *task.Task) {
	err := s.update(context.Background(), t)
	if err != nil {
		s.log.Error(err)
	}
}

func (s *service) CreatTask(ctx context.Context, req *task.CreateTaskRequst) (
	*task.Task, error) {
	// 构造 Task 对象, 并且校验了入参
	t, err := task.NewTaskFromReq(req)
	if err != nil {
		return nil, err
	}

	// 通过secret secret对象
	secret, err := s.secret.DescribeSecret(ctx, secret.NewDescribeSecretRequest(req.SecretId))
	if err != nil {
		return nil, err
	}
	t.UpdateSecretDesc(secret.ShortDesc())

	// 如果不是vsphere 需要检查region
	if !(secret.Vendor.Equal(resource.Vendor_VSPHERE) || req.ResourceType.IsIn(resource.Type_BILL)) {
		// 校验Region 参数是否传入
		if req.Region == "" {
			return nil, exception.NewBadRequest("region required")
		}
		// 只同步允许的区域
		if !secret.IsAllowRegion(req.Region) {
			return nil, exception.NewBadRequest("this secret not allow sync region %s", req.Region)
		}
	}

	// 解密secret
	err = secret.DecryptAPISecret(conf.C().App.EncryptKey)
	if err != nil {
		s.log.Warnf("decrypt api secret error, %s", err)
	}

	// 资源同步
	syncCtx, _ := context.WithTimeout(context.Background(), time.Minute*30)
	switch req.ResourceType {
	case resource.Type_HOST:
		go s.syncHost(syncCtx, secret, t, s.SyncTaskCallback)
	}

	// 记录任务
	if err := s.insert(ctx, t); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *service) QueryTask(ctx context.Context, req *task.QueryTaskRequest) (*task.TaskSet, error) {
	return nil, nil
}
func (s *service) DescribeTask(ctx context.Context, req *task.DescribeTaskRequest) (*task.Task, error) {
	return nil, nil
}
func (s *service) QueryTaskRecord(ctx context.Context, req *task.QueryTaskRecordRequest) (*task.RecordSet, error) {
	return nil, nil
}
