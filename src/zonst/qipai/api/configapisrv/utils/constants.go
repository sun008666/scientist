package utils

const (
	On  = 1
	Off = 0

	ClubAppID        string = "3c49995cda17777b5fb5"             // appid
	ClubAppSecretKey string = "c6d25a864c00a154d98f7a13d57a21fd" // appkey
	LayoutTime              = "2006-01-02 15:04:05"
)

// UserCertForceStatus 玩家强制实名认证
type UserCertForceStatus int

// UserCertForceScope 玩家强制实名认证范围
type UserCertForceScope int

const (
	// 关闭强制实名认知
	CloseUserCertForce UserCertForceStatus = iota
	// 开启强制实名认证
	OpenUserCertForce
)

const (
	// 江西省外
	OutsideJiangXi UserCertForceScope = iota
	// 江西省内
	InsideJiangXi
	// 全国
	WholeCountry
)
const (
	ClubDBOrm    = "clubdborm"
	ClubDB       = "clubdb"
	TargetHaoPai = "hao_pai"
	QipaiDBOrm   = "qipaidborm"
	TaskDB       = "taskdb"
	TaskDBOrm    = "taskdborm"
)
