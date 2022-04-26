package utils

import "github.com/go-basic/uuid"

// 生成UUID
func GetUuid() string {
	uuid := uuid.New()
	return uuid
}
