package commands

import (
	"fmt"

	"github.com/jiangrx816/wechat/http/router"
	"github.com/spf13/viper"

	"github.com/jiangrx816/wechat/core/graceful"

	rxLog "github.com/jiangrx816/wechat/core/log"
	"github.com/jiangrx816/wechat/core/server"
	"github.com/urfave/cli/v2"
)

func Serve(c *cli.Context) error {
	port := viper.GetViper().GetString("port")
	if port == "" {
		rxLog.Sugar().Fatalf("yml port is empty %#v", port)
	}
	httpPort := fmt.Sprintf(":%s", port)
	// 运行HTTP服务
	graceful.Start(server.NewHttp(server.Addr(httpPort), server.Router(router.All())))

	graceful.Wait()
	return nil
}
