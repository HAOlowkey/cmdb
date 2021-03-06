package cvm

import (
	"time"

	"github.com/HAOlowkey/cmdb/apps/host"
	"github.com/HAOlowkey/cmdb/apps/resource"
	"github.com/HAOlowkey/cmdb/utils"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

type CvmOperator struct {
	client *cvm.Client
	log    logger.Logger
	*AccountIdGetter
}

type AccountIdGetter struct {
	accountId string
}

func (o *AccountIdGetter) WithAccountId(aid string) {
	o.accountId = aid
}

func (o *AccountIdGetter) GetAccountId() string {
	return o.accountId
}

func NewCvmOperator(client *cvm.Client) *CvmOperator {
	return &CvmOperator{
		client:          client,
		log:             zap.L().Named("CVM OPERATOR"),
		AccountIdGetter: &AccountIdGetter{},
	}
}

func (o *CvmOperator) transferSet(items []*cvm.Instance) *host.HostSet {
	set := host.NewHostSet()

	for i := range items {
		item := o.transferOne(items[i])
		set.Add(item)
	}

	return set
}

// 描述实例的信息: https://cloud.tencent.com/document/api/213/15753#Instance
func (o *CvmOperator) transferOne(ins *cvm.Instance) *host.Host {
	h := host.NewDefaultHost()
	h.Base.Vendor = resource.Vendor_TENCENT
	h.Base.Region = o.client.GetRegion()
	h.Base.Zone = utils.PtrStrV(ins.Placement.Zone)
	h.Base.CreateAt = o.parseTime(utils.PtrStrV(ins.CreatedTime))
	h.Base.Id = utils.PtrStrV(ins.InstanceId)

	h.Information.ExpireAt = o.parseTime(utils.PtrStrV(ins.ExpiredTime))
	h.Information.Type = utils.PtrStrV(ins.InstanceType)
	h.Information.Name = utils.PtrStrV(ins.InstanceName)
	h.Information.Status = utils.PtrStrV(ins.InstanceState)
	h.Information.Tags = transferTags(ins.Tags)
	h.Information.PublicIp = utils.SlicePtrStrv(ins.PublicIpAddresses)
	h.Information.PrivateIp = utils.SlicePtrStrv(ins.PrivateIpAddresses)
	h.Information.PayType = utils.PtrStrV(ins.InstanceChargeType)
	h.Information.SyncAccount = o.GetAccountId()

	h.Describe.Cpu = utils.PtrInt64(ins.CPU)
	h.Describe.Memory = utils.PtrInt64(ins.Memory)
	h.Describe.OsName = utils.PtrStrV(ins.OsName)
	h.Describe.SerialNumber = utils.PtrStrV(ins.Uuid)
	h.Describe.ImageId = utils.PtrStrV(ins.ImageId)
	if ins.InternetAccessible != nil {
		h.Describe.InternetMaxBandwidthOut = utils.PtrInt64(ins.InternetAccessible.InternetMaxBandwidthOut)
	}
	h.Describe.KeyPairName = utils.SlicePtrStrv(ins.LoginSettings.KeyIds)
	h.Describe.SecurityGroups = utils.SlicePtrStrv(ins.SecurityGroupIds)
	return h
}

func transferTags(tags []*cvm.Tag) map[string]string {
	return nil
}

func (o *CvmOperator) parseTime(t string) int64 {
	ts, err := time.Parse("2006-01-02T15:04:05Z", t)
	if err != nil {
		o.log.Errorf("parse time %s error, %s", t, err)
		return 0
	}

	return ts.UnixMilli()
}
