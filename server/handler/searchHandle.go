package handler

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"

	"gitee.com/openeuler/PilotGo-plugin-elk/elasticClient"
	"gitee.com/openeuler/PilotGo-plugin-elk/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/pluginclient"
	"gitee.com/openeuler/PilotGo-plugin-elk/service/cluster"
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
		Index  string                 `json:"index"`
		Id     string                 `json:"id"`
		Params map[string]interface{} `json:"params"`
	}{}
	err = json.Unmarshal(req_body_bytes, &req_body)
	if err != nil {
		newError(ctx, err)
		return
	}
	query_body := map[string]interface{}{
		"id":     req_body.Id,
		"params": req_body.Params,
	}

	data, err := cluster.ProcessLogTimeAixsData(req_body.Index, query_body)
	if err != nil {
		wrapError(ctx, err)
		return
	}
	response.Success(ctx, data, "")
}
