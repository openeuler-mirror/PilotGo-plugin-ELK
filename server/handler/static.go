/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Jun 24 09:21:47 2024 +0800
 */
//go:build !production
// +build !production

package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func StaticRouter(router *gin.Engine) {
	static := router.Group("/plugin/elk")
	{
		static.Static("/assets", "web/dist/static")
		static.StaticFile("/", "web/dist/index.html")

		// 解决页面刷新404的问题
		router.NoRoute(func(c *gin.Context) {
			if !strings.HasPrefix(c.Request.RequestURI, "/plugin/elk/api") {
				c.File("web/dist/index.html")
				return
			}
			c.AbortWithStatus(http.StatusNotFound)
		})
	}
}
