package cache

import (
	"fmt"
	"zonst/qipai/api/configapisrv/utils"

	"github.com/garyburd/redigo/redis"
)

const (
	UserCertForceStatusKey = `user:cert:force:status:%v` // string user:cert:force:status:$game_id
)

// UpdateUserCertForceStatus 跟新强制实名认证开启状态
func UpdateUserCertForceStatus(pool *redis.Pool, gameID int, status utils.UserCertForceStatus) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do(`SET`, fmt.Sprintf(UserCertForceStatusKey, gameID), status)
	return err
}

// GetUserCertForceStatus 获取强制实名认证开启状态
func GetUserCertForceStatus(pool *redis.Pool, gameID int) (utils.UserCertForceStatus, error) {
	conn := pool.Get()
	defer conn.Close()
	status, err := redis.Int(conn.Do(`GET`, fmt.Sprintf(UserCertForceStatusKey, gameID)))
	if err != nil && err != redis.ErrNil {
		return 0, err
	}
	return utils.UserCertForceStatus(status), nil
}

// DeleteUserCertForceStatus 删除强制实名认证开启状态
func DeleteUserCertForceStatus(pool *redis.Pool, gameID int) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do(`DEL`, fmt.Sprintf(UserCertForceStatusKey, gameID))
	if err != nil && err != redis.ErrNil {
		return err
	}
	return nil
}
