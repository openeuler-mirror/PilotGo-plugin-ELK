/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Jun 24 09:21:47 2024 +0800
 */
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"

	"gitee.com/openeuler/PilotGo-plugin-elk/server/conf"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/pluginclient"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func InitWebServer() {
	if pluginclient.Global_Client == nil {
		err := errors.New("Global_Client is nil **errstackfatal**2") // err top
		errormanager.ErrorTransmit(pluginclient.Global_Context, err, true)
		return
	}

	go func() {
		engine := gin.Default()
		gin.SetMode(gin.ReleaseMode)
		pluginclient.Global_Client.RegisterHandlers(engine)
		InitRouter(engine)
		StaticRouter(engine)

		if conf.Global_Config.Elk.Https_enabled {
			err := engine.RunTLS(conf.Global_Config.Elk.Addr, conf.Global_Config.Elk.Public_certificate, conf.Global_Config.Elk.Private_key)
			if err != nil {
				err = errors.Errorf("%s: %s, %s **errstackfatal**2", err.Error(), conf.Global_Config.Elk.Public_certificate, conf.Global_Config.Elk.Private_key)
				errormanager.ErrorTransmit(pluginclient.Global_Context, err, true)
			}
		} else {
			err := engine.Run(conf.Global_Config.Elk.Addr)
			if err != nil {
				err = errors.Errorf("%s **errstackfatal**2", err.Error())
				errormanager.ErrorTransmit(pluginclient.Global_Context, err, true)
			}
		}
	}()
}

func InitRouter(router *gin.Engine) {
	api := router.Group("/plugin/elk/api")
	{
		api.POST("/create_policy", CreatePolicyHandle)

		api.POST("/log_clusterhost_timeaxis_data", SearchByTemplateHandle)
		api.POST("/log_hostprocess_timeaxis_data", SearchByTemplateHandle)
		api.POST("/log_stream_data", SearchByTemplateHandle)
		api.POST("/log_search", QueryHandler)
		api.POST("/log_advance_search", QueryHandler)
	}

	timeoutapi := router.Group("/plugin/elk/api")
	timeoutapi.Use(TimeoutMiddleware2(15 * time.Second))
	{

	}
}

func TimeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(12*time.Second),
		timeout.WithHandler(func(ctx *gin.Context) {
			ctx.Next()
		}),
		timeout.WithResponse(func(ctx *gin.Context) {
			ctx.JSON(http.StatusGatewayTimeout, gin.H{
				"code": http.StatusGatewayTimeout,
				"msg":  "server response timeout",
				"data": nil,
			})
		}),
	)
}

// 服务器响应超时中间件
func TimeoutMiddleware2(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer func() {
			if !c.GetBool("write") && ctx.Err() == context.DeadlineExceeded {
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}

			cancel()
		}()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
