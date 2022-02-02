package cvm

import (
	"github.com/HAOlowkey/cmdb/apps/host"
	"github.com/infraboard/mcube/flowcontrol/tokenbucket"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func newPager(size int, operator *CvmOperator, rate int) *pager {
	req := cvm.NewDescribeInstancesRequest()
	req.Limit = common.Int64Ptr(int64(size))
	return &pager{
		number:   1,
		size:     size,
		operator: operator,
		req:      req,
		log:      zap.L().Named("CVM Pager"),
		tb:       tokenbucket.NewBucketWithRate(1/float64(rate), 1),
	}
}

type pager struct {
	number   int
	size     int
	total    int64
	operator *CvmOperator
	req      *cvm.DescribeInstancesRequest
	log      logger.Logger
	tb       *tokenbucket.Bucket
}

func (p *pager) Next() *host.PagerResult {
	result := host.NewPagerResult()

	resp, err := p.operator.Query(p.nextReq())
	if err != nil {
		result.Err = err
		return result
	}
	p.total = resp.Total
	p.log.Debugf("get %d hosts", len(resp.Items))

	result.Data = resp
	result.HasNext = p.hasNext()

	p.number++
	return result
}

func (p *pager) nextReq() *cvm.DescribeInstancesRequest {
	p.tb.Wait(1)
	p.log.Debugf("请求第%d页数据", p.number)
	p.req.Offset = common.Int64Ptr(p.offset())
	return p.req
}

func (p *pager) hasNext() bool {
	return p.number*p.size < int(p.total)
}

func (p *pager) offset() int64 {
	return int64((p.number - 1) * p.size)
}
