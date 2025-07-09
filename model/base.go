package model

import (
	"github.com/jiangrx816/wechat/core/db"

	"gorm.io/gorm"
)

func Default() *gorm.DB { return db.MustGet("wechat") }
