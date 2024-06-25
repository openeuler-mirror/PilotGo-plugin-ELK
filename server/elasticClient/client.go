package elasticClient

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"gitee.com/openeuler/PilotGo-plugin-elk/conf"
	"gitee.com/openeuler/PilotGo-plugin-elk/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/pluginclient"
	elastic "github.com/elastic/go-elasticsearch/v7"
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
}

func (client *ElasticClient_v7) Search(index string, body io.Reader) ([]byte, error) {
	resp, err := client.Client.Search(
		client.Client.Search.WithContext(context.Background()),
		client.Client.Search.WithIndex(index),
		client.Client.Search.WithBody(body),
		client.Client.Search.WithTrackTotalHits(true),
		client.Client.Search.WithPretty(),
	)
	if err != nil {
		err = errors.Errorf("%+v **errstack**0", err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		err = errors.Errorf("%+v **errstack**0", resp.String())
		return nil, err
	} else {
		resp_body_bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			err = errors.Errorf("%+v **warn**0", err.Error())
			return nil, err
		}
		return resp_body_bytes, nil
	}
}
