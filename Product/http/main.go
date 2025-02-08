package main

import (
	"fmt"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/swagger"
	"github.com/joho/godotenv"
	"github.com/nacos-group/nacos-sdk-go/vo"
	swaggerFiles "github.com/swaggo/files"
	"http/routes"
	"log"
	"os"
	"project1/Common/config"
)

const (
	httpIp   = "127.0.0.1"
	httpPort = 8001
)

// @title HertzTest
// @version 1.0.0
// @description This is a demo using Hertz.
// @openapi: 3.0.0 // 此处加入 OpenAPI 版本字段
// @contact.name hertz-contrib
// @contact.url https://github.com/hertz-contrib
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
// @BasePath /
// @schemes http
func main() {
	// 获取进程 ID
	processID := os.Getpid()
	logFileName := fmt.Sprintf("./log/nacos-sdk-%d.log", processID)

	// 创建或打开日志文件
	file, err3 := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err3 != nil {
		log.Fatalf("Failed to open log file: %v", err3)
	}
	defer file.Close()

	// 设置日志输出
	log.SetOutput(file)

	// 记录日志
	log.Println("Service started")
	host := "192.168.88.130"
	port := uint64(8848)
	namespaceId := "public"  // 根据需要设置命名空间
	dataId := "product"      // 你的数据 ID
	group := "DEFAULT_GROUP" // 配置组

	nacosManager, err2 := config.NewNacosConfigManager(host, port, namespaceId, dataId, group)
	if err2 != nil {
		log.Fatalf("Error creating NacosConfigManager: %v", err2)
	}

	// 启动加载配置和监听
	nacosManager.Start("producthttp.env")
	// 加载环境变量
	err := godotenv.Load("./producthttp.env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	// 初始化Hertz服务器并设置端口为8000
	h := server.Default(server.WithHostPorts("0.0.0.0:8001"))

	// 注册路由
	routes.ProductRoutes(h)
	// 配置nacos客户端
	client := config.NewNacosClient()

	//strport := os.Getenv("PRODUCT_PORT")
	//intport, err := strconv.Atoi(strport)
	//if err != nil {
	//	fmt.Println("启动端口失败")
	//}
	// 注册服务到Nacos
	ip := httpIp    // 服务器IP，根据实际情况调整
	port = httpPort // Hertz服务器端口为8000

	_, err = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: "product",
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})

	if err != nil {
		log.Fatalf("register service instance failed: %v", err)
	}

	// 初始化Sentinel
	err = api.InitDefault()
	if err != nil {
		log.Fatalf("Sentinel init failed: %v", err)
	}
	// 批量注册限流规则
	apiPaths := []string{
		"GET:/product/getList",
		"GET:/product/getById/:productId",
		"GET:/product/getByName",
	}
	registerFlowRules(apiPaths, 10, 1000)

	url := swagger.URL("http://localhost:8001/swagger/index.html")
	h.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))
	// 启动Hertz服务器
	if err := h.Run(); err != nil {
		panic(err)
	}
}

// registerFlowRules 批量注册限流规则
func registerFlowRules(resources []string, threshold float64, intervalMs uint32) {
	rules := make([]*flow.Rule, 0, len(resources))
	for _, res := range resources {
		rules = append(rules, &flow.Rule{
			Resource:               res,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              threshold,
			StatIntervalInMs:       intervalMs,
		})
	}
	if _, err := flow.LoadRules(rules); err != nil {
		log.Fatalf("load rules failed: %v", err)
	}
}
