package utils

import (
	"testing"
	"OAuth2-demo/logger"
)

func TestPassWord(t *testing.T) {
	passwd := "123456"

	pwdSecret,err := PasswordHash(passwd)
	if err != nil {
		logger.Error("加密密码失败, err:",err)
		return
	}

	logger.Info("加密密码成功, pwdSecret:",pwdSecret)
	ok, err := PasswordVerify(pwdSecret,passwd)
	if err != nil {
		logger.Error("校验密码是否有效, 校验出错.")
		return
	}

	if ok {
		logger.Info("校验密码是否有效, 密码正确!")
	}else {
		logger.Error("校验密码是否有效, 密码错误!")
	}
}