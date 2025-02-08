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
	"project1/Common/config"
)

const (
	httpIp   = "127.0.0.1"
	httpPort = 8005
)

func main() {
	host := "192.168.88.130"
	port := uint64(8848)
	namespaceId := "public"  // 根据需要设置命名空间
	dataId := "payment"      // 你的数据 ID
	group := "DEFAULT_GROUP" // 配置组

	nacosManager, err2 := config.NewNacosConfigManager(host, port, namespaceId, dataId, group)
	if err2 != nil {
		log.Fatalf("Error creating NacosConfigManager: %v", err2)
	}

	// 启动加载配置和监听
	nacosManager.Start("paymenthttp.env")
	// 加载环境变量
	err := godotenv.Load("./paymenthttp.env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	// 初始化Hertz服务器并设置端口为8000
	h := server.Default(server.WithHostPorts("0.0.0.0:8005"))

	// 注册路由
	routes.PaymentRoutes(h)
	// 配置nacos客户端
	client := config.NewNacosClient()

	//strport := os.Getenv("PAYMENT_HTTP_PORT")
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
		ServiceName: "payment",
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
		"POST:/payment/create",
	}
	registerFlowRules(apiPaths, 10, 1000)

	url := swagger.URL("http://localhost:8005/swagger/index.html")
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
