/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-ELK licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: Wangjunqi123 <wangjunqi@kylinos.cn>
 * Date: Mon Jun 24 09:21:47 2024 +0800
 */
package global

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

func FileReadString(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.Errorf("%s **errstack**0", err.Error())
	}

	return string(content), nil
}

func FileReadBytes(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Errorf("%s **errstack**0", err.Error())
	}
	defer f.Close()

	var content []byte
	readbuff := make([]byte, 1024*4)
	for {
		n, err := f.Read(readbuff)
		if err != nil {
			if err == io.EOF {
				if n != 0 {
					content = append(content, readbuff[:n]...)
				}
				break
			}
			return nil, errors.Errorf("%s **errstack**0", err.Error())
		}
		content = append(content, readbuff[:n]...)
	}

	return content, nil
}
