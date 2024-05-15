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
	"net"
	"os"
	"os/signal"
	"syscall"
	"zg5/work/work07/server/consule"
	"zg5/work/work07/server/global"
	"zg5/work/work07/server/logic"
	"zg5/work/work07/server/proto/server"
)

func init() {
	InitViper()
	InitNacos()
	InitZapLog()
}

func InitViper() {
	v := viper.New()
	v.SetConfigFile("./config/config.yaml")
	err := v.ReadInConfig()
	if err != nil {
		zap.S().Info("Couldn't read config file")
	}
	v.Unmarshal(&global.NacosConfig)
}

func InitZapLog() {
	logger := zap.NewDevelopmentConfig()
	logger.OutputPaths = []string{
		"./logger/logger.log",
		"stderr",
		"stdout",
	}
	log, _ := logger.Build()
	zap.ReplaceGlobals(log)
}

func InitNacos() {
	//创建 clientConfig 的另一种方式
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(global.NacosConfig.NamespaceId), //When namespace is public, fill in the blank string here.
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{ //TODO: 22222
		{
			IpAddr: global.NacosConfig.Host,
			Port:   uint64(global.NacosConfig.Port),
		},
	}
	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err := clients.NewConfigClient( //TODO: 33333
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic("配置信息交互失败！")
	}
	content, err := configClient.GetConfig(vo.ConfigParam{ //TODO:44444
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	yaml.Unmarshal([]byte(content), &global.ServerConfig)
}

func Consul() consule.RegistryClient {
	global.ServerConfig.Consul.Id = uuid.NewString()
	var (
		err      error
		BaseHost string
		BasePort int
	)
	global.ConsulClient, err = api.NewClient(api.DefaultConfig())
	if err != nil {
		panic("api_consul: failed to create Consul client")
	}
	consulClient := consule.NewConsulClient(global.ServerConfig.Grpc.Host, global.ServerConfig.Grpc.Port)
	serviceList, _ := consulClient.FilterConsulByName(global.ServerConfig.Consul.Name)
	for _, v := range serviceList {
		BaseHost = v.Address
		BasePort = v.Port
	}
	if BaseHost != "" || BasePort != 0 {
		zap.S().Warn("服务已注册")
	}
	err = consulClient.RegisterConsul(global.ServerConfig.Consul.Id, global.ServerConfig.Consul.Name, global.ServerConfig.Consul.Tags)
	if err != nil {
		panic(err)
	}
	return consulClient
}

func InitGrpc() {
	GrpcServer := grpc.NewServer()
	Server := &logic.InitService{}
	server.RegisterServerServer(GrpcServer, Server)
	grpc_health_v1.RegisterHealthServer(GrpcServer, health.NewServer())
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.ServerConfig.Grpc.Host, global.ServerConfig.Grpc.Port))
	if err != nil {
		panic("gRpc服务监听失败！")
	}
	consulClient := Consul()
	go func() {
		err = GrpcServer.Serve(listen)
		panic(err)
	}()
	zap.S().Info("已启动rpc服务正在监听中...", global.ServerConfig.Grpc.Host+":", global.ServerConfig.Grpc.Port)
	SignChan := make(chan os.Signal, 1)
	signal.Notify(SignChan, syscall.SIGINT, syscall.SIGTERM)
	<-SignChan
	err = consulClient.DeregisterConsulById(global.ServerConfig.Consul.Id)
	if err != nil {
		panic(err)
	}
	GrpcServer.Stop()
}
