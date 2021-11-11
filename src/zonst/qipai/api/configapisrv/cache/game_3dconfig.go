package cache

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/lib/pq"
	"strconv"
)

const (
	//
	// threeDLobbyGameUpdateKey
	/**
	 * type: hash
	 * key: three:lobby:update:$game_id
	 * field: $game_id
	 * value:  (具体含义见 models.game_3dconfig GameAreaID 字段)
	 */
	threeDLobbyGameUpdateKey = "three:lobby:update:%d"
)

// ThreeDLobbyGameUpdate 子游戏缓存的结构体
type ThreeDLobbyGameUpdate struct {
	GameAreaID pq.Int64Array `json:"game_area_id"`
}

// GenKey1 获取key
func (*ThreeDLobbyGameUpdate) genKey1(gameID int) string {
	return fmt.Sprintf(threeDLobbyGameUpdateKey, gameID)
}

//
func (*ThreeDLobbyGameUpdate) genField(gameID int, gameAreaID pq.Int64Array) string {
	return fmt.Sprintf(strconv.Itoa(gameID), gameAreaID)
}

// CreateList 往缓存里面添加游戏一二的数据
func (a ThreeDLobbyGameUpdate) CreateList(conn redis.Conn, gameID int) error {
	key := a.genKey1(gameID)
	bytes, _ := json.Marshal(a)
	_, err := conn.Do(`HSET`, key, gameID, bytes)
	if err != nil {
		return err
	}
	_, _ = conn.Do(`EXPIRE`, key, 86400)
	return nil
}

// FindByGameID 从缓存当中获取GameAreaID1和GameAreaID2
func (a *ThreeDLobbyGameUpdate) FindByGameID(conn redis.Conn, gameID int) (ThreeDLobbyGameUpdate, bool, error) {
	key := a.genKey1(gameID)
	gameArea, err := redis.Bytes(conn.Do(`HGET`, key, gameID))
	var userInfo ThreeDLobbyGameUpdate
	if err != nil {
		return userInfo, false, err
	}
	err = json.Unmarshal(gameArea, &userInfo)
	return userInfo, err == nil, err
}

// DeleteCache 删除缓存
func (a *ThreeDLobbyGameUpdate) DeleteCache(conn redis.Conn, gameID int) error {
	key := a.genKey1(gameID)
	_, err := conn.Do(`DEL`, key)
	return err
}
