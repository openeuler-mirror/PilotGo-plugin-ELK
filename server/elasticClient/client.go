/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Jun 24 09:21:47 2024 +0800
 */
package elasticClient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"gitee.com/openeuler/PilotGo-plugin-elk/server/conf"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/server/pluginclient"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var Global_elastic *ElasticClient_v7

type ElasticClient_v7 struct {
	Client *elastic.Client
	Ctx    context.Context
}

func InitElasticClient() {
	addresses := []string{}
	if conf.Global_Config.Elasticsearch.Https_enabled {
		addresses = append(addresses, fmt.Sprintf("https://%s", conf.Global_Config.Elasticsearch.Addr))
	} else {
		addresses = append(addresses, fmt.Sprintf("http://%s", conf.Global_Config.Elasticsearch.Addr))
	}
	cfg := elastic.Config{
		Addresses: addresses,
		Username:  conf.Global_Config.Elasticsearch.Username,
		Password:  conf.Global_Config.Elasticsearch.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: 5 * time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		},
	}

	es_client, err := elastic.NewClient(cfg)
	if err != nil {
		err = errors.Errorf("failed to init elastic client: %+v **errstackfatal**0", err.Error()) // err top
		errormanager.ErrorTransmit(pluginclient.Global_Context, err, true)
		return
	}

	Global_elastic = &ElasticClient_v7{
		Client: es_client,
		Ctx:    pluginclient.Global_Context,
	}
}

// 通过dsl查询
func (client *ElasticClient_v7) SearchByDsl(index string, body io.Reader) ([]byte, error) {
	resp, err := client.Client.Search(
		client.Client.Search.WithContext(client.Ctx),
		client.Client.Search.WithIndex(index),
		client.Client.Search.WithBody(body),
		client.Client.Search.WithTrackTotalHits(true),
		client.Client.Search.WithPretty(),
	)
	return client.processApiResult(resp, err)
}

// 通过调用template模板查询
func (client *ElasticClient_v7) SearchByTemplate(index string, querybody map[string]interface{}) ([]byte, error) {
	query_body_bytes, err := json.Marshal(querybody)
	if err != nil {
		err = errors.Errorf("%s **errstack**0", err.Error())
		return nil, err
	}
	query_body_reader := bytes.NewReader(query_body_bytes)
	resp, err := client.Client.SearchTemplate(
		query_body_reader,
		client.Client.SearchTemplate.WithContext(context.Background()),
		client.Client.SearchTemplate.WithIndex(index),
		client.Client.SearchTemplate.WithPretty(),
	)
	return client.processApiResult(resp, err)
}

// 处理elasticsearch client接口的返回值
func (client *ElasticClient_v7) processApiResult(_resp *esapi.Response, _err error) ([]byte, error) {
	if _err != nil {
		_err = errors.Errorf("%+v **errstack**0", _err.Error())
		return nil, _err
	}
	defer _resp.Body.Close()
	if _resp.IsError() {
		_err = errors.Errorf("%+v **errstack**0", _resp.String())
		return nil, _err
	} else {
		resp_body_bytes, err := io.ReadAll(_resp.Body)
		if err != nil {
			err = errors.Errorf("%s **errstack**0", err.Error())
			return nil, err
		}
		return resp_body_bytes, nil
	}
}
