package yarx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seaung/yarx-go/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "yarx",
		Short:        "yarx 一个web漏洞扫描平台",
		Long:         "yarx 一个基于go语言开发的web在线漏洞扫描平台",
		SilenceUsage: true, // 命令出错时，不打印帮助信息。不需要打印帮助信息，设置为 true 可以保持命令出错时一眼就能看到错误信息
		RunE: func(cmd *cobra.Command, args []string) error { // 指定调用 cmd.Execute() 时，执行的 Run 函数，函数执行失败会返回错误信息
			logger.Init(initLogger()) // 初始化日志器
			defer logger.Sync()       // Sync 将缓存中的日志刷新到磁盘文件中

			return run()
		},
		Args: func(cmd *cobra.Command, args []string) error { // 这里设置命令运行时，不需要指定命令行参数
			for _, item := range args {
				if len(item) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		},
	}

	// 以下设置，使得 initConfig 函数在每个命令运行时都会被调用以读取配置
	cobra.OnInitialize(initConfig)

	// Cobra 支持持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
	cmd.PersistentFlags().StringVarP(&CFG, "config", "c", "", "请提供相关的配置文件路径")

	return cmd
}

// 程序入口启动函数
func run() error {
	// 初始化数据库连接
	if err := initStore(); err != nil {
		return err
	}

	// 设置gin模式
	gin.SetMode(viper.GetString("app.mode"))

	// 创建gin实例
	g := gin.New()

	// 中间件列表
	wms := []gin.HandlerFunc{gin.Recovery()}

	// 注册中间件
	g.Use(wms...)

	// 注册路由
	if err := installRoutes(g); err != nil {
		return err
	}

	// 创建http服务
	httpSrv := startHttpSecureServer(g)

	quit := make(chan os.Signal, 1)                                    // 定义信号
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // 监听信号
	<-quit                                                             // 这里不会阻塞

	logger.Warnw("Shutdown now server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 设置超市时间
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil { // 优雅退出程序
		logger.Warnw("Shutdown server error", err)
		return err
	}

	logger.Infow("Shutdown server success !!!")

	return nil
}

// 启动一个http server后台
func startHttpSecureServer(g *gin.Engine) *http.Server {
	httpsrv := &http.Server{
		Addr:    viper.GetString("app.port"),
		Handler: g,
	}

	go func() {
		if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorw("Server Exit !!! ", err)
		}
	}()

	return httpsrv
}
