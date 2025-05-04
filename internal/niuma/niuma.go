package niuma

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewNiuMaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "niuma",
		Short:        "niuma 后台程序",
		Long:         "niu ma 任务后台调度程序",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
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

	cobra.OnInitialize()

	cmd.PersistentFlags().StringVar(&CFG, "config", "/home/root/.niuma/prod.yml", "请指定牛马配置文件")

	return cmd
}

func start() error {
	return nil
}
