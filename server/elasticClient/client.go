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
	"strings"
	"time"

	"github.com/pkg/errors"

	"gitee.com/openeuler/PilotGo-plugin-elk/conf"
	"gitee.com/openeuler/PilotGo-plugin-elk/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/global/template"
	"gitee.com/openeuler/PilotGo-plugin-elk/pluginclient"
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
	if conf.Global_Config.Elk.Https_enabled {
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
			ResponseHeaderTimeout: time.Second,
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

	Global_elastic.initSearchTemplate()
}

// 在elasticsearch中添加查询模板
func (client *ElasticClient_v7) initSearchTemplate() {
	for key, value := range template.DSL_template_map {
		reqbody := strings.NewReader(value)
		_, err := client.Client.PutScript(
			key,
			reqbody,
			client.Client.PutScript.WithContext(client.Ctx),
			client.Client.PutScript.WithPretty(),
		)
		if err != nil {
			err = errors.Errorf("fail to put script: %s, %s **warn**0", key, err.Error()) // err top
			errormanager.ErrorTransmit(pluginclient.Global_Context, err, false)
		}
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
