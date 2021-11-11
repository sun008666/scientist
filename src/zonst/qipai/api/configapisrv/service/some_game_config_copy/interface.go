package some_game_config_copy

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
)

// Interface 某些平台配置 接口
type Interface interface {
	Copy(param Param) error
}

// Param 某些平台配置服务参数定义
type Param struct {
	API         string
	JwtToken    string
	SrcGameID   int32
	DestGameIDs []int32
	QipaiDBOrm  *gorm.DB
	TaskDB      *sqlx.DB
	TaskDBOrm   *gorm.DB
}

type CopyType int32

const (
	SIGN_AWARD                CopyType = 1
	SHOP_CONFIG               CopyType = 2
	TASK                      CopyType = 3
	ROOMCARD_PAY_LIST         CopyType = 4
	DIAMONDS_PAY_LIST         CopyType = 5
	SUBSIDY_CONFIG            CopyType = 6
	WECHAT_SHARE_BONUS_CONFIG CopyType = 7
)

// NotFoundCopyType 未找到要复制的类型
var NotFoundCopyType = errors.New(`not found copy type`)

// InterfaceFactory 结构工厂方法
func InterfaceFactory(t CopyType) (Interface, error) {
	switch t {
	case SHOP_CONFIG:
		return ShopConfig{}, nil
	case DIAMONDS_PAY_LIST:
		return DiamondsPayList{}, nil
	case ROOMCARD_PAY_LIST:
		return RoomCardPayList{}, nil
	case SIGN_AWARD:
		return SignAward{}, nil
	case TASK:
		return Task{}, nil
	case SUBSIDY_CONFIG:
		return SubsidyConfig{}, nil
	case WECHAT_SHARE_BONUS_CONFIG:
		return WechatShareBonusConfig{}, nil
	default:
		return nil, NotFoundCopyType
	}
}
