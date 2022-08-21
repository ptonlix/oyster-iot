package services

import (
	"os"
	"oyster-iot/init/constants"
	"path/filepath"

	"github.com/beego/beego/v2/core/logs"
)

type SysLogService struct {
}

func (s *SysLogService) GetSyslogFileList() (*[]string, error) {
	var files []string

	root := constants.LogFileDir
	err := filepath.Walk(root, s.visit(&files))
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	return &files, err
}

func (s *SysLogService) visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logs.Error(err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".log" {
			return nil
		}
		*files = append(*files, info.Name())
		return nil
	}
}

func (s *SysLogService) GetSyslogFile(fileName string) (*[]byte, error) {

	content, err := os.ReadFile(constants.LogFileDir + fileName)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	return &content, nil
}
