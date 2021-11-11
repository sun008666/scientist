package cache

import (
	"errors"
	"sync"
	"zonst/qipai/api/configapisrv/models"

	"github.com/go-xweb/log"
)

type TblGameDistrictConfig struct {
	sync.RWMutex
	Table map[int32][]*models.GameDistrictConfig
}

func NewKSCache() *TblGameDistrictConfig {
	// config := make([]*GameDistrictConfig, 0)
	obj := &TblGameDistrictConfig{Table: make(map[int32][]*models.GameDistrictConfig)}

	return obj
}
func (t *TblGameDistrictConfig) Get(gameID int32) ([]*models.GameDistrictConfig, error) {
	t.RLock()
	defer t.RUnlock()
	obj, ok := t.Table[gameID]
	if !ok {
		return nil, errors.New("key not found")
	}

	return obj, nil
}

func (t *TblGameDistrictConfig) Set(gameID int32, obj []*models.GameDistrictConfig) {
	t.Lock()
	defer t.Unlock()
	t.Table[gameID] = obj
}

func (t *TblGameDistrictConfig) Init(m map[int32][]*models.GameDistrictConfig) {
	t.Lock()
	defer t.Unlock()
	t.Table = m
}

// GetDistrictConfig 获取当前地区的微信公众号和客服微信、跑马灯内容
func (t *TblGameDistrictConfig) GetDistrictConfig(province, city, district string, gameID int32) (*models.GameDistrictConfig, error) {
	t.RLock()
	defer t.RUnlock()
	obj, ok := t.Table[gameID]
	if !ok {
		log.Debugf("key:%#v not found", gameID)
		return nil, errors.New("key not found")
	}

	for _, v := range obj {
		if v.Province == province && v.City == city {
			if len(v.DistrictArray) == 0 {
				return v, nil
			}
			for _, va := range v.DistrictArray {
				if va == district {
					return v, nil
				}

			}

		}
	}

	return nil, errors.New("key not found")
}
