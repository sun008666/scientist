package cache

import (
	"fmt"
	"zonst/qipai/api/configapisrv/utils"

	"github.com/garyburd/redigo/redis"
)

const (
	UserCertForceScopeKey = `user:cert:force:scope:%v` // string user:cert:force:scope:$game_id
)

// UpdateUserCertForceScope 跟新强制实名认证范围
func UpdateUserCertForceScope(pool *redis.Pool, gameID int, status utils.UserCertForceScope) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do(`SET`, fmt.Sprintf(UserCertForceScopeKey, gameID), status)
	return err
}

// GetUserCertForceScope 获取强制实名认证范围
func GetUserCertForceScope(pool *redis.Pool, gameID int) (utils.UserCertForceScope, error) {
	conn := pool.Get()
	defer conn.Close()
	status, err := redis.Int(conn.Do(`GET`, fmt.Sprintf(UserCertForceScopeKey, gameID)))
	if err != nil && err != redis.ErrNil {
		return 0, err
	}
	return utils.UserCertForceScope(status), nil
}

// DeleteUserCertForceScope 删除强制实名认证范围
func DeleteUserCertForceScope(pool *redis.Pool, gameID int) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do(`DEL`, fmt.Sprintf(UserCertForceScopeKey, gameID))
	if err != nil && err != redis.ErrNil {
		return err
	}
	return nil
}
