package commands

import (
	"github.com/jiangrx816/wechat/core/db"
	"github.com/jiangrx816/wechat/model"

	"github.com/urfave/cli/v2"
)

// 定义表模型与注释的关联结构
type TableWithComment struct {
	Model   interface{} // 表模型，例如 &model.NcLog{}
	Comment string      // 表注释
}

func Migrate() *cli.Command {
	return &cli.Command{
		Name:  "migrate",
		Usage: "数据库迁移",
		Subcommands: []*cli.Command{
			{
				Name:        "up",
				Usage:       "自动迁移数据库",
				Description: "自动迁移数据库",
				Action: func(ctx *cli.Context) error {
					tx := db.MustGet("wechat").Debug()

					// 定义需要迁移的表及其注释
					tables := []TableWithComment{
						{
							Model:   &model.SBookName{},
							Comment: "COMMENT='绘本名称表'", // SBookName 表注释
						},
						{
							Model:   &model.SChinesePicture{},
							Comment: "COMMENT='中文绘本表'", // SChinesePicture 表注释
						},
						{
							Model:   &model.SChinesePictureInfo{},
							Comment: "COMMENT='中文绘本详情表'", // SChinesePictureInfo 表注释
						},
						{
							Model:   &model.SEnglishPicture{},
							Comment: "COMMENT='英文绘本表'", // SEnglishPicture 表注释
						},
						{
							Model:   &model.SEnglishPictureInfo{},
							Comment: "COMMENT='英文绘本详情表'", // SEnglishPictureInfo 表注释
						},
						// 添加其他表及注释...
					}

					// 遍历并逐个迁移表
					for _, table := range tables {
						// 设置表选项（包含注释和其他配置）
						err := tx.Set(
							"gorm:table_options",
							"ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 "+table.Comment, // 合并引擎、字符集和注释
						).AutoMigrate(table.Model)

						if err != nil {
							return err
						}
					}

					return nil
				},
			},
		},
	}
}
