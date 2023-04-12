package yarxgo

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
	"github.com/seaung/yarx-go/internal/pkg/constants"
	"github.com/seaung/yarx-go/internal/pkg/log"
	"github.com/seaung/yarx-go/pkg/token"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func NewYarxgoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "yarx-go",
		Short:        "yarx-go 一个golang开发的web项目",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Init(logOptions())
			defer log.Sync()

			return run()
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the yarx-go configuration file. Empty string for no configuration file")

	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return cmd
}

func run() error {
	if err := initStore(); err != nil {
		return err
	}

	token.Init(viper.GetString("jwt-secret"), constants.XUsernameKey)
	gin.SetMode(viper.GetString("runmode"))

	g := gin.New()

	if err := initRouters(g); err != nil {
		return err
	}

	httpServ := runHttpServer(g)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServ.Shutdown(ctx); err != nil {
		log.Errorw("Insecure Server forced to shutdown", err, err)
		return err
	}

	return nil
}

func runHttpServer(g *gin.Engine) *http.Server {
	httpsrv := &http.Server{
		Addr:    viper.GetString("addr"),
		Handler: g,
	}

	log.Info("Start to listening the incomming request on http address", "addr", viper.GetString("addr"))

	go func() {
		if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()

	return httpsrv
}
