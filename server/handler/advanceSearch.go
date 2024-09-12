package handler

import (
	"encoding/json"
	"fmt"

	"gitee.com/openeuler/PilotGo-plugin-elk/server/elasticClient"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/pluginclient"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/service/model"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// ES查询条件转换器："与"查询
func buildElasticQuery(conditions map[string]interface{}) (string, error) {
	var query map[string]interface{}
	// 初始化bool查询
	boolQuery := map[string]interface{}{
		"bool": map[string][]map[string]interface{}{
			"must": make([]map[string]interface{}, 0),
		},
	}

	// 遍历条件，将每个条件添加到must数组中
	for key, value := range conditions {
		matchQuery := map[string]interface{}{
			"match": map[string]interface{}{
				key: value,
			},
		}
		mustSlice := boolQuery["bool"].(map[string][]map[string]interface{})["must"]
		mustSlice = append(mustSlice, matchQuery)
		boolQuery["bool"].(map[string][]map[string]interface{})["must"] = mustSlice
	}

	// 如果没有条件，则使用match_all查询
	if len(boolQuery["bool"].(map[string][]map[string]interface{})["must"]) == 0 {
		boolQuery["bool"] = map[string]interface{}{
			"must": []map[string]interface{}{
				{"match_all": map[string]interface{}{}},
			},
		}
	}

	query = map[string]interface{}{
		"query": boolQuery,
	}

	// 将查询转换为JSON字符串
	jsonQuery, err := json.MarshalIndent(query, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonQuery), nil
}

func LogAdvanceSearchHandel(ctx *gin.Context) {
	if elasticClient.Global_elastic.Client == nil {
		err := errors.New("global_elastic is null **warn**0")
		errormanager.ErrorTransmit(pluginclient.Global_Context, err, false)
		response.Fail(ctx, nil, err.Error())
		return
	}
	// 处理请求
	defer ctx.Request.Body.Close()
	searchParams := struct {
		Params model.SearchParams `json:"params"`
	}{}
	// 尝试将请求体中的JSON数据绑定到searchParams
	if err := ctx.ShouldBindJSON(&searchParams); err != nil {
		ctx.JSON(400, gin.H{"error": "解析JSON数据失败", "detail": err.Error()})
		return
	}
}

/**
	* 时间范围查询
**/
// 假设 conditions 中的值可能是：
// - 字符串、数字等，用于 match 查询
// - []interface{}，表示时间范围（如 [startTime, endTime]），用于 range 查询
func buildElasticQuery1(conditions map[string]interface{}) (string, error) {
	var query map[string]interface{}
	boolQuery := map[string]interface{}{
		"bool": map[string][]map[string]interface{}{
			"must": make([]map[string]interface{}, 0),
		},
	}

	for key, value := range conditions {
		switch v := value.(type) {
		case []interface{}:
			// 检查是否为时间范围（这里假设长度为2的切片是时间范围）
			if len(v) == 2 {
				// 时间处理
				startTime, endTime := v[0], v[1]
				startTimeStr, endTimeStr := fmt.Sprintf("%v", startTime), fmt.Sprintf("%v", endTime)

				rangeQuery := map[string]interface{}{
					"range": map[string]map[string]interface{}{
						key: {
							"gte": startTimeStr,
							"lte": endTimeStr,
						},
					},
				}
				boolQuery["bool"].(map[string][]map[string]interface{})["must"] = append(boolQuery["bool"].(map[string][]map[string]interface{})["must"], rangeQuery)
			} else {
				// 如果不是长度为2的切片，则可能是一个错误或不支持的类型，可以处理为错误或忽略，后续处理
				continue
			}
		default:
			// 对于非切片类型，使用 match 查询
			matchQuery := map[string]interface{}{
				"match": map[string]interface{}{
					key: value,
				},
			}
			boolQuery["bool"].(map[string][]map[string]interface{})["must"] = append(boolQuery["bool"].(map[string][]map[string]interface{})["must"], matchQuery)
		}
	}

	// 如果没有条件，则使用 match_all 查询
	if len(boolQuery["bool"].(map[string][]map[string]interface{})["must"]) == 0 {
		boolQuery["bool"] = map[string]interface{}{
			"must": []map[string]interface{}{
				{"match_all": map[string]interface{}{}},
			},
		}
	}

	query = map[string]interface{}{
		"query": boolQuery,
	}

	jsonQuery, err := json.MarshalIndent(query, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonQuery), nil
}
