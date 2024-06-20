package handler

import (
	"context"
	"encoding/json"
	"io"

	"github.com/pkg/errors"

	"gitee.com/openeuler/PilotGo-plugin-elk/elasticClient"
	"gitee.com/openeuler/PilotGo-plugin-elk/errormanager"
	"gitee.com/openeuler/PilotGo-plugin-elk/pluginclient"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func SearchHandle(ctx *gin.Context) {
	index := ctx.Query("index")
	if index == "" {
		err := errors.Errorf("%+v **warn**0", "index is null")
		errormanager.ErrorTransmit(pluginclient.Global_Context, err, false)
		response.Fail(ctx, nil, err.Error())
		return
	}

	defer ctx.Request.Body.Close()
	if elasticClient.Global_elastic.Client == nil {
		err := errors.Errorf("%+v **warn**0", "global_elastic is null")
		errormanager.ErrorTransmit(pluginclient.Global_Context, err, false)
		response.Fail(ctx, nil, err.Error())
		return
	}
	resp, err := elasticClient.Global_elastic.Client.Search(
		elasticClient.Global_elastic.Client.Search.WithContext(context.Background()),
		elasticClient.Global_elastic.Client.Search.WithIndex(index),
		elasticClient.Global_elastic.Client.Search.WithBody(ctx.Request.Body),
		elasticClient.Global_elastic.Client.Search.WithTrackTotalHits(true),
		elasticClient.Global_elastic.Client.Search.WithPretty(),
	)
	defer resp.Body.Close()
	if err != nil {
		err = errors.Errorf("%+v **warn**0", err.Error())
		errormanager.ErrorTransmit(pluginclient.Global_Context, err, false)
		response.Fail(ctx, nil, errors.Cause(err).Error())
		return
	}
	if resp.IsError() {
		err = errors.Errorf("%+v **warn**0", resp.String())
		errormanager.ErrorTransmit(pluginclient.Global_Context, err, false)
		response.Fail(ctx, nil, resp.String())
		return
	} else {
		resp_body_data := map[string]interface{}{}
		resp_body_bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			err = errors.Errorf("%+v **warn**0", err.Error())
			errormanager.ErrorTransmit(pluginclient.Global_Context, err, false)
			response.Fail(ctx, nil, resp.String())
			return
		}
		json.Unmarshal(resp_body_bytes, &resp_body_data)
		response.Success(ctx, resp_body_data, "")
	}
}
