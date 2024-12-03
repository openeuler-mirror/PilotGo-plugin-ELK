/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Jun 24 09:21:47 2024 +0800
 */
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
