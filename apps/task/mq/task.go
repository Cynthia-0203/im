package main

import (
	"flag"
	"fmt"
	// "sync"

	"gim/apps/task/mq/internal/config"
	"gim/apps/task/mq/internal/handler"
	"gim/apps/task/mq/internal/svc"
	// "gim/pkg/configserver"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/dev/task.yaml", "the config file")
// var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	Run(c)
	// err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
	// 	ETCDEndpoints:  "192.168.88.131:3379",
	// 	ProjectKey:     "98c6f2c2287f4c73cea3d40ae7ec3ff2",
	// 	Namespace:      "task",
	// 	Configs:        "task-mq.yaml",
	// 	ConfigFilePath: "./etc/conf",
	// 	LogLevel:       "DEBUG",
	// })).MustLoad(&c, func(bytes []byte) error {
	// 	var c config.Config
	// 	configserver.LoadFromJsonBytes(bytes, &c)
	
	// 	wg.Add(1)
	// 	go func(c config.Config) {
	// 		defer wg.Done()
	
	// 		Run(c)
	// 	}(c)
	// 	return nil
	// })
	// if err != nil {
	// 	panic(err)
	// }
	
	// wg.Add(1)
	// go func(c config.Config) {
	// 	defer wg.Done()
	
	// 	Run(c)
	// }(c)
	
	// wg.Wait()
}

func Run(c config.Config) {
	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	listen := handler.NewListen(ctx)

	serviceGroup := service.NewServiceGroup()
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	fmt.Println("Starting mqueue at ",c.ListenOn)
	serviceGroup.Start()
}