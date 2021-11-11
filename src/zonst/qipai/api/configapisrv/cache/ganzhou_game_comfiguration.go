package cache

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

const (
	//
	// ganzhouGameConfiguration
	/**
	 * type: string
	 * key: ganzhou:game:configuration
	 * value:  (具体含义见 models.game_3dconfig GameAreaID 字段)
	 */
	ganzhouGameConfiguration = "ganzhou:game:configuration"
)

// GanzhouGameConfiguration 子游戏缓存的结构体
type GanzhouGameConfiguration struct {
	ID         int `json:"id"`
	GameID     int `json:"game_id"`
	GameAreaID int `json:"game_area_id"`
	SubGameID  int `json:"sub_game_id"`
}

// FindGanConfiguration 从缓存当中获取赣州游戏的配置
func (g *GanzhouGameConfiguration) FindGanConfiguration(conn redis.Conn) ([]GanzhouGameConfiguration, bool, error) {
	gameAreaConfiguration, e := redis.Bytes(conn.Do(`GET`, ganzhouGameConfiguration))
	var ganzhouConfiguration []GanzhouGameConfiguration
	if e != nil {
		return ganzhouConfiguration, false, e
	}
	e = json.Unmarshal(gameAreaConfiguration, &ganzhouConfiguration)
	return ganzhouConfiguration, e == nil, e
}

// CreateList 往缓存里面添加地区配置数据
func (g GanzhouGameConfiguration) CreateList(conn redis.Conn, c []GanzhouGameConfiguration) error {
	bytes, _ := json.Marshal(g)
	_, err := conn.Do(`SET`, ganzhouGameConfiguration, bytes)
	if err != nil {
		return err
	}
	_, _ = conn.Do(`EXPIRE`, ganzhouGameConfiguration, 86400)
	return nil
}

// DeleteCache 删除缓存
func DeleteCache(conn redis.Conn) error {
	_, err := conn.Do(`DEL`, ganzhouGameConfiguration)
	return err
}
