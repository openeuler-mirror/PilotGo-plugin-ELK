package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"gitee.com/openeuler/PilotGo-plugin-elk/server/elasticClient"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/pluginclient"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	mapping = `
	{
		 "mappings": {
            "properties": {
                "v": {
                    "type": "text",
                    "fields": {
                        "keyword": {
                            "type": "keyword",
							"normalizer": "lowercase"
                        }
                    }
                }
            }
        }
	}`
)

// 创建索引请求
func CreateeIndex(ctx *gin.Context, indexName string) {
	defer ctx.Request.Body.Close()
	req := esapi.IndicesCreateRequest{
		Index: indexName,
		// 其它参数
	}
	res, err := req.Do(context.Background(), elasticClient.Global_elastic.Client)

	// 关闭响应体
	defer res.Body.Close()
	// 检查响应状态
	if err != nil {
		err = errors.Errorf("%+v **warn**0", err.Error())
		errormanager.ErrorTransmit(pluginclient.Global_Context, err, false)
		response.Fail(ctx, nil, errors.Cause(err).Error())
		return
	}
	if res.IsError() {
		err = errors.Errorf("%+v **warn**0", res.String())
		errormanager.ErrorTransmit(pluginclient.Global_Context, err, false)
		response.Fail(ctx, nil, res.String())
		return
	} else {
		resp_body_data := map[string]interface{}{}
		resp_body_bytes, err := io.ReadAll(res.Body)
		if err != nil {
			err = errors.Errorf("%+v **warn**0", err.Error())
			errormanager.ErrorTransmit(pluginclient.Global_Context, err, false)
			response.Fail(ctx, nil, res.String())
			return
		}
		json.Unmarshal(resp_body_bytes, &resp_body_data)
		response.Success(ctx, resp_body_data, "")
	}
}

/*
*
创建索引
*/
func CreateNewIndex(indexName string) bool {
	exists, err := elasticClient.Global_elastic.Client.Indices.Exists([]string{indexName})
	if err != nil {
		fmt.Printf("indexName Existed ! err is %s\n", err)
		return false
	}
	/**
	不存在就创建
	*/
	if exists.StatusCode != 200 {
		res, err := elasticClient.Global_elastic.Client.Indices.Create(indexName, elasticClient.Global_elastic.Client.Indices.Create.WithBody(strings.NewReader(mapping)))
		if err != nil {
			log.Fatal(err)
			return false
		}
		defer res.Body.Close()

		var createRes map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&createRes); err != nil {
			log.Fatal(err)
			return false
		}

		if acknowledged, ok := createRes["acknowledged"].(bool); ok && acknowledged {
			fmt.Printf("indexName create success: %s\n", res.String())
		} else {
			fmt.Printf("indexName create fail: %s\n", res.String())
			return false
		}
	}
	return true
}

// DeleteIndex 删除索引
func DeleteIndex(ctx *gin.Context, indexName []string) {
	// _, err := EsClient.Indices.Delete(indexName).Do(context.Background())   v8方式不适用
	req := esapi.IndicesDeleteRequest{
		Index: indexName,
	}
	_, err := req.Do(context.Background(), elasticClient.Global_elastic.Client)
	if err != nil {
		err = errors.Errorf("%+v **warn**0", err.Error())
		errormanager.ErrorTransmit(pluginclient.Global_Context, err, false)
		response.Fail(ctx, nil, errors.Cause(err).Error())
		return
	}
	// fmt.Printf("delete index successed,indexName:%s", indexName)
	return
}
