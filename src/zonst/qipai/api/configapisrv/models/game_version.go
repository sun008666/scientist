package models

import "github.com/jinzhu/gorm"

const (
	CloseStatus = 1
	Openstatus  = 2
)

// GamePackageTable 游戏版本表
type GamePackageTable struct {
	ID               int    `gorm:"column:id" json:"id"`
	GameID           int    `gorm:"column:game_id" json:"game_id"`
	GameName         string `gorm:"column:game_name" json:"game_name"`
	PackageVersion   string `gorm:"column:package_version" json:"package_version"`
	PackageVersionNO int    `gorm:"column:package_version_no" json:"package_version_no"`
	PackageOS        string `gorm:"column:package_os" json:"package_os"`
	PackageURL       string `gorm:"column:package_url" json:"package_url"`
	PackageName      string `gorm:"column:package_name" json:"package_name"`
	PackageIcon      string `gorm:"column:package_icon" json:"package_icon"`
	PackageFullImage string `gorm:"column:package_full_image" json:"package_full_image"`
	PackageMD5       string `gorm:"column:package_md5" json:"package_md5"`
	PackageSize      int64  `gorm:"column:package_size" json:"package_size"`
	LogTime          string `gorm:"column:log_time" json:"log_time"`
	Remark           string `gorm:"column:remark" json:"remark"`
	Status           int    `gorm:"column:status" json:"status"` //状态 1:关闭 2:开启
}

func (g *GamePackageTable) TableName() string {
	return "game_package"
}

type OnGameVersionListResponce struct {
	ID          int    `json:"id"`
	GameID      int    `json:"game_id"`
	GameName    string `json:"game_name"`
	Version     string `json:"version"`
	PackageOS   string `json:"package_os"`
	PackageMD5  string `json:"package_md5"`
	PackageSize string `json:"package_size"`
	LogTime     string `json:"log_time"`
	Remark      string `json:"remark"`
	Status      int    `json:"status"`
}

type PackageVersionNO struct {
	NO int `gorm:"column:no" json:"no"`
}

func GetGameVersionPackageByID(qipaidb *gorm.DB, id int32) (*GamePackageTable, error) {
	data := &GamePackageTable{}
	if err := qipaidb.Debug().Where("id=?", id).First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func UpdateGameVersionPackageStatus(qipaidb *gorm.DB, openPackageID, status int32) error {
	data := GamePackageTable{}
	//tx:=qipaidb.Begin()

	if err := qipaidb.Debug().Table(data.TableName()).Where("id=?", openPackageID).Update("status", status).Error; err != nil {
		//tx.Rollback()
		return err
	}

	//if err:=tx.Debug().Table(data.TableName()).Where("id=?",closePackageID).Update("status",CloseStatus).Error;err!=nil{
	//	tx.Rollback()
	//	return err
	//}
	//tx.Commit()
	return nil
}
