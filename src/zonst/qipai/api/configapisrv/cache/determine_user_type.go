package cache

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

const (
	//
	// determineUserTypeKey
	/**
	 * type: string
	 * key: three:lobby:update:$user_id
	 * value:  (具体含义见 models.determine_user_type BindGameID 字段)
	 */
	determineUserTypeKey = "determine:user:type:%d"
)

type DeterMineUserTypeFind struct {
	BindGameID int `json:"bind_game_id"`
}

func (d *DeterMineUserTypeFind) genKey(userID int) string {
	return fmt.Sprintf(determineUserTypeKey, userID)
}

// FindUserTypeByUserID 查询用户所属地区
func (d *DeterMineUserTypeFind) FindUserTypeByUserID(conn redis.Conn, userID int) (DeterMineUserTypeFind, bool, error) {
	key := d.genKey(userID)
	userType, e := redis.Bytes(conn.Do(`GET`, key))
	var determineUserType DeterMineUserTypeFind
	if e != nil {
		return determineUserType, false, e
	}
	e = json.Unmarshal(userType, &determineUserType)
	if e != nil {
		return determineUserType, false, e
	}
	return determineUserType, e == nil, nil
}

// CreatList 往缓存里面加入用户类型
func (d DeterMineUserTypeFind) CreatList(conn redis.Conn, userID int) error {
	key := d.genKey(userID)
	b, _ := json.Marshal(&d)
	_, e := conn.Do(`SET`, key, b)
	if e != nil {
		return e
	}
	_, _ = conn.Do(`EXPIRE`, key, 3600)
	return nil
}
