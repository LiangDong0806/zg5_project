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
	"zg5/work/work09/server/consul"
	"zg5/work/work09/server/global"
	"zg5/work/work09/server/logic"
	"zg5/work/work09/server/proto/server"
)

func init() {
	initViper()
	initNacos()
	initZapLog()
}

func initViper() {
	v := viper.New()
	v.SetConfigFile("./config/rpc.yaml")
	err := v.ReadInConfig()
	if err != nil {
	}
	v.Unmarshal(&global.NacosConf)
}

func initZapLog() {
	log := zap.NewDevelopmentConfig()
	log.OutputPaths = []string{
		"./logger/message.log",
		"stderr",
		"stdout",
	}
	logs, _ := log.Build()
	zap.ReplaceGlobals(logs)
}

func initNacos() {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{ //TODO: 11111
		NamespaceId:         global.NacosConf.NamespaceId, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{ //TODO: 22222
		{
			IpAddr:      global.NacosConf.Host,
			ContextPath: "/nacos",
			Port:        uint64(global.NacosConf.Port),
			Scheme:      "http",
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
		//panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{ //TODO:44444
		DataId: global.NacosConf.DataId,
		Group:  global.NacosConf.Group,
	})
	yaml.Unmarshal([]byte(content), &global.RpcALLConf)
	zap.S().Info("client configuration：", &content)
}

func initConsul() consul.RegistryClient {
	global.RpcALLConf.Consul.Id = uuid.NewString()
	var (
		err      error
		BaseHost string
		BasePort int
	)
	global.ConsulClient, err = api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	consulClient := consul.NewConsulClient(global.RpcALLConf.Grpc.Host, global.RpcALLConf.Grpc.Port)
	serviceList, err := consulClient.FilterConsulByName(global.RpcALLConf.Consul.Name)
	if err != nil {
		panic(err)
	}
	for _, s := range serviceList {
		BaseHost = s.Address
		BasePort = s.Port
	}
	if BaseHost != "" || BasePort != 0 {
		zap.S().Info("theServiceIsRegistered")
	}
	err = consulClient.RegisterConsul(global.RpcALLConf.Consul.Id, global.RpcALLConf.Consul.Name, global.RpcALLConf.Consul.Tags)
	if err != nil {
		panic(err)
	}
	return consulClient
}

func InitGrpc() {
	GrpcServer := grpc.NewServer()

	Server := &logic.RpcServer{}
	server.RegisterServerServer(GrpcServer, Server)
	grpc_health_v1.RegisterHealthServer(GrpcServer, health.NewServer())
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d",
		global.RpcALLConf.Grpc.Host, global.RpcALLConf.Grpc.Port))
	if err != nil {
		panic(err)
	}
	consulClient := initConsul()
	go func() {
		err = GrpcServer.Serve(listen)
		if err != nil {
			panic(err)
		}
	}()
	zap.S().Info("127.0.0.1:8080...")
	SignChan := make(chan os.Signal, 1)
	signal.Notify(SignChan, syscall.SIGINT, syscall.SIGTERM)
	<-SignChan
	err = consulClient.DeregisterConsulByID(global.RpcALLConf.Consul.Id)
	if err != nil {
		panic(err)
	}
	GrpcServer.Stop()
}
