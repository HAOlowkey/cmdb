package all

import (
	// 注册所有HTTP服务模块, 暴露给框架HTTP服务器加载
	_ "github.com/HAOlowkey/cmdb/apps/book/http"
	_ "github.com/HAOlowkey/cmdb/apps/host/http"
	_ "github.com/HAOlowkey/cmdb/apps/resource/http"
	_ "github.com/HAOlowkey/cmdb/apps/secret/http"
	_ "github.com/HAOlowkey/cmdb/apps/task/http"
)
