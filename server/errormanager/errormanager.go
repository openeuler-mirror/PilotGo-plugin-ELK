package errormanager

import (
	"fmt"
	"io"
	"os"
	"strings"

	"gitee.com/openeuler/PilotGo-plugin-elk/server/conf"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"github.com/pkg/errors"
)

var Global_ErrorManager *ErrorManager

type ErrorManager struct {
	ErrCh chan *ElkPluginError

	Out io.Writer
}

func InitErrorManager() {
	Global_ErrorManager = &ErrorManager{
		ErrCh: make(chan *ElkPluginError, 20),
	}

	switch conf.Global_Config.Logopts.Driver {
	case "stdout":
		Global_ErrorManager.Out = os.Stdout
	case "file":
		logfile, err := os.OpenFile(conf.Global_Config.Logopts.Path, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
		Global_ErrorManager.Out = logfile
	}

	go func(ch <-chan *ElkPluginError) {
		for elkerr := range ch {
			if elkerr.Err != nil {
				errarr := strings.Split(errors.Cause(elkerr.Err).Error(), "**")
				if len(errarr) < 2 {
					logger.Error("topoerror type required in root error (err: %+v)", elkerr.Err)
					os.Exit(1)
				}

				switch errarr[1] {
				case "debug": // 只打印最底层error的message，不展开错误链的调用栈
					logger.Debug("%+v\n", strings.Split(errors.Cause(elkerr.Err).Error(), "**")[0])
				case "warn": // 只打印最底层error的message，不展开错误链的调用栈
					logger.Warn("%+v\n", strings.Split(errors.Cause(elkerr.Err).Error(), "**")[0])
				case "errstack": // 打印错误链的调用栈
					fmt.Fprintf(Global_ErrorManager.Out, "%+v\n", elkerr.Err)
					// errors.EORE(err)
				case "errstackfatal": // 打印错误链的调用栈，并结束程序
					fmt.Fprintf(Global_ErrorManager.Out, "%+v\n", elkerr.Err)
					// errors.EORE(err)
					elkerr.Cancel()
				default:
					fmt.Printf("only support \"debug warn errstack errstackfatal\" error type: %+v\n", elkerr.Err)
					os.Exit(1)
				}
			}
		}
	}(Global_ErrorManager.ErrCh)
}
