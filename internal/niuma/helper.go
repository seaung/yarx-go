package niuma

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	server "github.com/RichardKnop/machinery/v2"
	redisbackend "github.com/RichardKnop/machinery/v2/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v2/brokers/redis"
	eagerlock "github.com/RichardKnop/machinery/v2/locks/eager"
	"github.com/seaung/yarx-go/pkg/machinery"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	CFG string
)

const (
	userHomeDir     = ".niuma"
	defaultFileName = "niuma.yml"
)

// 初始化machinery配置信息文件
func initConfig() {
	if CFG != "" {
		viper.SetConfigFile(CFG)
	} else {
		home, err := os.UserHomeDir()

		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(home, userHomeDir))
		viper.AddConfigPath(".")

		viper.SetConfigType("yaml")
		viper.SetConfigName(defaultFileName)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

// 初始化machinery后台
func initMarchineyDemon(ctx context.Context) (*server.Server, error) {
	opts := &machinery.MachineryOptions{
		Broker:       viper.GetString("app.redis.broker"),
		Backend:      viper.GetString("app.redis.backend"),
		DefaultQueue: viper.GetString("app.redis.default_queue"),
	}

	cfg, _ := machinery.NewMachineryConfig(opts)

	broker := redisbroker.New(cfg, viper.GetString("app.redis.broker"), "", "", 0)
	backend := redisbackend.New(cfg, viper.GetString("app.redis.backend"), "", "", 1)
	lock := eagerlock.New()

	instance := server.NewServer(cfg, broker, backend, lock)

	// 注册任务
	_ = registerTasks(instance)

	wk := instance.NewWorker(viper.GetString("app.redis.default_queue"), viper.GetInt("app.redis.worker_num"))

	go func() {
		go func() {
			<-ctx.Done()
		}()
		wk.Launch()
	}()

	return instance, nil
}

// 注册任务
func registerTasks(srv *server.Server) error {
	return srv.RegisterTasks(map[string]any{
		"scanner": func() { fmt.Println("scanner ...") },
	})
}
