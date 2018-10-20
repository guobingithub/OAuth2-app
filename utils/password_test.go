package utils

import (
	"testing"
	"OAuth2-demo/logger"
	"crypto/md5"
	"io"
	"fmt"
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

func TestMd5(t *testing.T) {
	passwd := "abc123"

	//pwdSecret密文存入数据库
	pwdSecret := Md5(passwd)

	logger.Info("Md5加密成功, 密码:",pwdSecret)

	//根据username从数据库取出pwdSecret
	if Md5(passwd) != pwdSecret {
		logger.Error("校验密码是否有效, 密码错误!")
	}else {
		logger.Info("校验密码是否有效, 密码正确!")
	}

	pwdReSecret := Md5(Md5(passwd))
	logger.Info("Md5二次加密成功, 密码:",pwdReSecret)
}

func TestM(t *testing.T) {
	str := "abc123"

	//方法一
	data:=[]byte(str)
	has:= md5.Sum(data)
	md5str1 := fmt.Sprintf("%x",has) //将[]byte转成16进制

	fmt.Println(md5str1)

	//方法二
	w := md5.New()
	io.WriteString(w,str)              //将str写入到w中
	md5str2 := fmt.Sprintf("%x",w.Sum(nil)) //w.Sum(nil)将w的hash转成[]byte格式

	fmt.Println(md5str2)
}