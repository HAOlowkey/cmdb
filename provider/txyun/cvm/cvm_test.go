package cvm_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/HAOlowkey/cmdb/provider/txyun/connectivity"
	"github.com/HAOlowkey/cmdb/provider/txyun/cvm"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/stretchr/testify/assert"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	txcvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

var operator *cvm.CvmOperator

func TestQueryCvmHost(t *testing.T) {
	should := assert.New(t)
	req := txcvm.NewDescribeInstancesRequest()
	set, err := operator.Query(req)
	should.NoError(err)
	fmt.Println(set)

}

func TestPageQueryCVMInstances(t *testing.T) {
	should := assert.New(t)

	pg := operator.PageQuery(cvm.NewPageQueryRequest(5, 5))
	hasNext := true
	for hasNext {
		ps := pg.Next()
		should.NoError(ps.Err)
		fmt.Println(ps.Data)
		hasNext = ps.HasNext
	}

}

func init() {
	var secretID, secretKey string
	if secretID = os.Getenv("TX_CLOUD_SECRET_ID"); secretID == "" {
		panic("empty TX_CLOUD_SECRET_ID")
	}

	if secretKey = os.Getenv("TX_CLOUD_SECRET_KEY"); secretKey == "" {
		panic("empty TX_CLOUD_SECRET_KEY")
	}

	client := connectivity.NewTencentCloudClient(secretID, secretKey, regions.Shanghai)
	operator = cvm.NewCvmOperator(client.CvmClient())

	zap.DevelopmentSetup()
}
