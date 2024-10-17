package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"gitee.com/openeuler/PilotGo-plugin-elk/server/elasticClient"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/pluginclient"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/service/model"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func QueryHandler(ctx *gin.Context) {
	if elasticClient.Global_elastic.Client == nil {
		err := errors.New("global_elastic is null **warn**0") // err top
		errormanager.ErrorTransmit(pluginclient.Global_Context, err, false)
		response.Fail(ctx, nil, err.Error())
		return
	}
	// 处理请求
	defer ctx.Request.Body.Close()
	req_body_bytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		newError(ctx, err)
		return
	}

	req_body := struct {
		Params model.SearchParams `json:"params"`
	}{}

	err = json.Unmarshal(req_body_bytes, &req_body)
	if err != nil {
		newError(ctx, err)
		return
	}

	// 构建查询体
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"range": map[string]interface{}{
							"@timestamp": map[string]string{
								"gte":    req_body.Params.StartTime,
								"lte":    req_body.Params.EndTime,
								"format": "strict_date_optional_time",
							},
						},
					},
				},
			},
		},
	}

	// 将查询体编码为JSON
	queryBytes, err := json.Marshal(query)
	if err != nil {
		newError(ctx, err)
		return
	}

	queryStr := string(queryBytes)

	// 执行搜索请求
	res, err := elasticClient.Global_elastic.Client.Search(
		elasticClient.Global_elastic.Client.Search.WithContext(context.Background()),
		elasticClient.Global_elastic.Client.Search.WithIndex(req_body.Params.Index),
		elasticClient.Global_elastic.Client.Search.WithBody(strings.NewReader(queryStr)),
		// todo...（其他选项如分页、排序等)
	)
	if err != nil {
		newError(ctx, err)
		return
	}
	bodyBytes, _ := io.ReadAll(res.Body)
	fmt.Printf("Response body: %s\n", bodyBytes)
	response.Success(ctx, bodyBytes, "")
}
