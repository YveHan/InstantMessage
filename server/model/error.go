package model

import (
	"errors"
)

var (
	ERROR_USER_NOTEXISTS = errors.New("用戶不存在")
	ERROR_USER_EXISTS = errors.New("用戶存在")
	ERROR_USER_PWD = errors.New("用戶密碼錯誤")
)

