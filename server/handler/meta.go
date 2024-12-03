/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Tue Jun 25 11:11:50 2024 +0800
 */
package handler

import (
	"gitee.com/openeuler/PilotGo-plugin-elk/server/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/pluginclient"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func newError(ctx *gin.Context, err error) {
	err = errors.Errorf("fail to search data: %s **errstack**0", err.Error()) // err top
	errormanager.ErrorTransmit(pluginclient.Global_Context, err, false)
	response.Fail(ctx, nil, err.Error())
}

func wrapError(ctx *gin.Context, err error) {
	err = errors.Wrap(err, " **0") // err top
	errormanager.ErrorTransmit(pluginclient.Global_Context, err, false)
	response.Fail(ctx, nil, err.Error())
}
