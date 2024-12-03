/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: yajun <yajun@kylinos.cn>
 * Date: Thu Sep 26 13:48:18 2024 +0800
 */
package handler

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_buildESQuery(t *testing.T) {
	type args struct {
		conditions string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test-1",
			args: args{
				conditions: "name:John",
			},
			want: `{"query":{"bool":{"must":[{"term":{"name":"John"}}]}}}`,
		},
		{
			name: "test-2",
			args: args{
				conditions: "name:John,age:25",
			},
			want: `{"query":{"bool":{"must":[{"term":{"name":"John"}},{"term":{"age":25}}]}}}`,
		},
		{
			name: "test-3",
			args: args{
				conditions: "name:John,age:",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildESQuery(tt.args.conditions)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildESQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.JSONEq(t, tt.want, got)
		})
	}
}

func Test_convertToESTermsQuery(t *testing.T) {
	queryStr := "A:1|2|3"
	esQuery, err := convertToESTermsQuery(queryStr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("ES Query:", esQuery)

	// 另一个示例
	queryStr2 := "B:a|b|c"
	esQuery2, err := convertToESTermsQuery(queryStr2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("ES Query 2:", esQuery2)
}

func Test_parseExpression1(t *testing.T) {
	// expr := "A=a,B=(C=c|d)|e"  
	// cond, err := parseExpression(expr)  
	// if err != nil {  
	// 	fmt.Println("Error:", err)  
	// 	return  
	// }  
	// fmt.Println(cond)  
	// fmt.Println(cond.Evaluate(map[string]string{"A": "a", "B": "e", "C": "d"})) // 应该返回 true  
	expr := "(A=a,B=b)|(C=c|D=d)"  
    cond, err := parseExpression1(expr)  
    if err != nil {  
        fmt.Println("Error:", err)  
        return  
    }  
    fmt.Println(cond) // 这将打印出表达式树的字符串表示，但可能不是您期望的格式 
}

/**
* 字符串分割查询
* 格式：key1:value1,key2:value2...
 */
 func buildESQuery(conditions string) (string, error) {
	// 使用strings.Split分割conditions字符串
	pairs := strings.Split(conditions, ",")
	// 初始化一个map来存储键值对
	queryMap := make(map[string]interface{})
	mustClause := make([]map[string]interface{}, 0)
	for _, pair := range pairs {
		// 进一步分割键值对
		keyVal := strings.Split(pair, ":")
		if len(keyVal) != 2 {
			return "", fmt.Errorf("invalid condition format: %s", pair)
		}
		// 构建一个匹配条件的map
		condition := make(map[string]interface{})
		condition["term"] = map[string]interface{}{keyVal[0]: keyVal[1]}
		// 将这个条件添加到mustClause中
		mustClause = append(mustClause, condition)
	}
	// 构建最终的查询
	queryMap["query"] = map[string]interface{}{
		"bool": map[string][]map[string]interface{}{
			"must": mustClause,
		},
	}
	// 将map转换为JSON字符串
	queryJSON, err := json.Marshal(queryMap)
	if err != nil {
		return "", err
	}
	return string(queryJSON), nil
}

// QueryBuilder 结构体用于构建Elasticsearch查询
type QueryBuilder struct {
	Bool struct {
		Must []map[string]interface{} `json:"must"`
	} `json:"bool"`
}

func parseAndBuildQuery(input string) ([]byte, error) {
	// 分离键值对
	pairs := strings.Split(input, "，")
	var mustConditions []map[string]interface{}

	for _, pair := range pairs {
		keyVal := strings.SplitN(pair, "：", 2)
		if len(keyVal) != 2 {
			return nil, fmt.Errorf("invalid key-value pair: %s", pair)
		}

		key := keyVal[0]
		values := strings.Split(keyVal[1], "|")

		// 为每个值构建一个"or"查询
		var orConditions []map[string]interface{}
		for _, value := range values {
			orConditions = append(orConditions, map[string]interface{}{
				"term": map[string]interface{}{
					key: value,
				},
			})
		}

		// 将"or"条件包装为"must"
		mustConditions = append(mustConditions, map[string]interface{}{
			"bool": map[string]interface{}{
				"should":               orConditions,
				"minimum_should_match": 1, // 至少匹配一个
			},
		})
	}

	// 构建最终查询
	query := QueryBuilder{
		Bool: struct {
			Must []map[string]interface{} `json:"must"`
		}{
			Must: mustConditions,
		},
	}

	// 转换为JSON
	return json.MarshalIndent(query, "", "  ")
}

// 定义一个结构体，表示ES的terms查询
type TermsQuery struct {
	Terms struct {
		Field  string   `json:"field"`
		Values []string `json:"values"`
	} `json:"terms"`
}

// 将形如"A:1|2|3"或"A:a|b|c"的字符串转换为ES的terms查询
func convertToESTermsQuery(queryStr string) (string, error) {
	// 分割字符串以获取字段名和值
	parts := strings.SplitN(queryStr, ":", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid query format")
	}
	fieldName := parts[0]
	valuesStr := parts[1]

	// 分割值字符串以获取值的数组
	values := strings.Split(valuesStr, "|")

	// 构建ES查询
	esQuery := TermsQuery{
		Terms: struct {
			Field  string   `json:"field"`
			Values []string `json:"values"`
		}{
			Field:  fieldName,
			Values: values,
		},
	}

	// 将查询转换为JSON字符串
	queryJSON, err := json.Marshal(esQuery)
	if err != nil {
		return "", err
	}

	return string(queryJSON), nil
}

// 假设这是一个递归地构建的ES查询的一部分
type ESQueryPart interface{}

// 一个简单的ES查询条件
type TermQuery struct {
	Term struct {
		Field string `json:"field"`
		Value string `json:"value"`
	} `json:"term"`
}

// 一个嵌套的ES查询（假设它已经被解析为某种结构）
// 在实际应用中，这个可能是一个更复杂的结构体
type NestedQuery struct {
	// 这里只是简单示例，实际中可能包含多个字段和查询条件
	Bool struct {
		Must []ESQueryPart `json:"must,omitempty"`
	} `json:"bool,omitempty"`
}

// 递归函数，解析查询字符串并构建ES查询
func parseQueryString(queryStr string) ([]ESQueryPart, error) {
	// 假设查询字符串是由逗号分隔的字段=值对，可能包含嵌套查询
	parts := strings.Split(queryStr, ",")
	var queries []ESQueryPart

	for _, partStr := range parts {
		// 检查是否为嵌套查询的格式（这里只是简单模拟）
		if strings.HasPrefix(partStr, "=") && strings.HasSuffix(partStr, ")") {
			// 假设我们有一个函数可以从字符串中提取嵌套查询（这里只是添加一个占位符）
			// ...
			// 由于我们不能直接从字符串中解析，这里只是添加一个NestedQuery的实例
			queries = append(queries, NestedQuery{})
		} else {
			// 处理简单的字段=值对
			parts := strings.SplitN(partStr, "=", 2)
			if len(parts) != 2 {
				return nil, errors.New("invalid query format")
			}
			field, value := parts[0], parts[1]
			queries = append(queries, TermQuery{
				Term: struct {
					Field string `json:"field"`
					Value string `json:"value"`
				}{
					Field: field,
					Value: value,
				},
			})
		}
	}

	return queries, nil
}

// parseExpression 解析整个表达式
// func parseExpression(expr string) (Condition, error) {
// 	// 使用正则表达式来简化处理，这里假设表达式格式正确
// 	// 更好的实现可能需要一个完整的解析器
// 	andRE := regexp.MustCompile(`\([^()\|,]+,([^()\|,]+)\)`)
// 	orRE := regexp.MustCompile(`\([^()\|,]+\|([^()\|,]+)\)`)

// 	// 处理或逻辑
// 	for orRE.MatchString(expr) {
// 		matches := orRE.FindStringSubmatch(expr)
// 		if len(matches) > 1 {
// 			left, err := parseExpression(expr[:orRE.FindStringIndex(expr)[0]])
// 			if err != nil {
// 				return nil, err
// 			}
// 			right, err := parseExpression(matches[1])
// 			if err != nil {
// 				return nil, err
// 			}
// 			expr = expr[orRE.FindStringIndex(expr)[1]:]
// 			expr = fmt.Sprintf("%s|%s%s", left.String(), right.String(), expr)
// 			if strings.HasPrefix(expr, "|") {
// 				expr = expr[1:]
// 			}
// 			orNode := OrNode{
// 				Children: []Condition{left, right},
// 			}
// 			return parseExpression(fmt.Sprintf("(%s)", expr))
// 		}
// 	}

// 	// 处理与逻辑
// 	for andRE.MatchString(expr) {
// 		matches := andRE.FindStringSubmatch(expr)
// 		if len(matches) > 1 {
// 			left, err := parseExpression(expr[:andRE.FindStringIndex(expr)[0]])
// 			if err != nil {
// 				return nil, err
// 			}
// 			right, err := parseExpression(matches[1])
// 			if err != nil {
// 				return nil, err
// 			}
// 			expr = expr[andRE.FindStringIndex(expr)[1]:]
// 			expr = fmt.Sprintf("%s,%s%s", left.String(), right.String(), expr)
// 			if strings.HasPrefix(expr, ",") {
// 				expr = expr[1:]
// 			}
// 			andNode := AndNode{
// 				Children: []Condition{left, right},
// 			}
// 			return parseExpression(fmt.Sprintf("(%s)", expr))
// 		}
// 	}

// 	// 处理单个条件
// 	cond, err := parseCondition(expr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return cond, nil
// }