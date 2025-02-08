package config

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"log"
	"os"
	"strconv"
)

func NewNacosClient() naming_client.INamingClient {
	// 这个是谁调用就用谁的配置文件，比如User调用借用User的http下的配置文件
	fmt.Println("NACOS_PORT", os.Getenv("NACOS_PORT"))
	fmt.Println("NACOS_HOST", os.Getenv("NACOS_HOST"))
	num, err2 := strconv.Atoi(os.Getenv("NACOS_PORT"))
	fmt.Println("num:", num)
	if err2 != nil {
		log.Fatal("Nacos Port Error")
	}
	// 配置Nacos客户端
	serverConfig := []constant.ServerConfig{
		{
			IpAddr: os.Getenv("NACOS_HOST"),
			Port:   uint64(num),
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         "public", // 可选，根据你的命名空间来设置
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./log",
		CacheDir:            "./cache",
	}

	// 创建Nacos客户端
	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: serverConfig,
	})

	if err != nil {
		log.Fatalf("create Nacos client error: %v", err)
	}
	return client
}
