package pluginclient

import (
	"context"
	"errors"
	"fmt"

	"gitee.com/openeuler/PilotGo-plugin-elk/server/conf"
	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
)

var Global_Client *client.Client

var Global_Context context.Context

func InitPluginClient() {
	if conf.Global_Config != nil && conf.Global_Config.Elk.Https_enabled {
		PluginInfo.Url = fmt.Sprintf("https://%s", conf.Global_Config.Elk.Addr_target)
	} else if conf.Global_Config != nil && !conf.Global_Config.Elk.Https_enabled {
		PluginInfo.Url = fmt.Sprintf("http://%s", conf.Global_Config.Elk.Addr_target)
	} else {
		err := errors.New("Global_Config is nil")
		logger.Fatal("%+v", err)
	}

	Global_Client = client.DefaultClient(PluginInfo)

	getExtentions()

	Global_Context = context.Background()
}

// 注册插件扩展点
func getExtentions() {
	var ex []common.Extention
	pe1 := &common.PageExtention{
		Type:       common.ExtentionPage,
		Name:       "elk集群部署",
		URL:        "/deploy",
		Permission: "plugin.elk.page/menu",
	}
	pe2 := &common.PageExtention{
		Type:       common.ExtentionPage,
		Name:       "agent状态监听",
		URL:        "/status",
		Permission: "plugin.elk.page/menu",
	}
	pe3 := &common.PageExtention{
		Type:       common.ExtentionPage,
		Name:       "policy配置",
		URL:        "/policy",
		Permission: "plugin.elk.page/menu",
	}
	ex = append(ex, pe1, pe2, pe3)
	Global_Client.RegisterExtention(ex)
}
