package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

var Cfg *Config

// Config 对应配置文件结构
type Config struct {
	Listen                  string                 `toml:"listen"`
	HotUpdateAddr           string                 `toml:"hotUpdateAddr"`
	PackageDownloadAddr     string                 `toml:"packageDownloadAddr"`
	DBServers               map[string]DBServer    `toml:"dbservers"`
	RedisServers            map[string]RedisServer `toml:"redisservers"`
	FileAddrConfig          FileAddrConfig         `toml:"fileAddrConfig"`
	Platform2FileAddrConfig Platform2CloudFileAddr `toml:"platform2FileAddrConfig"`
	PlatFormURL             string                 `toml:"platFormURL"`
	//CDNURL               string                 `toml:"cdnURL"`
	//CDNUserName          string                 `toml:"cdnUserName"`
	//CDNAPIKey            string                 `toml:"cdnAPIKey"`
	HotupdatedownloadURL string `toml:"hotupdatedownloadURL"`
	//CDNResp              string                 `toml:"cdnResp"`
	CloudFileAddr          CloudFileAddr          `toml:"cloudFileAddr"`
	Platform2CloudFileAddr Platform2CloudFileAddr `toml:"platform2CloudFileAddr"`
	BucketName             string                 `toml:"bucketName"`
	AppID                  string                 `toml:"appID"`
	SecretID               string                 `toml:"secretID"`
	SecretKey              string                 `toml:"secretKey"`
	CoroutineCount         int                    `toml:"coroutineCount"`
	GameListAddr           string                 `toml:"gameListAddr"`
	ConfigAddr             string                 `toml:"configAddr"`
	BackUpURL              string                 `toml:"backupUrl"`
	WechatKeysURL          string                 `toml:"wechatkeys_url"`
	ResourceKey            string                 `toml:"resource_key"`
	AuthorityURL           string                 `toml:"authority_url"`
	GlobalWhiteList        map[string]bool        `toml:"global_whitelist_game"` // 支持全局白名单时，白名单更新会同步到所有平台, key为game_id
	PlatformFlag           string                 `toml:"platform_flag"`
	EnvironmentFlag        string                 `toml:"environment_flag"`
	HuaweiBucketName       string                 `toml:"huawei_bucket_name"`
	Ak                     string                 `toml:"ak"`
	Sk                     string                 `toml:"sk"`
	Endpoint               string                 `toml:"endpoint"`
	SyncGameLibrary        bool                   `toml:"sync_game_library"`
	SyncPlatform           bool                   `toml:"sync_platform"`
	SyncPlatformGameIDs    map[string]bool        `toml:"sync_platform_game_ids"`
	SyncDomainName         string                 `toml:"sync_domain_name"`
	CDNAPIURL              string                 `toml:"cdnapi_url"`
	PushClientAddr         string                 `toml:"pushClientAddr"` //发送消息发客户端中转服务pushclientapisrv的IP
	CDNAPISrv              string                 `toml:"cdnapi_srv"`
	CdnSecretID            string                 `toml:"cdnSecretID"`
	CdnSecretKEY           string                 `toml:"cdnSecretKEY"`
	GameVersionPath        string                 `toml:"gameVersionPath"`
	AndroidPackageURL      string                 `toml:"androidPackageUrl"`
	CDNAPI                 string                 `toml:"cdn_api"`
	IosPackageURL          string                 `toml:"iosPackageUrl"`
	PackageIconURL         string                 `toml:"packageIconUrl"`
}

// FileAddrConfig 定义长传地址的配置
type FileAddrConfig struct {
	Image           string `toml:"image"`
	Patch           string `toml:"patch"`
	PatchNew        string `toml:"patchNew"`
	DifferencePatch string `toml:"differencePatch"` //差量补丁存放路径
	CompletePatch   string `toml:"completePatch"`   //全量补丁存放路径
	PatchMutil      string `toml:"patchMutil"`
	Package         string `toml:"package"`
	JSONPath        string `toml:"jsonPath"`
}

type Platform2FileAddrConfig struct {
	Patch    string `toml:"patch"`
	PatchNew string `toml:"patchNew"`
	JSONPath string `toml:"jsonPath"`
}

// CloudFileAddr 云端文件存储地址
type CloudFileAddr struct {
	Image      string `toml:"image"`
	Patch      string `toml:"patch"`
	PatchNew   string `toml:"patchNew"`
	PatchMutil string `toml:"patchMutil"`
	Package    string `toml:"package"`
	JSONPath   string `toml:"jsonPath"`
}

type Platform2CloudFileAddr struct {
	Patch      string `toml:"patch"`
	PatchNew   string `toml:"patchNew"`
	PatchMutil string `toml:"patchMutil"`
	JSONPath   string `toml:"jsonPath"`
}

//var FileConfig = struct {
//	FileAddr            FileAddrConfig `mapstructure:"file_addr"`
//	HotUpdateAddr       string         `mapstructure:"hot_update_addr"`
//	PackageDownloadAddr string         `mapstructure:"package_download_addr"`
//}{}

// UnmarshalConfig 解析toml配置
func UnmarshalConfig(tomlfile string) (*Config, error) {
	c := &Config{}
	if _, err := toml.DecodeFile(tomlfile, c); err != nil {
		return c, err
	}
	return c, nil
}

// DBServerConf 获取数据库配置
func (c Config) DBServerConf(key string) (DBServer, bool) {
	s, ok := c.DBServers[key]
	return s, ok
}

// RedisServerConf 获取数据库配置
func (c Config) RedisServerConf(key string) (RedisServer, bool) {
	s, ok := c.RedisServers[key]
	return s, ok
}

// GetListenAddr 监听地址
func (c Config) GetListenAddr() string {
	return c.Listen
}

// GetHotUpdateAddr 主动更新版本接口地址
func (c Config) GetHotUpdateAddr() string {
	return c.HotUpdateAddr
}

// GetPackageDownloadAddr 整包下载地址
func (c Config) GetPackageDownloadAddr() string {
	return c.PackageDownloadAddr
}

//GetFileAddrConfig 保存文件路径
func (c Config) GetFileAddrConfig() FileAddrConfig {
	return c.FileAddrConfig
}

//GetSyncGameLibrary 同步游戏库标志
func (c Config) GetSyncGameLibrary() bool {
	return c.SyncGameLibrary
}

//GetSyncDomainName 同步域名
func (c Config) GetSyncDomainName() string {
	return c.SyncDomainName
}

// RedisServer 表示 redis 服务器配置
type RedisServer struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

// DBServer 表示DB服务器配置
type DBServer struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	DBName   string `toml:"dbname"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

// ConnectString 表示连接数据库的字符串
func (m DBServer) ConnectString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		m.Host, m.Port, m.User, m.Password, m.DBName)
}
