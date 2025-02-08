package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"project1/Common/config"
	"strings"
)

// 发现服务
func discoverService(namingClient naming_client.INamingClient, serviceName string) ([]model.Instance, error) {
	instances, err := namingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		HealthyOnly: true,
	})
	return instances, err
}

// 创建反向代理
func reverseProxy(target *url.URL) http.Handler {
	return httputil.NewSingleHostReverseProxy(target)
}

func gatewayHandler(namingClient naming_client.INamingClient, w http.ResponseWriter, r *http.Request) {
	// 打印请求的 URI
	log.Println("Received request URI:", r.RequestURI)

	// 将 r.RequestURI 按照 "/" 分割
	parts := strings.Split(r.RequestURI, "/")
	length := len(parts)
	parts[length-1] = strings.Split(parts[length-1], "?")[0]
	// 确保有足够的部分来获取服务名称
	if len(parts) < 2 {
		http.Error(w, "Invalid request URI", http.StatusBadRequest)
		return
	}

	serviceName := parts[1] // 获取第一个部分作为服务名称
	log.Printf("Routing request to service: %s\n", serviceName)

	// 发现服务实例
	instances, err := discoverService(namingClient, serviceName)
	if err != nil || len(instances) == 0 {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	// 选择一个实例（这里简单选第一个）
	target := instances[0]
	log.Printf("Selected instance: %s:%d\n", target.Ip, target.Port)

	targetURL := fmt.Sprintf("http://%s:%d", target.Ip, target.Port)
	log.Printf("Forwarding to: %s\n", targetURL)

	targetParsedURL, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusInternalServerError)
		return
	}

	// 创建反向代理
	proxy := reverseProxy(targetParsedURL)

	// 更新请求的 URL 和请求头
	r.URL.Path = strings.Join(parts[1:], "/")     // 移除服务名称部分
	r.Host = targetParsedURL.Host                 // 设置目标主机
	r.Header.Set("X-Forwarded-Host", r.Host)      // 设置转发主机头
	r.Header.Set("X-Forwarded-For", r.RemoteAddr) // 转发原始请求的 IP

	// Log the modified request to make sure it's correct
	log.Printf("Forwarding Request: %v", r)

	// 执行代理请求并返回响应
	proxy.ServeHTTP(w, r)
}

func main() {
	// 加载配置
	godotenv.Load("config.env")
	// 创建 Nacos 客户端
	namingClient := config.NewNacosClient()

	// 设置 HTTP 服务器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		gatewayHandler(namingClient, w, r)
	})

	// 启动 HTTP 服务器
	log.Println("Starting gateway on port 9000...")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatal(err)
	}
}
