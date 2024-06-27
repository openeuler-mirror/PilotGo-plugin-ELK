package handler

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/pkg/errors"

	"gitee.com/openeuler/PilotGo-plugin-elk/server/elasticClient"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/global"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/pluginclient"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/service/cluster"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

// 查询日志时间轴相关数据
func Search_LogTimeAxisDataHandle(ctx *gin.Context) {
	if elasticClient.Global_elastic.Client == nil {
		err := errors.New("global_elastic is null **warn**0") // err top
		errormanager.ErrorTransmit(pluginclient.Global_Context, err, false)
		response.Fail(ctx, nil, err.Error())
		return
	}
	// process request body
	defer ctx.Request.Body.Close()
	req_body_bytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		newError(ctx, err)
		return
	}
	req_body := struct {
		Id     string                 `json:"id"`
		Params map[string]interface{} `json:"params"`
	}{}
	err = json.Unmarshal(req_body_bytes, &req_body)
	if err != nil {
		newError(ctx, err)
		return
	}
	params := map[string]interface{}{}
	for k, v := range req_body.Params {
		if strings.Contains(k, "query_range") {
			params[k] = global.GetTime_Timestamp2DateTime(int64(v.(float64)))
			continue
		}
		params[k] = v
	}
	
	// TODO: 判断索引
	index := "logs-*"

	query_body := map[string]interface{}{
		"id":     req_body.Id,
		"params": params,
	}
	data, err := cluster.ProcessLogTimeAixsData(index, query_body)
	if err != nil {
		wrapError(ctx, err)
		return
	}
	response.Success(ctx, data, "")
}
