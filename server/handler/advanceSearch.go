package handler

import (
	"encoding/json"
	"fmt"
	"strings"

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

// Condition:条件树的接口
type Condition interface {
	Evaluate(context map[string]string) bool
	String() string
}

// ConditionNode:条件节点
type ConditionNode struct {
	Key   string
	Value string
}

func (c ConditionNode) Evaluate(context map[string]string) bool {
	return context[c.Key] == c.Value
}

func (c ConditionNode) String() string {
	return fmt.Sprintf("%s=%s", c.Key, c.Value)
}

// AndNode:与逻辑
type AndNode struct {
	Children []Condition
}

func (a AndNode) Evaluate(context map[string]string) bool {
	for _, child := range a.Children {
		if !child.Evaluate(context) {
			return false
		}
	}
	return true
}

func (a AndNode) String() string {
	parts := make([]string, len(a.Children))
	for i, child := range a.Children {
		parts[i] = child.String()
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, ","))
}

// OrNode:或逻辑
type OrNode struct {
	Children []Condition
}

func (o OrNode) Evaluate(context map[string]string) bool {
	for _, child := range o.Children {
		if child.Evaluate(context) {
			return true
		}
	}
	return false
}

func (o OrNode) String() string {
	parts := make([]string, len(o.Children))
	for i, child := range o.Children {
		parts[i] = child.String()
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, "|"))
}

// 解析单个条件
func parseCondition(expr string) (Condition, error) {
	parts := strings.SplitN(expr, "=", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid condition format")
	}
	return ConditionNode{Key: parts[0], Value: parts[1]}, nil
}

// 解析整个表达式并返回构建的条件树
func parseExpression1(expr string) (Condition, error) {

	// 解析为单个条件
	cond, err := parseCondition(expr)
	if err == nil {
		return cond, nil
	}

	// 解析为 And 表达式
	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") && strings.Contains(expr, ",") {
		innerExpr := expr[1 : len(expr)-1] // 去掉外层的括号
		parts := strings.Split(innerExpr, ",")
		children := make([]Condition, len(parts))
		for i, part := range parts {
			child, err := parseExpression1(strings.TrimSpace(part)) // 递归解析每个部分
			if err != nil {
				return nil, err
			}
			children[i] = child
		}
		return AndNode{Children: children}, nil
	}

	// 解析为 Or 表达式
	if strings.Contains(expr, "|") {
		parts := strings.Split(expr, "|")
		children := make([]Condition, len(parts))
		for i, part := range parts {
			child, err := parseExpression1(strings.TrimSpace(part)) // 递归解析每个部分
			if err != nil {
				return nil, err
			}
			children[i] = child
		}
		return OrNode{Children: children}, nil
	}

	// 无匹配项，返回错误
	return nil, errors.New("invalid expression format")
}
