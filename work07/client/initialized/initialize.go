package initialized

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"zg5/work/work07/client/api_consul"
	"zg5/work/work07/client/api_global"
	"zg5/work/work07/client/routers"
	"zg5/work/work07/client/service"
	"zg5/work/work07/server/proto/server"
)

func init() {
	InitZapLog()
	InitViper()
	InitNacos()
	initRouters()

	service.InitMysql()
	service.InitRedis()
}

func InitViper() {
	v := viper.New()
	v.SetConfigFile("./config/api.yaml")
	err := v.ReadInConfig()
	if err != nil {
		zap.S().Info("Couldn't read config file'")
	}
	v.Unmarshal(&api_global.NacosConfig)
	zap.S().Info(api_global.NacosConfig, "fdslkfnuidsvodsnvd")
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
		constant.WithNamespaceId(api_global.NacosConfig.NamespaceId), //When namespace is public, fill in the blank string here.
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{ //TODO: 22222
		{
			IpAddr: api_global.NacosConfig.Host,
			Port:   uint64(api_global.NacosConfig.Port),
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
		DataId: api_global.NacosConfig.DataId,
		Group:  api_global.NacosConfig.Group,
	})
	yaml.Unmarshal([]byte(content), &api_global.ClientConfig)
}

func InitConsul() api_consul.RegistryClient {
	api_global.ClientConfig.Consul.Id = uuid.NewString()
	var (
		err error
	)

	api_global.ConsulClient, err = api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	//consul服务发现
	consulClinet := api_consul.NewConsulClient(api_global.ClientConfig.ApiConf.Host, api_global.ClientConfig.ApiConf.Port)
	i := consulClinet.AgentHealthServiceByName(api_global.ClientConfig.Consul.UserConsulName)
	//拨号
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", i[0].Service.Address, i[0].Service.Port), grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	api_global.ServerClient = server.NewServerClient(conn)
	//注册
	err = consulClinet.RegisterConsul(api_global.ClientConfig.Consul.Id, api_global.ClientConfig.Consul.Name,
		api_global.ClientConfig.Consul.Tags)
	if err != nil {
		zap.S().Info("注册失败")
		panic(err)
	}
	return consulClinet
}

func InitGrpc() {
	consulClient := InitConsul()
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", api_global.ClientConfig.ApiConf.Host, api_global.ClientConfig.ApiConf.Port),
		Handler: api_global.Router,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	//优雅的关闭
	SigChan := make(chan os.Signal, 1)
	signal.Notify(SigChan, syscall.SIGINT, syscall.SIGTERM)
	<-SigChan
	err := consulClient.DeregisterConsulById(api_global.ClientConfig.Consul.Id)
	if err != nil {
		panic(err)
	}
}

func initRouters() {
	api_global.Router = gin.New()
	Group := api_global.Router.Group("vv1")
	routers.InitRouter(Group)
}
