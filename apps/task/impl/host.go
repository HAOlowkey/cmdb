package impl

import (
	"context"
	"fmt"

	"github.com/HAOlowkey/cmdb/apps/host"
	"github.com/HAOlowkey/cmdb/apps/resource"
	"github.com/HAOlowkey/cmdb/apps/secret"
	"github.com/HAOlowkey/cmdb/apps/task"
	txConn "github.com/HAOlowkey/cmdb/provider/txyun/connectivity"
	cvmOp "github.com/HAOlowkey/cmdb/provider/txyun/cvm"
)

func (s *service) syncHost(ctx context.Context, secret *secret.Secret, t *task.Task, cb SyncTaskCallback) {
	var (
		pager host.Pager
	)

	// 处理任务状态, 把状态标记为Running
	t.Run()
	defer func() {
		// 因为sync 是异步的，跑在goroutine中
		// 如果该goroutine 有panic, goroutine外部是无法捕捉的
		// 使用 recover 来捕捉 当前goroutine抛出的panic
		if err := recover(); err != nil {
			s.log.Errorf("pannic, %v", err)
			t.Failed(fmt.Sprintf("pannic, %v", err))
		} else {
			t.Completed()
		}
		cb(t)
	}()

	// 获取一个 查询 CVM Pagger
	switch secret.Vendor {
	case resource.Vendor_TENCENT:
		s.log.Debugf("sync txyun cvm ...")
		client := txConn.NewTencentCloudClient(secret.ApiKey, secret.ApiSecret, t.Region)
		// 校验了客户端是否 正常使用，补充AccountId
		// 每一个Resource, 都有一个syncAccount
		if err := client.Check(); err != nil {
			t.Failed(err.Error())
			return
		}
		s.log.Debugf("check account %s", client.AccountID())
		operater := cvmOp.NewOperator(client.CvmClient())
		operater.WithAccountId(client.AccountID())
		req := cvmOp.NewPageQueryRequest(int(secret.RequestRate))
		pager = operater.PageQuery(req)
	default:
		t.Failed(fmt.Sprintf("unsuport vendor %s", secret.Vendor))
		return
	}

	// 分页查询数据
	if pager != nil {
		hasNext := true
		for hasNext {
			p := pager.Next()
			hasNext = p.HasNext

			if p.Err != nil {
				t.Failed(fmt.Sprintf("sync error, %s", p.Err))
				return
			}

			// 调用host服务保持数据
			for i := range p.Data.Items {
				target := p.Data.Items[i]
				// 补充管理信息
				target.Base.SecretId = secret.Id
				s.SyncHost(ctx, target, t)
			}
		}
	}
}

func (s *service) SyncHost(ctx context.Context, ins *host.Host, t *task.Task) {
	// 往host 服务同步了一个资源
	h, err := s.host.SyncHost(ctx, ins)

	// 通过详情使用Record记录
	var detail *task.Record
	if err != nil {
		s.log.Warnf("save host error, %s", err)
		detail = task.NewSyncFailedRecord(t.Id, ins.Base.Id, ins.Information.Name, err.Error())
	} else {
		s.log.Debugf("save host %s to db", h.ShortDesc())
		detail = task.NewSyncSucceedRecord(t.Id, ins.Base.Id, ins.Information.Name)
	}

	t.AddDetail(detail)
	if err := s.insertTaskDetail(ctx, detail); err != nil {
		s.log.Errorf("update detail error, %s", err)
	}
}
