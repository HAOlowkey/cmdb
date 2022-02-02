package cvm

import (
	"gitee.com/go-course/cmdb/utils"
	"github.com/HAOlowkey/cmdb/apps/host"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func (o *CvmOperator) Query(req *cvm.DescribeInstancesRequest) (*host.HostSet, error) {
	resp, err := o.client.DescribeInstances(req)
	if err != nil {
		return nil, err
	}

	set := o.transferSet(resp.Response.InstanceSet)
	set.Total = utils.PtrInt64(resp.Response.TotalCount)

	return set, nil
}

func NewPageQueryRequest(size, reqPerSecond int) *PageQueryRequest {
	return &PageQueryRequest{
		ReqPerSecond: reqPerSecond,
		size:         size,
	}
}

type PageQueryRequest struct {
	ReqPerSecond int
	size         int
}

func (o *CvmOperator) PageQuery(req *PageQueryRequest) host.Pager {
	return newPager(req.size, o, req.ReqPerSecond)
}
