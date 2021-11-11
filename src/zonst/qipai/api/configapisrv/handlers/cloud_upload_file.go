package handlers

import (
	"bytes"
	"context"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"zonst/qipai/api/configapisrv/config"
	"zonst/qipai/libcos"

	"github.com/go-xweb/log"
)

func CloudUploadFile1(f *multipart.FileHeader, fileName string) error {
	tomlConfig := config.Cfg
	bucketName := tomlConfig.BucketName
	appID := tomlConfig.AppID
	secretID := tomlConfig.SecretID
	secretKey := tomlConfig.SecretKey
	u, _ := url.Parse("https://" + bucketName + "-" + appID + ".cos.ap-shanghai.myqcloud.com")
	client := libcos.NewClient1(u, appID, secretID, secretKey, 10)

	// 读取文件内容
	fileContent, err := f.Open()
	if err != nil {
		log.Errorf("CloudUploadFile1: Open: err:%v\n", err.Error())
		return err
	}
	fi, err := ioutil.ReadAll(fileContent)
	if err != nil {
		log.Errorf("CloudUploadFile1: ioutil.ReadAll: err:%v\n", err.Error())
		return err
	}
	fReader := bytes.NewReader(fi)

	// 上传
	err = client.Put(context.Background(), fileName, fReader)
	if err != nil {
		log.Errorf("CloudUploadFile1: Put: err:%+v\n", err)
		return err
	}
	return nil
}
