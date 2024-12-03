/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Jun 24 09:21:47 2024 +0800
 */
package errormanager

import (
	"context"
	"os"

	"gitee.com/openeuler/PilotGo/sdk/logger"
)

type ElkPluginError struct {
	Err    error
	Cancel context.CancelFunc
}

/*
@ctx:	插件服务端初始上下文（默认为pluginclient.Global_Context）

@err:	最终生成的error

@exit_after_print: 打印完错误链信息后是否结束主程序
*/
func ErrorTransmit(ctx context.Context, err error, exit_after_print bool) {
	if Global_ErrorManager == nil {
		logger.Error("globalerrormanager is nil")
		os.Exit(1)
	}

	if exit_after_print {
		cctx, cancelF := context.WithCancel(ctx)
		Global_ErrorManager.ErrCh <- &ElkPluginError{
			Err:    err,
			Cancel: cancelF,
		}
		<-cctx.Done()
		close(Global_ErrorManager.ErrCh)
		os.Exit(1)
	}

	Global_ErrorManager.ErrCh <- &ElkPluginError{
		Err: err,
	}
}
