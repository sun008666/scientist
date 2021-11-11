package service

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/go-xweb/log"
	"github.com/phinexdaz/ipapk"
	"howett.net/plist"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"zonst/qipai/api/configapisrv/config"
	db2 "zonst/qipai/api/configapisrv/dependency/db"
	"zonst/qipai/api/configapisrv/hwobs/main"
	"zonst/qipai/api/configapisrv/models"
	"zonst/qipai/api/configapisrv/utils"
	"zonst/qipai/libcos"
)

func GameVersionUploadIosRequest(fr io.Reader, size int64, gameId int, gameName string, remark string, tc *Trace) error {
	var uploadFuncByMultipartFileHeader func(io.Reader, string) error
	var uploadFuncByFileName func(string, string) error

	tc.StepStart("读取云环境")

	environmentFlag := config.Cfg.EnvironmentFlag
	switch environmentFlag {
	case "tencent-cos":
		uploadFuncByMultipartFileHeader = CloudUploadFile1
		uploadFuncByFileName = CloudUploadFile3
	case "huawei-obs":
		uploadFuncByMultipartFileHeader = UploadFileToHuaweiObsByMultipartFileHeader
		uploadFuncByFileName = UploadFileToHuaweiObsByFileName
	default:
		tc.StepFail("读取云环境", errorx.NewServiceError("云环境配置有误", 2))
		return errorx.NewServiceError("云环境配置有误", 2)
	}

	tc.StepOver("读取云环境")

	tc.StepStart("ipa转换")
	file := config.Cfg.GameVersionPath + strconv.Itoa(gameId) + "/ios/temp.ipa"

	ipa, e := ipapk.NewAppParser(file)
	if e != nil {
		os.Remove(file)

		tc.StepFail("ipa转换", errorx.Wrap(e))
		return errorx.Wrap(e)
	}
	tc.StepOver("ipa转换")

	tc.StepStart("计算md5")
	version := ipa.Build
	name := ipa.Name
	displayName := ipa.BundleId

	log.Printf("OnGameVersionUploadIosRequest: version:%v name:%v displayName:%v\n", version, name, displayName)
	//计算MD5
	md5, e := fileMD5(file)
	if e != nil {
		os.Remove(file)

		tc.StepFail("计算md5", errorx.Wrap(e))
		return errorx.Wrap(e)
	}

	tc.StepOver("计算md5")

	tc.StepStart("读取游戏包记录")

	db := db2.QipaiDB
	packageVersionNOObj := models.PackageVersionNO{}
	if e := db.Debug().Table("game_package").Where("game_id = ? and package_os = ? and package_version = ?", gameId, "ios", version).Select("max(package_version_no) as no").Find(&packageVersionNOObj).Error; e != nil {
		tc.StepFail("读取游戏包记录", errorx.Wrap(e))

		return errorx.Wrap(e)
	}

	tc.StepOver("读取游戏包记录")

	packageVersionNO := packageVersionNOObj.NO + 1
	log.Printf("packageVersionNO:%v", packageVersionNO)

	gameName = chineseToRune(gameName) + "mj"

	packageURL := config.Cfg.IosPackageURL + strconv.Itoa(gameId) + "/ios/" + gameName + "_" + version + "." + strconv.Itoa(packageVersionNO) + ".plist"

	packageIcon := config.Cfg.PackageIconURL + strconv.Itoa(gameId) + "/icons/" + version + "-icon.png"

	obj := &models.GamePackageTable{
		GameID:           gameId,
		GameName:         "中至" + gameName,
		PackageVersion:   version,
		PackageVersionNO: packageVersionNO,
		PackageOS:        "ios",
		PackageURL:       packageURL,
		PackageName:      name,
		PackageIcon:      packageIcon,
		PackageMD5:       md5,
		PackageSize:      size,
		LogTime:          time.Now().Format("2006-01-02 15:04:05"),
		Remark:           remark,
	}

	tc.StepStart("ipa重命名")

	ipaName := config.Cfg.GameVersionPath + strconv.Itoa(gameId) + "/ios/" + gameName + "_" + version + "." + strconv.Itoa(packageVersionNO) + ".ipa"
	e = os.Rename(file, ipaName)
	if e != nil {
		os.Remove(file)
		tc.StepFail("ipa重命名", errorx.Wrap(e))

		return errorx.Wrap(e)
	}
	tc.StepOver("ipa重命名")

	iconPath := config.Cfg.GameVersionPath + strconv.Itoa(gameId) + "/icons/"
	icon := iconPath + version + "-icon.png"

	tc.StepStart("ipa重命名")

	if !SaveIconImage(ipaName, iconPath, icon) {
		os.Remove(ipaName)

		tc.StepFail("ipa重命名", errorx.NewServiceError("拷贝Icon-57.png或AppIcon57x57.png文件失败,请检查包中是否存在该文件", 10))
		return errorx.NewServiceError("拷贝Icon-57.png或AppIcon57x57.png文件失败,请检查包中是否存在该文件", 10)
	}
	tc.StepOver("ipa重命名")

	tc.StepStart("创建plist")

	plistName := config.Cfg.GameVersionPath + strconv.Itoa(gameId) + "/ios/" + gameName + "_" + version + "." + strconv.Itoa(packageVersionNO) + ".plist"
	ipaURL := config.Cfg.PackageIconURL + strconv.Itoa(gameId) + "/ios/" + gameName + "_" + version + "." + strconv.Itoa(packageVersionNO) + ".ipa"
	if !CreatePlist(plistName, ipaURL, packageIcon, name, version, displayName) {
		os.Remove(ipaName)
		tc.StepFail("创建plist", errorx.NewServiceError("生成plist文件出错", 10))

		return errorx.NewServiceError("生成plist文件出错", 10)
	}
	tc.StepOver("创建plist")

	tc.StepStart("上传ipa")

	savedIpaCloud := "/" + strconv.Itoa(gameId) + "/ios/" + gameName + "_" + version + "." + strconv.Itoa(packageVersionNO) + ".ipa"

	if e := uploadFuncByMultipartFileHeader(fr, savedIpaCloud); e != nil {
		os.Remove(ipaName)
		tc.StepFail("上传ipa", errorx.Wrap(e))
		return errorx.Wrap(e)
	}

	tc.StepOver("上传ipa")

	saveIconCLoud := "/" + strconv.Itoa(gameId) + "/icons/" + version + "-icon.png"

	tc.StepStart("上传icon")

	if e := uploadFuncByFileName(icon, saveIconCLoud); e != nil {
		os.Remove(ipaName)
		tc.StepFail("上传icon", errorx.Wrap(e))

		return errorx.Wrap(e)
	}
	tc.StepOver("上传icon")

	savePlistCLoud := strings.Replace(savedIpaCloud, ".ipa", ".plist", 1)

	tc.StepStart("上传plist")

	if e = uploadFuncByFileName(plistName, savePlistCLoud); e != nil {
		os.Remove(ipaName)
		tc.StepFail("上传plist", errorx.Wrap(e))

		return errorx.Wrap(e)
	}
	tc.StepOver("上传plist")

	tc.StepStart("创建游戏包记录")

	if e := db.Debug().Table("game_package").Create(obj).Error; e != nil {
		os.Remove(ipaName)
		tc.StepFail("创建游戏包记录", errorx.Wrap(e))

		return errorx.Wrap(e)
	}

	tc.StepOver("创建游戏包记录")

	tc.StepStart("刷新ipaCDN")

	cfg := config.Cfg
	response := utils.RefreshCDNURL(cfg.CDNAPI + savedIpaCloud)
	if response != 0 {
		tc.StepFail("刷新ipaCDN", errorx.NewServiceError("刷新CDNIPA失败", 5))

		return errorx.NewServiceError("刷新CDNIPA失败", 5)
	}
	tc.StepOver("刷新ipaCDN")

	tc.StepStart("刷新iconCDN")

	response = utils.RefreshCDNURL(cfg.CDNAPI + saveIconCLoud)
	if response != 0 {
		tc.StepFail("刷新iconCDN", errorx.NewServiceError("刷新CDNIcon失败", 5))

		return errorx.NewServiceError("刷新CDNIcon失败", 5)
	}

	tc.StepOver("刷新iconCDN")

	tc.StepStart("刷新plistCDN")
	response = utils.RefreshCDNURL(cfg.CDNAPI + savePlistCLoud)
	if response != 0 {
		tc.StepFail("刷新plistCDN", errorx.NewServiceError("刷新CDNplist失败", 5))

		return errorx.NewServiceError("刷新CDNplist失败", 5)
	}

	tc.StepOver("刷新plistCDN")

	return nil
}

//ParseIpa : It parses the given ipa and returns a map from the contents of Info.plist in it
func ParseIpa(name string) (map[string]interface{}, error) {
	r, err := zip.OpenReader(name)
	if err != nil {
		log.Println("Error opening ipa/zip ", err.Error())
		return nil, err
	}
	defer r.Close()

	for _, file := range r.File {
		if file.FileInfo().Name() == "Info.plist" {
			rc, err := file.Open()
			if err != nil {
				log.Println("Error opening Info.plist in zip", err.Error())
				return nil, err
			}
			buf := make([]byte, file.FileInfo().Size())
			_, err = io.ReadFull(rc, buf)
			if err != nil {
				log.Println("Error reading Info.plist", err.Error())
				return nil, err
			}
			log.Printf("info_map:%v", string(buf))
			var info_map map[string]interface{}
			_, err = plist.Unmarshal(buf, &info_map)
			if err != nil {
				log.Println("Error reading Info.plist", err.Error())
				return nil, err
			}
			return info_map, nil
		}
	}
	return nil, errors.New("Info.plist not found")
}

// SaveIconImage 保存ios的icon图片
func SaveIconImage(ipaFile string, imgPath string, imgFile string) bool {
	if err := os.MkdirAll(imgPath, 0755); err != nil {
		log.Errorf("SaveIconImage: err: %v\n", err.Error())
		return false
	}
	r, err := zip.OpenReader(ipaFile)
	if err != nil {
		log.Errorf("SaveIconImage: err: %v\n", err.Error())
		return false
	}
	defer r.Close()

	for _, file := range r.File {
		fileName := file.FileInfo().Name()
		if fileName == "Icon-57.png" || fileName == "AppIcon57x57.png" {
			rc, err := file.Open()
			if err != nil {
				log.Errorf("SaveIconImage: err: %v\n", err.Error())
				return false
			}
			if _, err = os.Stat(imgFile); err == nil {
				os.Remove(imgFile)
			}
			f, err := os.Create(imgFile)
			if err != nil {
				log.Errorf("SaveIconImage: err: %v\n", err.Error())
				return false
			}
			defer f.Close()
			_, err = io.Copy(f, rc)
			if err != nil {
				log.Errorf("SaveIconImage: err: %v\n", err.Error())
				return false
			}
			return true
		}
	}
	return false
}

// CreatePlist 生成plist文件
func CreatePlist(plist string, ipaURL string, iconURL string, packageName string, version string, gameName string) bool {
	temp := `<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
	<dict>
		<key>items</key>
		<array>
			<dict>
				<key>assets</key>
				<array>
					<dict>
						<key>kind</key>
						<string>software-package</string>
						<key>url</key>
						<string><![CDATA[` + ipaURL + `]]></string>
					</dict>
					<dict>
						<key>kind</key>
						<string>full-size-image</string>
						<key>needs-shine</key>
						<true/>
						<key>url</key>
						<string></string>
					</dict>
					<dict>
						<key>kind</key>
						<string>display-image</string>
						<key>needs-shine</key>
						<true/>
						<key>url</key>
						<string><![CDATA[` + iconURL + `]]></string>
					</dict>
				</array>
				<key>metadata</key>
				<dict>
					<key>bundle-identifier</key>
					<string>` + gameName + `</string>
					<key>bundle-version</key>
					<string><![CDATA[` + version + `]]></string>
					<key>kind</key>
					<string>software</string>
					<key>title</key>
					<string><![CDATA[` + packageName + `]]></string>
				</dict>
			</dict>
		</array>
	</dict>
	</plist>`

	var d1 = []byte(temp)
	if err := ioutil.WriteFile(plist, d1, 0666); err != nil {
		log.Errorf("CreatePlist: err: %v\n", err.Error())
		return false
	}
	return true
}

func fileMD5(file string) (fileMd5 string, err error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	value := md5.Sum(f)
	return fmt.Sprintf("%x", value), err
}
func chineseToRune(str string) string {
	var temp = []rune(str)
	var obj string
	for i := 0; i < len(temp)-2; i++ {
		obj = obj + chineseToStr(temp[i])
	}
	return obj
}

func chineseToStr(chinese rune) string {
	var py utils.Pinyin

	// 初始化，载入汉字拼音映射文件
	py.Init("pinyin_table.txt")

	return py.GetPinyin(chinese, false)
}

func CloudUploadFile1(fr io.Reader, fileName string) error {
	tomlConfig := config.Cfg
	bucketName := tomlConfig.BucketName
	appID := tomlConfig.AppID
	secretID := tomlConfig.SecretID
	secretKey := tomlConfig.SecretKey
	u, _ := url.Parse("https://" + bucketName + "-" + appID + ".cos.ap-shanghai.myqcloud.com")
	client := libcos.NewClient1(u, appID, secretID, secretKey, 10)

	// 上传
	e := client.Put(context.Background(), fileName, fr)
	if e != nil {
		return errorx.Wrap(e)
	}
	return nil
}

func CloudUploadFile2(fileByte []byte, fileName string) error {
	tomlConfig := config.Cfg
	bucketName := tomlConfig.BucketName
	appID := tomlConfig.AppID
	secretID := tomlConfig.SecretID
	secretKey := tomlConfig.SecretKey
	u, _ := url.Parse("https://" + bucketName + "-" + appID + ".cos.ap-shanghai.myqcloud.com")
	client := libcos.NewClient1(u, appID, secretID, secretKey, 10)

	fReader := bytes.NewReader(fileByte)

	// 上传
	err := client.Put(context.Background(), fileName, fReader)
	if err != nil {
		log.Errorf("CloudUploadFile2: Put: err:%+v\n", err)
		return err
	}
	return nil
}

func CloudUploadFile3(filePath string, fileName string) error {
	tomlConfig := config.Cfg
	bucketName := tomlConfig.BucketName
	appID := tomlConfig.AppID
	secretID := tomlConfig.SecretID
	secretKey := tomlConfig.SecretKey
	u, e := url.Parse("http://" + bucketName + "-" + appID + ".cos.ap-shanghai.myqcloud.com")
	if e != nil {
		return errorx.Wrap(e)
	}
	client := libcos.NewClient1(u, appID, secretID, secretKey, 10)

	file, e := os.Open(filePath)
	if e != nil {
		return errorx.Wrap(e)
	}
	defer file.Close()
	stats, e := file.Stat()
	if e != nil {
		return errorx.Wrap(e)
	}
	fileByte := make([]byte, stats.Size())
	file.Read(fileByte)

	fReader := bytes.NewReader(fileByte)

	// 上传
	e = client.Put(context.Background(), fileName, fReader)
	if e != nil {
		return errorx.Wrap(e)
	}
	return nil
}
func UploadFileToHuaweiObsByMultipartFileHeader(fr io.Reader, fileName string) error {
	conf := config.Cfg
	//bucketName := conf.BucketName
	bucketName := conf.HuaweiBucketName
	ak := conf.Ak
	sk := conf.Sk
	endpoint := conf.Endpoint

	buf, e := ioutil.ReadAll(fr)
	if e != nil {
		return errorx.Wrap(e)
	}

	e = libhwobs.PutObject(bucketName, ak, sk, endpoint, fileName, buf)
	if e != nil {
		return errorx.Wrap(e)
	}
	return nil
}

func UploadFileToHuaweiObsByFileName(localPath, fileName string) error {
	conf := config.Cfg
	//bucketName := conf.BucketName
	bucketName := conf.HuaweiBucketName
	ak := conf.Ak
	sk := conf.Sk
	endpoint := conf.Endpoint
	e := libhwobs.PutFile(bucketName, ak, sk, endpoint, fileName, localPath)
	if e != nil {
		return errorx.Wrap(e)
	}
	return nil
}
