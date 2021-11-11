package handlers

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"zonst/logging"
	"zonst/qipai/api/configapisrv/config"
	"zonst/qipai/api/configapisrv/dependency/db"
	"zonst/qipai/api/configapisrv/dependency/errs"
	"zonst/qipai/api/configapisrv/dependency/redistool"
	"zonst/qipai/api/configapisrv/middlewares"
	"zonst/qipai/api/configapisrv/models"
	"zonst/qipai/api/configapisrv/remote_api"
	"zonst/qipai/api/configapisrv/service"
	"zonst/qipai/api/configapisrv/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-xweb/log"
	"github.com/jinzhu/gorm"
	"github.com/shogo82148/androidbinary/apk"
)

const (
	IOSTYPE     = "ios"
	ANDIORDTYPE = "android"
)

type OnGameVersionUploadAndroidReq struct {
	GameID   int    `form:"game_id" binding:"required"`
	GameName string `form:"game_name" binding:"required"`
	Remark   string `form:"remark" binding:"required"`
}

// OnGameVersionUploadAndroidRequest 游戏版本上传-安卓
func OnGameVersionUploadAndroidRequest(c *gin.Context) {
	var uploadFunc func(*multipart.FileHeader, string) error
	environmentFlag := config.Cfg.EnvironmentFlag
	switch environmentFlag {
	case "tencent-cos":
		uploadFunc = CloudUploadFile1
	case "huawei-obs":
		uploadFunc = UploadFileToHuaweiObsByMultipartFileHeader
	default:
		log.Errorf("OnGameVersionUploadAndroidRequest: 云环境配置有误\n")
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "云环境配置有误"})
		return
	}

	req := &OnGameVersionUploadAndroidReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnGameVersionUploadAndroidRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnGameVersionUploadAndroidRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	f, err := c.FormFile("file")
	if err != nil {
		log.Errorf("OnGameVersionUploadAndroidRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	log.Printf("apk_name:%v", f.Filename)
	if f.Filename[len(f.Filename)-3:] != "apk" {
		log.Errorf("OnGameVersionUploadAndroidRequest: err: %v\n", "上传格式错误")
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "上传格式错误"})
		return
	}
	log.Printf("保存%v到临时文件temp.apk", f.Filename)
	if !SaveFile(c, req.GameID, f) {
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "保存apk包出错"})
		return
	}
	config := c.MustGet("config").(*config.Config)
	file := config.GameVersionPath + strconv.Itoa(req.GameID) + "/android/temp.apk"
	//获取版本号
	version, err := ParseApk(file)
	if err != nil {
		os.Remove(file)
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "解析apk包出错"})
		return
	}
	//计算MD5
	md5, err := FileMD5(file)
	if err != nil {
		os.Remove(file)
		log.Errorf("OnGameVersionUploadAndroidRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}

	db := c.MustGet("qipaidb").(*gorm.DB)
	// packageVersionNO := 0
	packageVersionNOObj := models.PackageVersionNO{}
	// db.Table("game_package").Where("game_id = ? and package_os = ? and package_version = ?", req.GameID, "android", version).Count(&packageVersionNO)
	if err := db.Debug().Table("game_package").Where("game_id = ? and package_os = ? and package_version = ?", req.GameID, "android", version).Select("max(package_version_no) as no").Find(&packageVersionNOObj).Error; err != nil {
		log.Errorf("OnGameVersionUploadAndroidRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	packageVersionNO := packageVersionNOObj.NO + 1
	log.Printf("packageVersionNO:%v", packageVersionNO)
	gameName := ChineseToRune(req.GameName) + "mj"
	packageURL := config.AndroidPackageURL + strconv.Itoa(req.GameID) + "/android/" + gameName + "_" + version + "." + strconv.Itoa(packageVersionNO) + ".apk"

	obj := &models.GamePackageTable{
		GameID:           req.GameID,
		GameName:         "中至" + req.GameName,
		PackageVersion:   version,
		PackageVersionNO: packageVersionNO,
		PackageOS:        "android",
		PackageURL:       packageURL,
		PackageMD5:       md5,
		PackageSize:      f.Size,
		LogTime:          time.Now().Format("2006-01-02 15:04:05"),
		Remark:           req.Remark,
	}

	apkName := config.GameVersionPath + strconv.Itoa(req.GameID) + "/android/" + gameName + "_" + version + "." + strconv.Itoa(packageVersionNO) + ".apk"
	err = os.Rename(file, apkName)
	if err != nil {
		os.Remove(file)
		log.Errorf("OnGameVersionUploadAndroidRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}

	savedCloud := "/" + strconv.Itoa(req.GameID) + "/android/" + gameName + "_" + version + "." + strconv.Itoa(packageVersionNO) + ".apk"
	//if err := CloudUploadFile1(c, f, savedCloud); err != nil {
	//	os.Remove(apkName)
	//	log.Errorf("OnGameVersionUploadAndroidRequest: CloudUploadFile1: err: %v\n", err.Error())
	//	c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "apk上传至云失败"})
	//	return
	//}
	if err := uploadFunc(f, savedCloud); err != nil {
		os.Remove(apkName)
		log.Errorf("OnGameVersionUploadAndroidRequest: uploadFunc: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "apk上传至云失败"})
		return
	}

	log.Printf("OnGameVersionUploadAndroidRequest: obj: %+v\n", obj)
	if err := db.Debug().Table("game_package").Create(obj).Error; err != nil {
		os.Remove(apkName)
		log.Errorf("OnGameVersionUploadAndroidRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	response := utils.RefreshCDNURL(config.CDNAPI + savedCloud)
	if response != 0 {
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "刷新CDN失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0"})
}

func SaveFile(c *gin.Context, gameID int, f *multipart.FileHeader) bool {
	gameVersionPath := c.MustGet("config").(*config.Config).GameVersionPath + strconv.Itoa(gameID) + "/android"
	if err := os.MkdirAll(gameVersionPath, 0755); err != nil {
		log.Errorf("SaveFile: err: %v\n", err.Error())
		return false
	}

	saved := gameVersionPath + "/temp.apk"
	log.Printf("temp.apk路径:%v", saved)
	if _, err := os.Stat(saved); err == nil {
		log.Errorf("SaveFile: err: %v\n", "文件已存在")
		return false
	}

	if err := c.SaveUploadedFile(f, saved); err != nil {
		log.Errorf("SaveFile: err: %v\n", err.Error())
		return false
	}
	return true
}

func ParseApk(file string) (string, error) {
	pkg, err := apk.OpenFile(file)
	defer pkg.Close()
	if err != nil {
		log.Errorf("ParseApk: err: %v", err.Error())
		return "", err
	}
	return pkg.Manifest().VersionName.String()
}

func FileMD5(file string) (fileMd5 string, err error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	value := md5.Sum(f)
	return fmt.Sprintf("%x", value), err
}

func ChineseToRune(str string) string {
	var temp = []rune(str)
	var obj string
	for i := 0; i < len(temp)-2; i++ {
		obj = obj + ChineseToStr(temp[i])
	}
	return obj
}

func ChineseToStr(chinese rune) string {
	var py utils.Pinyin

	// 初始化，载入汉字拼音映射文件
	py.Init("pinyin_table.txt")

	return py.GetPinyin(chinese, false)
}

type OnGameVersionUploadIosReq struct {
	GameID   int    `form:"game_id" binding:"required"`
	GameName string `form:"game_name" binding:"required"`
	Remark   string `form:"remark" binding:"required"`
}

// OnGameVersionUploadIosRequest 游戏版本上传 ios
func OnGameVersionUploadIosRequest(c *gin.Context) {
	logging.Debugf("OnGameVersionUploadIosRequest,开始解析参数\n")

	req := &OnGameVersionUploadIosReq{}
	// userID := middlewares.GetToken(c).UserID
	// defer LogStatUserID("OnGameVersionUploadIosRequest", c, userID, req, time.Now())
	if e := c.Bind(req); e != nil {
		errs.SaveError(errorx.Wrap(e), map[string]interface{}{
			"req": req,
		})
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": e.Error()})
		return
	}

	conn := redistool.WhiteListCacheRedisPool.Get()
	defer conn.Close()

	var lockkey = fmt.Sprintf("configapisrv:uploadioslock:%s:%d:%s:%s", config.Node.Mode, req.GameID, req.GameName, req.Remark)

	if !redistool.Once(
		conn,
		lockkey,
		6*60,
	) {
		c.JSON(200, gin.H{
			"errmsg": fmt.Sprintf("上一个任务进行中,请勿重复执行。(game_i=%d,game_name=%s,remark=%s)", req.GameID, req.GameName, req.Remark),
			"errno":  "-1",
		})
		return
	}

	fr, fh, e := c.Request.FormFile("file")
	if e != nil {
		errs.SaveError(errorx.Wrap(e), map[string]interface{}{
			"req": req,
		})

		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": e.Error()})
		return
	}

	if fh.Filename[len(fh.Filename)-3:] != "ipa" {
		errs.SaveError(errorx.NewServiceError(`上传格式错误 'f.Filename[len(f.Filename)-3:] != "ipa"'`, 1), map[string]interface{}{
			"req": req,
		})
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "上传格式错误"})
		return
	}

	if e := SaveIPAFile(c, req.GameID, fh); e != nil {

		errs.SaveError(errorx.Wrap(e), map[string]interface{}{
			"req": req,
		})
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "保存ipa包出错"})
		return
	}

	go func() {
		defer func() {
			if e := recover(); e != nil {
				errs.SaveError(errorx.NewFromStringf("panic recovers from %v", e))
			}
		}()

		defer func() {
			conn := redistool.WhiteListCacheRedisPool.Get()
			conn.Do("del", lockkey)
		}()

		var ul = models.UploadLog{
			Label:     "ios",
			RequestId: fmt.Sprintf("%s:%d:%d", req.GameName, req.GameID, time.Now().UnixNano()),
			GameName:  req.GameName,
			GameId:    req.GameID,
			Remark:    req.Remark,
			Md5:       "",
		}
		if e := db.QipaiDB.Model(models.UploadLog{}).Create(&ul).Error; e != nil {
			errs.SaveError(errorx.Wrap(e))
			return
		}

		tc := service.NewTrace(ul.RequestId)

		if e := service.GameVersionUploadIosRequest(fr, fh.Size, req.GameID, req.GameName, req.Remark, tc); e != nil {
			errs.SaveError(errorx.Wrap(e))

			ul.DB().Model(ul).Where("id=?", ul.Id).Updates(map[string]interface{}{
				"vstate":      3,
				"fail_reason": errorx.Wrap(e).Error(),
				"finish_at":   time.Now(),
			})

			return
		}

		ul.DB().Model(ul).Where("id=?", ul.Id).Updates(map[string]interface{}{
			"vstate":    2,
			"finish_at": time.Now(),
		})
	}()

	c.JSON(http.StatusOK, gin.H{"errno": "0"})

}

// SaveIPAFile ios包保存到服务器
func SaveIPAFile(c *gin.Context, gameID int, f *multipart.FileHeader) error {
	gameVersionPath := config.Cfg.GameVersionPath + strconv.Itoa(gameID) + "/ios"
	if e := os.MkdirAll(gameVersionPath, 0755); e != nil {
		return errorx.Wrap(e)
	}

	saved := gameVersionPath + "/temp.ipa"
	os.Remove(saved)
	if _, err := os.Stat(saved); err == nil {
		return errorx.NewServiceError("文件已存在", 101)
	}

	if e := c.SaveUploadedFile(f, saved); e != nil {
		return errorx.Wrap(e)
	}
	return nil
}

type OnGameVersionUploadUrlReq struct {
	GameID   int    `json:"game_id" binding:"required"`
	GameName string `json:"game_name" binding:"required"`
	IosUrl   string `json:"ios_url" binding:"required"`
}

// OnGameVersionUploadUrlRequest 上传ios url下载地址
func OnGameVersionUploadUrlRequest(c *gin.Context) {
	req := &OnGameVersionUploadUrlReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnGameVersionUploadUrlRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnGameVersionUploadUrlRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	obj := &models.GamePackageTable{
		GameID:         req.GameID,
		GameName:       req.GameName,
		PackageVersion: "1.0",
		PackageOS:      "ios",
		PackageURL:     req.IosUrl,
		LogTime:        time.Now().Format("2006-01-02 15:04:05"),
	}
	log.Printf("OnGameVersionUploadUrlRequest: obj: %+v\n", obj)
	db := c.MustGet("qipaidb").(*gorm.DB)
	n := 0
	if db.Table("game_package").Where("game_name = ? and package_size = ?", req.GameName, 0).Count(&n); n != 0 {
		if err := db.Table("game_package").Where("game_name = ? and package_size = ?", req.GameName, 0).Update(obj).Error; err != nil {
			log.Errorf("OnGameVersionUploadUrlRequest: err: %v\n", err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
			return
		}
	} else {
		if err := db.Table("game_package").Create(obj).Error; err != nil {
			log.Errorf("OnGameVersionUploadUrlRequest: err: %v\n", err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0"})
}

type OnGameVersionListReq struct {
	GameID      int    `json:"game_id" form:"game_id" binding:"required"`
	PackageType string `json:"package_type" form:"package_type"`
}

func OnGameVersionListRequest(c *gin.Context) {
	req := &OnGameVersionListReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnGameVersionListRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnGameVersionListRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	db := c.MustGet("qipaidb").(*gorm.DB)
	obj := []models.GamePackageTable{}
	// if err := db.Table("game_package").Where("game_id = ?", req.GameID).Order("package_version desc, package_version_no desc, package_os").Find(&obj).Error; err != nil && err != gorm.ErrRecordNotFound {
	sql := db.Table("game_package").Where("game_id = ?", req.GameID)
	if req.PackageType != "" {
		if req.PackageType != IOSTYPE && req.PackageType != ANDIORDTYPE {
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "包的类型只能是ios或者android"})
			return
		}
		sql = sql.Where("package_os=?", req.PackageType)
	}
	if err := sql.Order("log_time desc").Find(&obj).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("OnGameVersionListRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	data := []models.OnGameVersionListResponce{}
	for _, v := range obj {
		t, _ := time.Parse("2006-01-02T15:04:05Z", v.LogTime)
		temp := models.OnGameVersionListResponce{
			ID:          v.ID,
			GameID:      v.GameID,
			GameName:    v.GameName,
			Version:     fmt.Sprintf("%v.%v", v.PackageVersion, v.PackageVersionNO),
			PackageOS:   v.PackageOS,
			PackageMD5:  v.PackageMD5,
			PackageSize: fmt.Sprintf("%.2fM", float64(v.PackageSize)/1024/1024),
			LogTime:     t.Format("2006-01-02 15:04:05"),
			Remark:      v.Remark,
			Status:      v.Status,
		}
		data = append(data, temp)
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0", "data": data})
}

type OnGameVersionDeleteReq struct {
	ID int `json:"id" binding:"required"`
}

func OnGameVersionDeleteRequest(c *gin.Context) {
	req := &OnGameVersionDeleteReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnGameVersionDeleteRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnGameVersionDeleteRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	db := c.MustGet("qipaidb").(*gorm.DB)
	if err := db.Table("game_package").Where("id = ?", req.ID).Delete(&models.GamePackageTable{}).Error; err != nil {
		log.Errorf("OnGameVersionDeleteRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0"})
}

type OnUpdateGameVersionStatusReq struct {
	OpenPackageID  int32 `json:"open_package_id" form:"open_package_id"`   //要开启的包
	ClosePackageID int32 `json:"close_package_id" form:"close_package_id"` //要关闭的包
}

//加一个接口设置上传包状态的设置，哪个包的状态要开，哪个包的状态要关，
func OnUpdateGameVersionStatusRequest(c *gin.Context) {
	req := &OnUpdateGameVersionStatusReq{}
	userID := middlewares.GetToken(c).UserID
	jwtToken := c.Request.Header.Get("x-xq5-jwt")
	defer LogStatUserID("OnUpdateGameVersionStatusRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnUpdateGameVersionStatusRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "参数不匹配," + err.Error()})
		return
	}
	if req.ClosePackageID <= 0 && req.OpenPackageID <= 0 {
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "参数不匹配开启和关闭的包必须有一个"})
		return
	}
	qipaidb := c.MustGet("qipaidb").(*gorm.DB)
	//开启某一个包
	if req.OpenPackageID > 0 {
		//获取要开启的包的信息
		openGameVersion, err := models.GetGameVersionPackageByID(qipaidb, req.OpenPackageID)
		if err != nil {
			log.Errorf("OnUpdateGameVersionStatusRequest: err: %v\n", err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
			return
		}
		htmlConfig, err := models.GetHtmlPageByGameID(qipaidb, openGameVersion.GameID)
		if err != nil {
			log.Errorf("OnUpdateGameVersionStatusRequest-err: %v\n", err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
			return
		}
		//先修改包的状态
		err = models.UpdateGameVersionPackageStatus(qipaidb, req.OpenPackageID, models.Openstatus)
		if err != nil {
			log.Errorf("OnUpdateGameVersionStatusRequest: err: %v\n", err.Error())
			return
		}
		//更新下载地址
		//package_os  ios android 判断安卓还是ios
		if openGameVersion.PackageOS == "ios" {
			//要开启的是唯一一个地址
			IosUrlArr := make([]string, 0)
			if strings.Contains(htmlConfig.IosUrl, "|") {
				IosUrlArr = strings.Split(htmlConfig.IosUrl, "|")
				IosUrlArr = append(IosUrlArr, openGameVersion.PackageURL)
			} else {
				if len(htmlConfig.IosUrl) > 0 {
					IosUrlArr = append(IosUrlArr, htmlConfig.IosUrl)
				}
				IosUrlArr = append(IosUrlArr, openGameVersion.PackageURL)
			}
			htmlConfig.IosUrl = DistinctURL(IosUrlArr)

		} else {
			//安卓包
			androidUrl := make([]string, 0)
			if strings.Contains(htmlConfig.AndroidUrl, "|") {
				androidUrl = strings.Split(htmlConfig.AndroidUrl, "|")
			} else {
				if len(htmlConfig.AndroidUrl) > 0 {
					androidUrl = append(androidUrl, htmlConfig.AndroidUrl)
				}
			}
			androidUrl = append(androidUrl, openGameVersion.PackageURL)
			htmlConfig.AndroidUrl = DistinctURL(androidUrl)
		}
		//然后在调用修改下载页面的接口，将ios或者安卓包的下载地址刷新一下
		err = remote_api.UpdateDownloadHtmlURL(htmlConfig.ID, htmlConfig.GameID, htmlConfig.GameName, htmlConfig.EnName, htmlConfig.AndroidUrl, htmlConfig.IosUrl, htmlConfig.AdUrl, jwtToken, htmlConfig.IsAdShow, htmlConfig.IsAndroidRedirect, htmlConfig.IsIosRedirect)
		if err != nil {
			//修改不成功就改回来状态 对调一下参数
			err1 := models.UpdateGameVersionPackageStatus(qipaidb, req.OpenPackageID, models.CloseStatus)
			if err1 != nil {
				log.Errorf("OnUpdateGameVersionStatusRequest-err: %v\n", err1.Error())
			}
			log.Errorf("OnUpdateGameVersionStatusRequest-err: %v\n", err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"errno": "0", "errmsg": "开启成功"})
		return
	}
	//关闭某一个包
	err := closePackage(qipaidb, jwtToken, req)
	if err != nil {
		log.Errorf("OnUpdateGameVersionStatusRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0", "errmsg": "关闭成功"})
	return
}

func closePackage(qipaidb *gorm.DB, jwtToken string, req *OnUpdateGameVersionStatusReq) error {
	//获取要关闭的包的信息
	openGameVersion, err := models.GetGameVersionPackageByID(qipaidb, req.ClosePackageID)
	if err != nil {
		log.Errorf("closePackage-err: %v\n", err.Error())
		return err
	}
	//先修改包的状态
	err = models.UpdateGameVersionPackageStatus(qipaidb, req.ClosePackageID, models.CloseStatus)
	if err != nil {
		log.Errorf("closePackage-err: %v\n", err.Error())
		return err
	}

	htmlConfig, err := models.GetHtmlPageByGameID(qipaidb, openGameVersion.GameID)
	if err != nil {
		log.Errorf("closePackage-err: %v\n", err.Error())
		return err
	}
	//找到之前的包的地址,删除掉
	//package_os  ios android 判断安卓还是ios
	if openGameVersion.PackageOS == "ios" {
		if !strings.Contains(htmlConfig.IosUrl, "|") {
			if htmlConfig.IosUrl == openGameVersion.PackageURL {
				return errors.New("必须要有一个开启状态的更新包")
			}
		}
		IosUrlArr := strings.Split(htmlConfig.IosUrl, "|")
		htmlConfig.IosUrl = DeleteOneURL(IosUrlArr, openGameVersion.PackageURL)
	} else {
		//要关闭的是唯一一个地址
		if !strings.Contains(htmlConfig.AndroidUrl, "|") {
			if htmlConfig.AndroidUrl == openGameVersion.PackageURL {
				return errors.New("必须要有一个开启状态的更新包")
			}
			return errors.New("当前要关闭的更新包未开启")
		}
		androidURLArr := strings.Split(htmlConfig.AndroidUrl, "|")
		htmlConfig.AndroidUrl = DeleteOneURL(androidURLArr, openGameVersion.PackageURL)
	}

	//然后在调用修改下载页面的接口，将ios或者安卓包的下载地址刷新一下
	err = remote_api.UpdateDownloadHtmlURL(htmlConfig.ID, htmlConfig.GameID, htmlConfig.GameName, htmlConfig.EnName, htmlConfig.AndroidUrl, htmlConfig.IosUrl, htmlConfig.AdUrl, jwtToken, htmlConfig.IsAdShow, htmlConfig.IsAndroidRedirect, htmlConfig.IsIosRedirect)
	if err != nil {
		//修改不成功就改回来状态 对调一下参数
		err1 := models.UpdateGameVersionPackageStatus(qipaidb, req.ClosePackageID, models.Openstatus)
		if err1 != nil {
			log.Errorf("closePackage-err: %v\n", err1.Error())
		}
		return err
	}
	return nil
}

//去重url
func DistinctURL(arr []string) string {
	if len(arr) == 0 {
		return ""
	}
	if len(arr) == 1 {
		return arr[0]
	}
	newArr := make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	if len(newArr) == 1 {
		return newArr[0]
	}
	result := ""
	arrLen := len(newArr) - 1
	for i, v := range newArr {
		if i != arrLen {
			result += v + "|"
		} else {
			result += v
		}
	}
	return result
}

//删除某个url
func DeleteOneURL(arr []string, oneURL string) string {
	if len(arr) == 0 {
		return ""
	}
	if len(arr) == 1 {
		if arr[0] == oneURL {
			return ""
		}
		return arr[0]
	}
	newArr := make([]string, 0)
	oldLen := len(arr)
	for i := 0; i < oldLen; i++ {
		//删除当前关闭的
		if arr[i] == oneURL {
			continue
		}
		newArr = append(newArr, arr[i])
	}
	if len(newArr) == 1 {
		return newArr[0]
	}
	result := ""
	arrLen := len(newArr) - 1
	for i, v := range newArr {
		if i != arrLen {
			result += v + "|"
		} else {
			result += v
		}
	}
	return result
}
