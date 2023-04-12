package yarxgo

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/seaung/yarx-go/internal/pkg/log"
	"github.com/seaung/yarx-go/internal/yarx-go/store"
	"github.com/seaung/yarx-go/pkg/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	recommendedHomeDir = ".yarxgo"

	defaultConfigName = "yarxgo.yml"
)

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()

		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(home, recommendedHomeDir))

		viper.AddConfigPath(".")

		viper.SetConfigType("yaml")

		viper.SetConfigName(defaultConfigName)
	}

	viper.AutomaticEnv()

	viper.SetEnvPrefix("YARX_GO")

	replacer := strings.NewReplacer(".", "_")

	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		log.Errorw("Failed to read viper configuration file", "err", err)
	}

	log.Debug("Using config file", "file", viper.ConfigFileUsed())
}

func initStore() error {
	dbOptions := &db.MySQLOptions{
		Host:                  viper.GetString("db.host"),
		Username:              viper.GetString("db.username"),
		Password:              viper.GetString("db.password"),
		Database:              viper.GetString("db.database"),
		MaxIdleConnections:    viper.GetInt("db.max-idle-connections"),
		MaxOpenConnections:    viper.GetInt("db.max-open-connections"),
		MaxConnectionLifeTime: viper.GetDuration("db.max-connections-life-time"),
		LogLevel:              viper.GetInt("db.log-level"),
	}

	ins, err := db.NewMySQL(dbOptions)
	if err != nil {
		return err
	}

	_ = store.NewStore(ins)

	return nil
}

func logOptions() *log.Options {
	return &log.Options{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
	}
}
