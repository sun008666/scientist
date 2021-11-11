package cache

import (
	"zonst/qipai/api/configapisrv/utils"

	"github.com/garyburd/redigo/redis"
)

// 返回值key表示game_id，value表示禁止状态，1表示已禁止
func BatchGetGameForbidStatus(cache *redis.Pool, gameIDs []int32) (map[int32]int32, error) {
	conn := cache.Get()
	defer conn.Close()

	key := "forbidden:games"
	s, err := redis.Ints(conn.Do("SMEMBERS", key))
	if err != nil {
		return nil, err
	}

	m := make(map[int32]int32, len(gameIDs)+1)
	for _, v := range gameIDs {
		m[v] = utils.Off
	}

	for _, v := range s {
		if _, ok := m[int32(v)]; ok {
			m[int32(v)] = utils.On
		}
	}
	return m, nil
}

func UpdateGameForbidStatus(cache *redis.Pool, gameID, forbid int32) error {
	conn := cache.Get()
	defer conn.Close()

	key := "forbidden:games"
	var err error

	if forbid == utils.On {
		_, err = conn.Do("SADD", key, gameID)
	} else {
		_, err = conn.Do("SREM", key, gameID)
	}

	return err
}
