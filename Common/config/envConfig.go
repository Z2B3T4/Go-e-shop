package config

import (
	"bytes"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
)

type NacosConfigManager struct {
	client      config_client.IConfigClient
	dataId      string
	group       string
	configMap   map[string]interface{}
	configMutex sync.Mutex
}

// NewNacosConfigManager - 创建 Nacos 配置管理器
func NewNacosConfigManager(host string, port uint64, namespaceId, dataId, group string) (*NacosConfigManager, error) {
	serverConfig := []constant.ServerConfig{
		{
			IpAddr: host,
			Port:   port,
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         namespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./log",
		CacheDir:            "./cache",
	}

	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: serverConfig,
	})

	if err != nil {
		return nil, fmt.Errorf("create Nacos config client error: %v", err)
	}

	return &NacosConfigManager{
		client:    client,
		dataId:    dataId,
		group:     group,
		configMap: make(map[string]interface{}),
	}, nil
}

// Start - 开始加载配置和监听变化
func (ncm *NacosConfigManager) Start(filename string) {
	ncm.loadConfig(filename)

	// 在 goroutine 中监听配置变化
	go ncm.listenConfig(filename)
}

// loadConfig - 加载配置
func (ncm *NacosConfigManager) loadConfig(filename string) {
	ncm.configMutex.Lock()
	defer ncm.configMutex.Unlock()

	content, err := ncm.client.GetConfig(vo.ConfigParam{
		DataId: ncm.dataId,
		Group:  ncm.group,
	})
	if err != nil {
		log.Fatalf("Failed to get config from Nacos: %v", err)
	}

	// 将 YAML 内容解析为 map
	err = yaml.Unmarshal([]byte(content), &ncm.configMap)
	if err != nil {
		log.Fatalf("Error parsing YAML content: %v", err)
	}

	// 转换为 .env 格式
	envContent := ncm.convertMapToEnv()

	// 写入 .env 文件
	err = ncm.writeEnvToFile(filename, envContent)
	if err != nil {
		log.Fatalf("Error writing to .env file: %v", err)
	}

	log.Println("Configuration successfully written to .env file.")
}

// listenConfig - 监听配置变化
func (ncm *NacosConfigManager) listenConfig(filename string) {
	err := ncm.client.ListenConfig(vo.ConfigParam{
		DataId: ncm.dataId,
		Group:  ncm.group,
		OnChange: func(namespace, group, dataId, data string) {
			log.Printf("Config has changed: %s", data)
			ncm.loadConfig(filename)
		},
	})

	if err != nil {
		log.Fatalf("Error adding listener: %v", err)
	}
}

// convertMapToEnv - 将 map 转换为 .env 格式字符串
func (ncm *NacosConfigManager) convertMapToEnv() string {
	var envBuffer bytes.Buffer
	for key, value := range ncm.configMap {
		envBuffer.WriteString(fmt.Sprintf("%s=%v\n", key, value))
	}
	return envBuffer.String()
}

// writeEnvToFile - 将内容写入 .env 文件
func (ncm *NacosConfigManager) writeEnvToFile(filename, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}
