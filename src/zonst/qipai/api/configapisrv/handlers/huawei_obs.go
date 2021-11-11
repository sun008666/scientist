package handlers

import (
	"io/ioutil"
	"mime/multipart"
	"zonst/qipai/api/configapisrv/config"

	"github.com/gin-gonic/gin"
	"github.com/go-xweb/log"

	//libhwobs "zonst/qipai/libhwobs/main"
	libhwobs "zonst/qipai/api/configapisrv/hwobs/main"
)

//UploadFileToHuaweiObs 上传文件到华为云环境存储桶
func UploadFileToHuaweiObsByFileName(localPath, fileName string) error {
	conf := config.Cfg
	//bucketName := conf.BucketName
	bucketName := conf.HuaweiBucketName
	ak := conf.Ak
	sk := conf.Sk
	endpoint := conf.Endpoint
	err := libhwobs.PutFile(bucketName, ak, sk, endpoint, fileName, localPath)
	if err != nil {
		log.Errorf("UploadToHuaweiObs: err:%v", err)
		return err
	}
	return nil
}

//UploadFileToHuaweiObsByMultipartFileHeader 上传文件到华为云环境存储桶
func UploadFileToHuaweiObsByMultipartFileHeader(f *multipart.FileHeader, fileName string) error {
	conf := config.Cfg
	//bucketName := conf.BucketName
	bucketName := conf.HuaweiBucketName
	ak := conf.Ak
	sk := conf.Sk
	endpoint := conf.Endpoint
	// 读取文件内容
	fileContent, err := f.Open()
	if err != nil {
		log.Errorf("UploadFileToHuaweiObsByMultipartFileHeader: Open: err:%v\n", err.Error())
		return err
	}
	fileBytes, err := ioutil.ReadAll(fileContent)
	if err != nil {
		log.Errorf("UploadFileToHuaweiObsByMultipartFileHeader: ioutil.ReadAll: err:%v\n", err.Error())
		return err
	}
	err = libhwobs.PutObject(bucketName, ak, sk, endpoint, fileName, fileBytes)
	return nil
}

//UploadObjectToHuaweiObs 上传对象到华为云环境存储桶
func UploadObjectToHuaweiObs(c *gin.Context, fileByte []byte, fileName string) error {
	conf := c.MustGet("config").(*config.Config)
	//bucketName := conf.BucketName
	bucketName := conf.HuaweiBucketName
	ak := conf.Ak
	sk := conf.Sk
	endpoint := conf.Endpoint
	err := libhwobs.PutObject(bucketName, ak, sk, endpoint, fileName, fileByte)
	return err
}
