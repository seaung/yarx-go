package yarx

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/seaung/yarx-go/internal/yarx/store"
	"github.com/seaung/yarx-go/pkg/db"
	"github.com/seaung/yarx-go/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	userHomeDir       = ".yarx"
	defaultConfigName = "dev.yml"
)

var (
	CFG string
)

// 初始化配置文件
func initConfig() {
	if CFG != "" {
		// 从命令行选项中指定配置文件，进行配置的读取
		viper.SetConfigFile(CFG)
	} else {
		// 否则通过下面的设置读取配置文件
		home, err := os.UserHomeDir() // 读取用户当前的家目录
		cobra.CheckErr(err)           // 如果获取失败就调用CheckErr方法打印错误信息并退出程序

		viper.AddConfigPath(filepath.Join(home, userHomeDir)) // 将家目录下的.yarx所在路径添加进配置文件的搜索路径中

		viper.AddConfigPath(".")               // 将当前路径也添加到搜索路径里面
		viper.SetConfigFile("yaml")            // 设置配置文件的类型为yaml格式
		viper.SetConfigName(defaultConfigName) // 设置配置文件的默认名称
	}

	viper.AutomaticEnv()       // 设置自动读取环境变量
	viper.SetEnvPrefix("YARX") // 设置读取的环境变量前缀为YARX

	// 将 viper.Get(key) key 字符串中 '.' 和 '-' 替换为 '_'
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil { // 读取配置文件
		log.Fatal(err)
	}
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

// 初始化数据库
func initStore() error {
	dbOptions := &db.SQLOptions{
		Host:    viper.GetString("db.host"),
		Name:    viper.GetString("db.name"),
		Pass:    viper.GetString("db.pass"),
		DB:      viper.GetString("db.db"),
		MaxIdle: viper.GetInt("db.max_idle"),
		MaxOpen: viper.GetInt("db.max_open"),
		MaxLife: viper.GetDuration("db.max_lite"),
		Level:   viper.GetInt("db.level"),
	}

	instance, err := db.NewConnection(dbOptions)
	if err != nil {
		return err
	}

	_ = store.NewStore(instance)
	return nil
}
