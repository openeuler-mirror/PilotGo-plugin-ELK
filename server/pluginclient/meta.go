package pluginclient

import "gitee.com/openeuler/PilotGo/sdk/plugin/client"

const Version = "1.0.1"

var PluginInfo = &client.PluginInfo{
	Name:        "elk",
	Version:     Version,
	Description: "connect PilotGo and elk",
	Author:      "wangjunqi",
	Email:       "wangjunqi@kylinos.cn",
	Url:         "", // 客户端访问插件服务端的地址
	PluginType:  "micro-app",
}
