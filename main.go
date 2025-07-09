package main

import (
	"os"

	"github.com/jiangrx816/wechat/core/commands"
	rxMysql "github.com/jiangrx816/wechat/core/db"
	rxLog "github.com/jiangrx816/wechat/core/log"
	"github.com/jiangrx816/wechat/utils"

	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var configFile string

func main() {
	app := cli.NewApp()
	app.Action = commands.Serve
	app.Before = initConfig
	app.Commands = commands.Commands
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config",
			Value:       "", // 默认从config目录读取
			Usage:       "specify the location of the configuration file",
			Required:    false,
			Destination: &configFile,
		},
	}
	if err := app.Run(os.Args); err != nil {
		rxLog.Sugar().Fatal(err)
	}
}

func initConfig(*cli.Context) error {
	viper.SetDefault("app", "wechat")
	if err := utils.LoadConfigInFile(configFile); err != nil {
		return err
	}
	if err := rxLog.InitFromViper(); err != nil {
		return err
	}
	if err := rxMysql.InitMysqlDB(); err != nil {
		return err
	}
	return nil
}
