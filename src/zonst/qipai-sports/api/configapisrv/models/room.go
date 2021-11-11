package models

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Room 游戏房间服务配置
type Room struct {
	GameID         int32  `json:"game_id" db:"game_id"`                   // 游戏ID
	GameAreaID     int32  `json:"game_area_id" db:"game_area_id"`         // 子游戏ID
	RoomID         int32  `json:"room_id" db:"room_id"`                   // 房间ID
	RoomType       int32  `json:"room_type" db:"room_type"`               // 房间类型
	RoomServerIP   string `json:"room_server_ip" db:"room_server_ip"`     // 房间服务所在ip
	RoomServerPort int32  `json:"room_server_port" db:"room_server_port"` // 房间服务端口
	MasterNodeID   int32  `json:"master_node_id" db:"master_node_id"`     // 主节点ID
	DeskNums       int32  `json:"desk_nums" db:"desk_nums"`               // 桌子数量
	SeatNums       int32  `json:"seat_nums" db:"seat_nums"`               // 座位数量
	MinUserNums    int32  `json:"min_user_nums" db:"min_user_nums"`       // 最小游戏人数
	ClubNodeID     int32  `json:"club_node_id" db:"club_node_id"`         // 俱乐部服务ID
	MatchNodeID    int32  `json:"match_node_id" db:"match_node_id"`       // 比赛服务ID
	CoinNodeID     int32  `json:"coin_node_id" db:"coin_node_id"`         // 金币场服务ID
}

// Validate 验证
func (r *Room) Validate() error {
	if r.GameID == 0 {
		return fmt.Errorf("Room:GameID(%v)==0", r.GameID)
	}
	if r.GameAreaID <= -1 {
		return fmt.Errorf("Room:GameAreaID(%v)<=-1", r.GameAreaID)
	}
	if r.DeskNums <= 0 {
		return fmt.Errorf("Room:DeskNums(%v)<=0", r.DeskNums)
	}
	if r.SeatNums < 2 {
		return fmt.Errorf("Room:SeatNums(%v)<2", r.SeatNums)
	}
	if r.MinUserNums < 2 {
		return fmt.Errorf("Room:MinUserNums(%v)<2", r.MinUserNums)
	}
	if r.MasterNodeID <= -1 {
		return fmt.Errorf("Room:MasterNodeID(%v)==0", r.MasterNodeID)
	}

	return nil
}

// GetRoom 获取房间服务信息
func GetRoom(ctx *gin.Context, roomID int32) (*Room, error) {

	db := ctx.MustGet("configdb").(*sqlx.DB)

	tpl := "SELECT game_id, game_area_id, id AS room_id, room_type, room_server_ip, room_server_port, " +
		"master_node_id, desk_nums, seat_nums, min_user_nums, club_node_id, match_node_id, coin_node_id " +
		"FROM game_room_sh " +
		"WHERE id=$1;"

	obj := &Room{}
	if err := db.Unsafe().Get(obj, tpl, roomID); err != nil {
		return nil, err
	}

	if err := obj.Validate(); err != nil {
		return nil, err
	}
	return obj, nil
}
