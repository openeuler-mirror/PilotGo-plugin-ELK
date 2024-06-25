package handler

import (
	"bytes"
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
		Index string      `json:"index"`
		DSL   interface{} `json:"dsl"`
	}{}
	err = json.Unmarshal(req_body_bytes, &req_body)
	if err != nil {
		newError(ctx, err)
		return
	}
	query_body_bytes, err := json.Marshal(req_body.DSL)
	if err != nil {
		newError(ctx, err)
		return
	}

	resp_body_bytes, err := elasticClient.Global_elastic.Search(req_body.Index, bytes.NewReader(query_body_bytes))
	if err != nil {
		wrapError(ctx, err)
		return
	}
	data, err := cluster.ProcessLogTimeAixsData(resp_body_bytes)
	if err != nil {
		wrapError(ctx, err)
		return
	}
	response.Success(ctx, data, "")
}
