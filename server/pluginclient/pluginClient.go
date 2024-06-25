package pluginclient

import (
	"context"

	"gitee.com/openeuler/PilotGo/sdk/common"
	"gitee.com/openeuler/PilotGo/sdk/plugin/client"
)

var Global_Client *client.Client

var Global_Context context.Context

func InitPluginClient() {
	Global_Client = client.DefaultClient(PluginInfo)

	// 注册插件扩展点
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

	Global_Context = context.Background()
}
