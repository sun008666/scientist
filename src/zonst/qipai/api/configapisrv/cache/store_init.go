package cache

import (
	"fmt"
	"zonst/qipai/api/configapisrv/config"
	"zonst/qipai/api/configapisrv/models"

	"github.com/go-xweb/log"

	"github.com/jinzhu/gorm"
)

// Store 内存存储
type Store struct {
	DistrictTable *TblGameDistrictConfig // 当前地区配置缓存
	cfg           *config.Config         // 配置
}

// NewStore 初始化存储器
func NewStore(cfg *config.Config) *Store {
	if cfg == nil {
		panic("cfg is nil")
	}

	obj := &Store{
		DistrictTable: NewKSCache(),
		cfg:           cfg}
	return obj
}

// Init 初始化
func (s *Store) Init() {
	db := newDBPool(s.cfg)
	// 创建DB对象
	initDistrictCache(db, s.DistrictTable)

}

// NewDBPool 初始化
func NewDBPool(cfg *config.Config) *gorm.DB {
	return newDBPool(cfg)
}

func newDBPool(cfg *config.Config) *gorm.DB {
	if cfg == nil {
		panic("newDBPool: cfg is nil")
	}

	dbName := "qipaidb"
	// 读取配置
	dbConfig, ok := cfg.DBServerConf(dbName)
	if !ok {
		panic(fmt.Sprintf("Postgres: %v no set.", dbName))
	}

	db, err := gorm.Open("postgres", dbConfig.ConnectString())
	if err != nil {
		panic(fmt.Sprintf("gorm.Open: err:%v", err))
	}

	return db

}

func initDistrictCache(db *gorm.DB, cacheObj *TblGameDistrictConfig) {
	objs, err := models.GetDistrictConfigList(db, 0, "", "")
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf("initDistrictCache: %v\n", err)
		return
	}

	tmpTable := make(map[int32][]*models.GameDistrictConfig)
	for _, v := range objs {
		if value, ok := tmpTable[v.GameID]; ok {
			value = append(value, v)
			tmpTable[v.GameID] = value
		} else {
			temp := make([]*models.GameDistrictConfig, 0)
			temp = append(temp, v)
			tmpTable[v.GameID] = temp

		}
	}
	cacheObj.Init(tmpTable)
	log.Debugf("initDistrictCache: 加载%v条数据\n", len(tmpTable))
}
