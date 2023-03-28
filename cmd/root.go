package cmd

import (
	"github.com/9d4/netpilot/database"
	"github.com/9d4/netpilot/server"
	"github.com/9d4/netpilot/worker"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "netpilot",
	Short: "Start netpilot",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := database.ConnectDB(&database.Config{
			Host:     v.GetString("DB_HOST"),
			Port:     v.GetInt("DB_PORT"),
			Database: v.GetString("DB_NAME"),
			Username: v.GetString("DB_USERNAME"),
			Password: v.GetString("DB_PASSWORD"),
		})
		if err != nil {
			jww.FATAL.Fatal(err)
		}

		_, err = database.ConnectRedis(&database.RedisConfig{
			Address:  v.GetString("REDIS_ADDRESS"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		})
		if err != nil {
			jww.FATAL.Fatal(err)
		}

		worker.RunBoardWorker()
		server.Start(server.NewConfig(v))
	},
}

var (
	v     = viper.NewWithOptions(viper.EnvKeyReplacer(strings.NewReplacer("-", "_")))
	flags = flag.NewFlagSet(rootCmd.Name(), flag.ContinueOnError)
)

func init() {
	loadLogger()
	loadEnv()
	loadFlags()
	loadConfig()
}

func loadLogger() {
	logWriter, err := os.OpenFile("netpilot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Println("Unable to create log file:", err)
	}

	if err == nil {
		jww.SetLogOutput(logWriter)
	}

	jww.SetFlags(log.Lshortfile)
	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelInfo)
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		jww.FATAL.Fatal("Error reading .env file")
	}
}

func loadFlags() {
	// flags installation go here
	rootCmd.PersistentFlags().AddFlagSet(flags)
	v.BindPFlags(rootCmd.PersistentFlags())
}

func loadConfig() {
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		return
	}
}
