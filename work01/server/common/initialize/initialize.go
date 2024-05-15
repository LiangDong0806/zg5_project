package initialize

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"gopkg.in/yaml.v2"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"zg5/work01/server/common/global"
	"zg5/work01/server/logic"
	"zg5/work01/server/models"
	server "zg5/work01/server/proto"
	"zg5/work01/server/vender/consul"
)

func init() {
	InitViper()
	InitZap()
	InitNacos()
	models.InitElastic()
}

// TODO: 第一步 Viper配置信息
func InitViper() {
	v := viper.New()
	v.SetConfigFile("./etc/dev.yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	v.Unmarshal(&global.NacosConfig)
}

// TODO: 第二步 Zap日志
func InitZap() {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{
		"./logger/server.logger",
		"stderr",
		"stdout",
	}
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)
}

// TODO: 第三步 Nacos动态配置
func InitNacos() {
	// 创建clientConfig
	log.Println(global.NacosConfig, "[][]][]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]")
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/logger",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   uint64(global.NacosConfig.Port),
		},
	}
	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		panic(err)
	}
	log.Println(content, "...................")
	yaml.Unmarshal([]byte(content), &global.ServerConfig)
	zap.S().Info(global.ServerConfig.Consul, global.ServerConfig.MySQL, "11111111111111111111111")
}

// TODO: 第四步 初始化 Consul 客户端
func InitConsul() consule.RegistryClient {
	global.ServerConfig.Consul.Id = uuid.New().String()
	var (
		err      error
		BaseHost string
		BasePort int
	)

	global.ConsulClient, err = api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	//过滤
	consulClient := consule.NewConsulClient(global.ServerConfig.Grpc.Host, global.ServerConfig.Grpc.Port)
	serviceList, err := consulClient.FilterConsulByName(global.ServerConfig.Consul.Name)
	if err != nil {
		panic(err)
	}
	for _, v := range serviceList {
		BaseHost = v.Address
		BasePort = v.Port
	}
	if BaseHost != "" || BasePort != 0 {
		zap.S().Info("服务已注册！")
	}
	//注册
	err = consulClient.RegisterConsul(global.ServerConfig.Consul.Id, global.ServerConfig.Consul.Name, global.ServerConfig.Consul.Tags)
	if err != nil {
		zap.S().Info("注册失败")
		panic(err)
	}
	return consulClient
}

// TODO: 第五步 初始化Grpc服务
func InitGrpc() {
	GrpcServer := grpc.NewServer()

	Server := &logic.ServerRpc{}
	server.RegisterServerServer(GrpcServer, Server)
	grpc_health_v1.RegisterHealthServer(GrpcServer, health.NewServer())
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Grpc.Host, global.ServerConfig.Grpc.Port))
	if err != nil {
		panic(err)
	}
	consulClient := InitConsul()
	go func() {
		err = GrpcServer.Serve(listen)
		if err != nil {
			panic(err)
		}
	}()
	zap.S().Info("保姆究极状态下已启动...")
	SignChan := make(chan os.Signal, 1)
	signal.Notify(SignChan, syscall.SIGINT, syscall.SIGTERM)
	<-SignChan
	err = consulClient.DeregisterConsul(global.ServerConfig.Consul.Id)
	if err != nil {
		panic(err)
	}
	GrpcServer.Stop()
}

//	func InitConsul() {
//		client, _ := api.NewClient(api.DefaultConfig())
//		result, _ := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%v"`, "Server"))
//		var BaseHost = ""
//		for _, v := range result {
//			BaseHost = v.Address
//		}
//		if BaseHost == "" {
//			err := client.Agent().ServiceRegister(&api.AgentServiceRegistration{
//				ID:      uuid.New().String(),
//				Name:    "Server",
//				Tags:    nil,
//				Port:    8001,
//				Address: "127.0.0.1",
//			})
//			if err != nil {
//				zap.S().Info("consule 注册失败")
//			} else {
//				zap.S().Info("consule 已注册")
//			}
//		} else {
//			zap.S().Info("consule 注册成功")
//		}
//
// }
