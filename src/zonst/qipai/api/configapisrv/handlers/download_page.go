package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
	"zonst/logging"
	"zonst/qipai/api/configapisrv/config"
	"zonst/qipai/api/configapisrv/dependency/errs"
	"zonst/qipai/api/configapisrv/middlewares"
	"zonst/qipai/api/configapisrv/models"
	"zonst/qipai/api/configapisrv/remote_api"
	"zonst/qipai/api/configapisrv/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-xweb/log"
	"github.com/jinzhu/gorm"
)

// OnDownloadPageSelectReq 下载页面-详情
type OnDownloadPageSelectReq struct {
	GameID int `json:"game_id" binding:"required"`
}

// OnDownloadPageSelectRequest 下载页面-详情
func OnDownloadPageSelectRequest(c *gin.Context) {
	req := &OnDownloadPageSelectReq{}
	// userID := middlewares.GetToken(c).UserID
	// defer LogStatUserID("OnGameVersionDeleteRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnDownloadPageSelectRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}

	db := c.MustGet("qipaidb").(*gorm.DB)
	temp := &models.HtmlPageTable{}

	if err := db.Table("html_page").Where("game_id = ?", req.GameID).First(temp).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("OnDownloadPageSelectRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	log.Printf("temp: %v\n", temp)

	if err := db.Debug().Table("game_package").Where("game_id = ? and package_os = ?", req.GameID, "android").Select("package_url, concat_ws('-', package_version, package_url) as url, remark").Order("package_version desc, package_version_no desc, remark").Find(&temp.AndroidList).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("OnDownloadPageSelectRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	if err := db.Debug().Table("game_package").Where("game_id = ? and package_os = ?", req.GameID, "ios").Select("package_url, concat_ws('-', package_version, package_url) as url, remark").Order("package_version desc, package_version_no desc, remark").Find(&temp.IosList).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("OnDownloadPageSelectRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	// log.Printf("andriod_list: %v\n", andriodList)
	c.JSON(http.StatusOK, gin.H{"errno": "0", "data": temp})
}

// OnDownloadPageUpdateReq 下载页面-修改
type OnDownloadPageUpdateReq struct {
	ID       int32  `form:"id"`
	GameID   string `form:"game_id"`
	GameName string `form:"game_name"`
	EnName   string `form:"en_name"`
	// AndroidID         int    `form:"android_id"`
	AndroidUrl string `form:"android_url"`
	// IosID             int    `form:"ios_id"`
	IosUrl            string `form:"ios_url"`
	AdURL             string `form:"ad_url"`
	IsAdShow          bool   `form:"is_ad_show"`
	IsAndroidRedirect bool   `form:"is_android_redirect"`
	IsIosRedirect     bool   `form:"is_ios_redirect"`
}

// OnDownloadPageUpdateRequest 下载页面-修改
func OnDownloadPageUpdateRequest(c *gin.Context) {
	req := &OnDownloadPageUpdateReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnDownloadPageUpdateRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnDownloadPageUpdateRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}

	success, imgEixst := UploadImg(c, req.GameID)
	if !success {
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "上传图片出错"})
	}

	db := c.MustGet("qipaidb").(*gorm.DB)

	// andriodUrl := models.PackageURL{}
	// if req.AndroidUrl != "" {
	// 	andriodUrl.Url = req.AndroidUrl
	// } else {
	// 	if err := db.Table("game_package").Select("package_url").Where("id = ?", req.AndroidID).Find(&andriodUrl).Error; err != nil {
	// 		log.Errorf("OnDownloadPageUpdateRequest: err: %v\n", err.Error())
	// 		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
	// 		return
	// 	}
	// }
	// iosUrl := models.PackageURL{}
	// if req.IosUrl != "" {
	// 	iosUrl.Url = req.IosUrl
	// } else {
	// 	if err := db.Table("game_package").Select("package_url").Where("id = ?", req.IosID).Find(&iosUrl).Error; err != nil {
	// 		log.Errorf("OnDownloadPageUpdateRequest: err: %v\n", err.Error())
	// 		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
	// 		return
	// 	}
	// }
	config := c.MustGet("config").(*config.Config)

	imgURL := config.ImgURL
	html := "/data/static_site/www4417com/" + req.GameID + ".html"
	headImg := imgURL + "/game/" + req.GameID + "/" + req.GameID + "_head.png"
	adImg := imgURL + "/game/" + req.GameID + "/" + req.GameID + "_ad.png"
	bodyImg := imgURL + "/game/" + req.GameID + "/" + req.GameID + "_body.png"

	if req.ID == 0 {
		gameID, _ := strconv.Atoi(req.GameID)
		page := &models.HtmlPageTable{
			GameID:            gameID,
			GameName:          req.GameName,
			EnName:            req.EnName,
			AndroidUrl:        req.AndroidUrl,
			IosUrl:            req.IosUrl,
			AdUrl:             req.AdURL,
			IsAdShow:          req.IsAdShow,
			IsAndroidRedirect: req.IsAndroidRedirect,
			IsIosRedirect:     req.IsIosRedirect,
		}
		if imgEixst[0] {
			page.HeadImg = headImg
		}
		if imgEixst[1] {
			page.BodyImg = bodyImg
		}
		if imgEixst[2] {
			page.AdImg = adImg
		}
		if err := db.Table("html_page").Create(page).Error; err != nil {
			log.Errorf("OnDownloadPageUpdateRequest: Create: err: %v\n", err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
			return
		}
	} else {
		db = db.Table("html_page").Where("id = ?", req.ID).Update(req).Update("is_ad_show", req.IsAdShow).Update("is_android_redirect", req.IsAndroidRedirect).Update("is_ios_redirect", req.IsIosRedirect)
		if imgEixst[0] {
			db.Update("head_img", headImg)
		}
		if imgEixst[1] {
			db.Update("body_img", bodyImg)
		}
		if imgEixst[2] {
			db.Update("ad_img", adImg)
		}
		if err := db.Error; err != nil {
			log.Errorf("OnDownloadPageUpdateRequest: Update1: err: %v\n", err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
			return
		}
	}

	//
	//if err := db.Table("html_page").Where("id = ?", req.ID).Update("head_img", headImg).Update("ad_img", adImg).Update("body_img", bodyImg).Error; err != nil {
	//	log.Errorf("OnDownloadPageUpdateRequest: Update2: err: %v\n", err.Error())
	//}
	htmlURL := config.HTMLURL
	if !models.CreateHTML(c, html, req.GameName, req.IsAndroidRedirect, req.IsIosRedirect, req.IsAdShow, req.AdURL, adImg, req.AndroidUrl, headImg, bodyImg, req.IosUrl, req.GameID, htmlURL) {
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "生成html出错"})
		return
	}
	uri := htmlURL + req.GameID + ".html"
	response := utils.RefreshCDNURL(uri)
	if response != 0 {
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "刷新CDN失败"})
		return
	}
	if imgEixst[0] {
		response = utils.RefreshCDNURL(headImg)
		if response != 0 {
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "刷新CDN失败"})
			return
		}
	}
	if imgEixst[1] {
		response = utils.RefreshCDNURL(bodyImg)
		if response != 0 {
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "刷新CDN失败"})
			return
		}
	}
	if imgEixst[2] {
		response = utils.RefreshCDNURL(adImg)
		if response != 0 {
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "刷新CDN失败"})
			return
		}
	}
	//responce = middlewares.RefreshCDNDir(c, dir)
	//if responce != 0 {
	//	c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "刷新CDN失败"})
	//	return
	//}
	//responce = middlewares.RefreshCDNDir(c, dir2)
	//if responce != 0 {
	//	c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "刷新CDN失败"})
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{"errno": "0"})
}

// UploadImg UploadImg
func UploadImg(c *gin.Context, gameID string) (bool, []bool) {
	imgExist := make([]bool, 3, 3)
	headF, err := c.FormFile("head_img")
	if err != nil && err.Error() != "http: no such file" {
		log.Errorf("UploadImg: err: %v\n", err.Error())
		return false, imgExist
	}
	if err == nil {
		imgExist[0] = true
		path := c.MustGet("config").(*config.Config).ImgPath + gameID
		img := path + "/" + gameID + "_head.png"
		if !SaveImg(c, headF, path, img) {
			return false, imgExist
		}
	}

	bodyF, err := c.FormFile("body_img")
	if err != nil && err.Error() != "http: no such file" {
		log.Errorf("UploadImg: err: %v\n", err.Error())
		return false, imgExist
	}
	if err == nil {
		imgExist[1] = true
		path := c.MustGet("config").(*config.Config).ImgPath + gameID
		img := path + "/" + gameID + "_body.png"
		if !SaveImg(c, bodyF, path, img) {
			return false, imgExist
		}
	}

	adF, err := c.FormFile("ad_img")
	if err != nil && err.Error() != "http: no such file" {
		log.Errorf("UploadImg: err: %v\n", err.Error())
		return false, imgExist
	}
	if err == nil {
		imgExist[2] = true
		path := c.MustGet("config").(*config.Config).ImgPath + gameID
		img := path + "/" + gameID + "_ad.png"
		if !SaveImg(c, adF, path, img) {
			return false, imgExist
		}
	}
	return true, imgExist
}

// SaveImg SaveImg
func SaveImg(c *gin.Context, f *multipart.FileHeader, path string, img string) bool {
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Errorf("SaveImg: err: %v\n", err.Error())
		return false
	}
	os.Remove(img)
	if err := c.SaveUploadedFile(f, img); err != nil {
		log.Errorf("SaveImg: err: %v\n", err.Error())
		return false
	}
	return true
}

func OnDownloadPageOpRequest(c *gin.Context) {
	api := c.Request.RequestURI
	req := &models.OnDownloadPageOpRequest{}
	if err := c.Bind(req); err != nil {
		log.Errorf("api:%v err:%v\n", api, err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "参数错误"})
		return
	}

	db := c.MustGet("qipaidb").(*gorm.DB)
	gameList, err := models.GetGameIDListCloseStatus(db, []int32{int32(req.GameID)})
	if err != nil {
		log.Errorf("OnDownloadPageOpRequest: 查询平台是否关闭了下载页面出错, req:%+v, err:%v\n", req, err)
		c.JSON(http.StatusOK, gin.H{"errno": "2", "errmsg": "查询平台是否关闭了下载页面出错"})
		return
	}
	if len(gameList) == 0 {
		c.JSON(http.StatusOK, gin.H{"errno": "2", "errmsg": "平台配置不存在"})
		return
	}
	if gameList[0].Close == req.Close {
		c.JSON(http.StatusOK, gin.H{"errno": "0"})
		return
	}

	// 记录日志
	cfg := c.MustGet("config").(*config.Config)
	workerInfo := middlewares.GetToken(c)
	beforeContentBytes, err := json.Marshal(gameList[0])
	if err != nil {
		logging.Errorf("OnDownloadPageOpRequest: json序列化更新前的配置出错, req:%+v, err:%v\n", req, err)
		c.JSON(http.StatusOK, gin.H{"errno": "2", "errmsg": "json序列化更新前的配置出错"})
		return
	}
	afterContentBytes, err := json.Marshal(req)
	if err != nil {
		logging.Errorf("OnDownloadPageOpRequest: json序列化更新后的配置出错, req:%+v, err:%v\n", req, err)
		c.JSON(http.StatusOK, gin.H{"errno": "2", "errmsg": "json序列化更新后的配置出错"})
		return
	}
	if err = remote_api.AddConfigUpdateLog(int32(req.GameID), remote_api.ModuleIDHomePageConfig, workerInfo.UserID,
		remote_api.OpTypeAdd, string(beforeContentBytes), string(afterContentBytes), cfg.ConfigUpdateLogAPI,
		utils.ClubAppID, utils.ClubAppSecretKey); err != nil {
		logging.Errorf("OnDownloadPageOpRequest: 写入日志出错, req:%+v, err:%+v\n", req, err)
		c.JSON(http.StatusOK, gin.H{"errno": "2", "errmsg": "写入日志出错"})
		return
	}

	html := fmt.Sprintf("/data/static_site/www4417com/%v.html", req.GameID)
	htmlb := fmt.Sprintf("/data/static_site/www4417com/%v.html.backup", req.GameID)
	indexHTML := "/data/static_site/www4417com/index.html"

	if req.Close {
		cmd := exec.Command("cp", html, htmlb)
		if _, err := cmd.Output(); err != nil {
			errs.SaveError(errorx.Wrap(err), map[string]interface{}{
				"req": req,
			})
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "修改失败"})
			return
		}
		cmd = exec.Command("cp", indexHTML, html)
		if _, err := cmd.Output(); err != nil {
			errs.SaveError(errorx.Wrap(err), map[string]interface{}{
				"req": req,
			})
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "修改失败"})
			return
		}
		db := c.MustGet("qipaidb").(*gorm.DB)
		if err := db.Table("html_page").Where("game_id = ?", req.GameID).Update("close", true).Error; err != nil {
			log.Errorf("api:%v request:%v err:%v\n", api, req, err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "修改失败"})
			return
		}
	} else {
		cmd := exec.Command("cp", htmlb, html)
		if _, err := cmd.Output(); err != nil {
			errs.SaveError(errorx.Wrap(err), map[string]interface{}{
				"req": req,
			})
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "修改失败"})
			return
		}
		db := c.MustGet("qipaidb").(*gorm.DB)
		if err := db.Table("html_page").Where("game_id = ?", req.GameID).Update("close", false).Error; err != nil {
			log.Errorf("api:%v request:%v err:%v\n", api, req, err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "修改失败"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0"})
}

type GetGameListCloseStatusReq struct {
	GameIDList []int32 `json:"game_id_list" form:"game_id_list"`
}

func (c *GetGameListCloseStatusReq) Validate() error {
	if len(c.GameIDList) == 0 {
		return fmt.Errorf("game_id_list不能为空")
	}
	return nil
}

func GetGameListCloseStatus(c *gin.Context) {
	api := c.Request.RequestURI
	req := &GetGameListCloseStatusReq{}
	if err := c.ShouldBind(req); err != nil {
		log.Errorf("api:%v err:%v\n", api, err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "参数错误"})
		return
	}

	if err := req.Validate(); err != nil {
		log.Errorf("api:%v err:%v\n", api, err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}

	db := c.MustGet("qipaidb").(*gorm.DB)
	gameList, err := models.GetGameIDListCloseStatus(db, req.GameIDList)
	if err != nil {
		log.Errorf("api:%v request:%v err:%v\n", api, req, err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "查询平台列表是否关闭下载页面出错"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0", "data": gameList})
}

func OnHtmlUpdateRequest(c *gin.Context) {

	db := c.MustGet("qipaidb").(*gorm.DB)
	pages, err := models.GetAllPage(db)
	if err != nil {
		log.Errorf("OnHtmlUpdateRequest: GetAllPage: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}

	var failList []string
	config := c.MustGet("config").(*config.Config)

	imgURL := config.ImgURL
	htmlURL := config.HTMLURL

	backUpHtmlFormat := `/data/static_site/www4417com/%v.html.backup`
	htmlFormat := `/data/static_site/www4417com/%v.html`
	for _, v := range *pages {
		gameID := strconv.Itoa(v.GameID)
		html := fmt.Sprintf(htmlFormat, v.GameID)
		if v.Close {
			html = fmt.Sprintf(backUpHtmlFormat, v.GameID)
		}
		headImg := imgURL + "/game/" + gameID + "/" + gameID + "_head.png"
		adImg := imgURL + "/game/" + gameID + "/" + gameID + "_ad.png"
		bodyImg := imgURL + "/game/" + gameID + "/" + gameID + "_body.png"
		if !models.CreateHTML(c, html, v.GameName, v.IsAndroidRedirect, v.IsIosRedirect, v.IsAdShow, v.AdUrl, adImg, v.AndroidUrl, headImg, bodyImg, v.IosUrl, gameID, htmlURL) {
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "生成html出错"})
			return
		}

		uri := htmlURL + gameID + ".html"
		response := utils.RefreshCDNURL(uri)
		if response != 0 {
			//重试
			response = utils.RefreshCDNURL(uri)
			if response != 0 {
				failList = append(failList, gameID)
			}
		}

	}

	if len(failList) == 0 {
		c.JSON(http.StatusOK, gin.H{"errno": "0"})
	} else {
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": fmt.Sprintf("%v平台刷新cdn失败，请稍后手动重试", failList)})
	}

	return
}

// OnDownloadUrlRequest 下载页面-详情
func OnDownloadUrlRequest(c *gin.Context) {
	req := &OnDownloadPageSelectReq{}
	if err := c.Bind(req); err != nil {
		log.Errorf("OnDownloadUrlRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}

	db := c.MustGet("qipaidb").(*gorm.DB)
	temp := &models.HtmlPageTable{}

	if err := db.Table("html_page").Select("ios_url").Where("game_id = ?", req.GameID).First(temp).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("OnDownloadUrlRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0", "data": temp.IosUrl})
}

// OnDownloadPageUpdateURLRequest 单独刷新下载的ios和anroid地址接口
func OnDownloadPageUpdateURLRequest(c *gin.Context) {
	req := &OnDownloadPageUpdateReq{}
	userID := middlewares.GetToken(c).UserID
	defer LogStatUserID("OnDownloadPageUpdateRequest", c, userID, req, time.Now())
	if err := c.Bind(req); err != nil {
		log.Errorf("OnDownloadPageUpdateRequest: err: %v\n", err.Error())
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
		return
	}

	db := c.MustGet("qipaidb").(*gorm.DB)

	config := c.MustGet("config").(*config.Config)

	imgURL := config.ImgURL
	html := "/data/static_site/www4417com/" + req.GameID + ".html"
	headImg := imgURL + "/game/" + req.GameID + "/" + req.GameID + "_head.png"
	adImg := imgURL + "/game/" + req.GameID + "/" + req.GameID + "_ad.png"
	bodyImg := imgURL + "/game/" + req.GameID + "/" + req.GameID + "_body.png"

	if req.ID == 0 {
		gameID, _ := strconv.Atoi(req.GameID)
		page := &models.HtmlPageTable{
			GameID:            gameID,
			GameName:          req.GameName,
			EnName:            req.EnName,
			AndroidUrl:        req.AndroidUrl,
			IosUrl:            req.IosUrl,
			AdUrl:             req.AdURL,
			IsAdShow:          req.IsAdShow,
			IsAndroidRedirect: req.IsAndroidRedirect,
			IsIosRedirect:     req.IsIosRedirect,
		}
		if err := db.Table("html_page").Create(page).Error; err != nil {
			log.Errorf("OnDownloadPageUpdateRequest: Create: err: %v\n", err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
			return
		}
	} else {
		db = db.Table("html_page").Where("id = ?", req.ID).Update(req).Update("is_ad_show", req.IsAdShow).Update("is_android_redirect", req.IsAndroidRedirect).Update("is_ios_redirect", req.IsIosRedirect)
		if err := db.Error; err != nil {
			log.Errorf("OnDownloadPageUpdateRequest: Update1: err: %v\n", err.Error())
			c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": err.Error()})
			return
		}
	}

	htmlURL := config.HTMLURL
	if !models.CreateHTML(c, html, req.GameName, req.IsAndroidRedirect, req.IsIosRedirect, req.IsAdShow, req.AdURL, adImg, req.AndroidUrl, headImg, bodyImg, req.IosUrl, req.GameID, htmlURL) {
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "生成html出错"})
		return
	}
	uri := htmlURL + req.GameID + ".html"
	response := utils.RefreshCDNURL(uri)
	if response != 0 {
		c.JSON(http.StatusOK, gin.H{"errno": "-1", "errmsg": "刷新CDN失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errno": "0"})
}
