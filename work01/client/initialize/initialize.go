package initialize

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"zg5/work01/client/consul"
	"zg5/work01/client/global"
	server "zg5/work01/server/proto"
)

func init() {
	initViper()
	initZapLog()
	initNacos()
}

func initViper() {
	v := viper.New()
	v.SetConfigFile("./config/api.yaml")
	err := v.ReadInConfig()
	if err != nil {
		log.Println("Error reading config file")
	}
	v.Unmarshal(&global.NacosConf)
}

func initZapLog() {
	log := zap.NewDevelopmentConfig()
	log.OutputPaths = []string{
		"./logger/logger.log",
		"stderr",
		"stdout",
	}
	logs, _ := log.Build()
	zap.ReplaceGlobals(logs)
}

func initNacos() {
	// 创建clientConfig的另一种方式
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(global.NacosConf.NamespaceId), //当namespace是public时，此处填空字符串。
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)
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
		panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{ //TODO:44444
		DataId: global.NacosConf.DataId,
		Group:  global.NacosConf.Group,
	})
	yaml.Unmarshal([]byte(content), &global.ApiConf)
}

func initConsul() consul.RegistryClient {
	global.ApiConf.Consul.Id = uuid.NewString()
	var err error
	global.ConsulClient, err = api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	consulClient := consul.NewConsulClient(global.ApiConf.ApiConf.Host, global.ApiConf.ApiConf.Port)
	car := consulClient.AgentHealthServiceByName(global.ApiConf.Consul.UserConsulName)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", car[0].Service.Address, car[0].Service.Port), grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	global.ServerClient = server.NewServerClient(conn)
	err = consulClient.RegisterConsul(global.ApiConf.Consul.Id, global.ApiConf.Consul.Name, global.ApiConf.Consul.Tags)
	if err != nil {
		panic(err)
	}
	return consulClient
}

func InitGrpc() {
	consulClient := initConsul()
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", global.ApiConf.ApiConf.Host, global.ApiConf.ApiConf.Port),
		Handler: global.Engine,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	zap.S().Info("端口监听中...", global.ApiConf.ApiConf.Host+":", global.ApiConf.ApiConf.Port)
	SignChan := make(chan os.Signal, 1)
	signal.Notify(SignChan, syscall.SIGINT, syscall.SIGTERM)
	<-SignChan
	err := consulClient.DeregisterConsulById(global.ApiConf.Consul.Id)
	if err != nil {
		panic(err)
	}
}
