package connectivity_test

import (
	"fmt"
	"os"
	"testing"

	"gitee.com/go-course/cmdb/provider/txyun/connectivity"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
)

var (
	client *connectivity.TencentCloudClient
)

func TestCvmClient(t *testing.T) {
	cvmConn := client.CvmClient()
	fmt.Print(cvmConn)
}

func init() {
	var secretID, secretKey string
	if secretID = os.Getenv("TX_CLOUD_SECRET_ID"); secretID == "" {
		panic("empty TX_CLOUD_SECRET_ID")
	}

	if secretKey = os.Getenv("TX_CLOUD_SECRET_KEY"); secretKey == "" {
		panic("empty TX_CLOUD_SECRET_KEY")
	}

	client = connectivity.NewTencentCloudClient(secretID, secretKey, regions.Shanghai)
}
