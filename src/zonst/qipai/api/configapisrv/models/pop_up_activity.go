package models

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-xweb/log"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// PopUpActivityTable 弹窗活动
type PopUpActivityTable struct {
	ID             int            `gorm:"id" json:"id"`
	GameID         int            `gorm:"game_id" json:"game_id"`
	Title          string         `gorm:"title" json:"title"`                   //标题
	ImgURL         string         `gorm:"img_url" json:"img_url"`               //图片链接
	StartTime      int            `gorm:"start_time" json:"start_time"`         //开始时间
	EndTime        int            `gorm:"end_time" json:"end_time"`             //结束时间
	LinkTo         int            `gorm:"link_to" json:"link_to"`               //跳转至， 任务类型
	LinkURL        string         `gorm:"link_url" json:"link_url"`             //外部链接URL
	IsShow         int            `gorm:"is_show" json:"is_show"`               //激活
	SendToClient   int            `gorm:"send_to_client" json:"send_to_client"` //显示
	ShowWho        int            `gorm:"column:show_who" json:"show_who"`      //
	Position       int            `gorm:"column:position" json:"position"`      //位置
	RedDot         int            `gorm:"column:red_dot" json:"red_dot"`
	PlayerType     pq.Int64Array  `gorm:"column:player_type" json:"player_type"`
	PlayerIdentity pq.Int64Array  `gorm:"column:player_identity" json:"player_identity"`
	ShowArea       pq.StringArray `gorm:"column:show_area" json:"show_area"`
	ShowGameAreaID pq.Int64Array  `gorm:"column:show_game_area_id" json:"show_game_area_id"`
	IsWhiteList    bool           `json:"is_white_list" gorm:"column:is_white_list"` // 是否是白名单    true false
	ShowSystem     string         `json:"show_system" gorm:"column:show_system"`     // 显示的操作系统 all android ios
	GameList       pq.Int64Array  `json:"game_list" gorm:"column:game_list"`
}

type PositionMax struct {
	NO int `gorm:"column:no" json:"no"`
}

func AddPopUpActivity(ctx *gin.Context, obj *PopUpActivityTable) error {
	db := ctx.MustGet("qipaidb").(*gorm.DB)
	if obj.ShowSystem != "all" && obj.ShowSystem != "android" && obj.ShowSystem != "ios" {
		log.Errorf("AddPopUpActivity: err: %v\n", "操作系统只能是all或android或ios")
		return errors.New("show_system只能取all或android或ios")
	}
	positionMax := PositionMax{}
	if err := db.Debug().Table("pop_up_activity").Where("game_id = ?", obj.GameID).Select("max(position) as no").Find(&positionMax).Error; err != nil {
		log.Errorf("AddPopUpActivity: err: %v", err.Error())
		return err
	}
	obj.Position = positionMax.NO + 1

	if err := db.Debug().Table("pop_up_activity").Create(obj).Error; err != nil {
		log.Errorf("AddPopUpActivity: err: %v\n", err.Error())
		return err
	}
	return nil
}

func UpdatePopActivity(c *gin.Context, id int, obj *PopUpActivityTable) error {
	db := c.MustGet("qipaidb").(*gorm.DB)
	if err := db.Debug().Table("pop_up_activity").Where("id = ?", id).Updates(obj).Update("is_show", obj.IsShow).Update("send_to_client", obj.IsShow).Update("is_white_list", obj.IsWhiteList).Error; err != nil {
		log.Errorf("UpdatePopActivity: err: %v\n", err.Error())
		return err
	}
	return nil
}

func GetPopActivityByID(qipaidb *gorm.DB, id int) (*PopUpActivityTable, error) {
	data := &PopUpActivityTable{}
	if err := qipaidb.Debug().Table("pop_up_activity").Where("id=?", id).First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
