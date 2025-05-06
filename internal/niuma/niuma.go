package niuma

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/seaung/yarx-go/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewNiuMaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "niuma",
		Short:        "niuma 后台程序",
		Long:         "niu ma 任务后台调度程序",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Init(initLogger())
			defer logger.Sync()

			return start()
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, item := range args {
				if len(item) > 0 {
					return fmt.Errorf("")
				}
			}
			return nil
		},
	}

	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().StringVar(&CFG, "config", "/home/root/.niuma/prod.yml", "请指定牛马配置文件")

	return cmd
}

// 初始化日志器
func initLogger() *logger.LoggerOptions {
	return &logger.LoggerOptions{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		OutputPaths:       viper.GetStringSlice("log.output"),
		Format:            viper.GetString("log.format"),
	}
}

func start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	instance, err := initMarchineyDemon(ctx)
	if err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit

	// 取消上下文，通知所有worker停止
	cancel()

	// 使用instance停止broker消费任务
	logger.Infow("正在停止broker消费任务...")
	instance.GetBroker().StopConsuming()

	// 等待一段时间，让worker有机会完成当前任务
	logger.Infow("等待worker完成当前任务...")
	time.Sleep(3 * time.Second)

	logger.Infow("Worker已安全退出")
	return nil
}
