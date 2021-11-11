package test

var (
	// 正式配置
	ClubdbConfigRelease = DBConfig{
		Host:     "10.133.196.213",
		Port:     5432,
		UserName: "qipai",
		Password: "qipai#xq5",
		DBName:   "clubdb",
	}

	ClubCacheConfigRelease = RedisConfig{
		Addr:     "10.66.206.3:6379",
		Password: "crs-bprrl43i:qipai0918",
		DB:       0,
	}

	// 测试配置
	ClubdbConfigTest = DBConfig{
		Host:     "127.0.0.1",
		Port:     5432,
		UserName: "postgres",
		Password: "123456",
		DBName:   "qipaidb",
	}

	ClubCacheConfigTest = RedisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}

	LogdbdbConfigTest = DBConfig{
		Host:     "127.0.0.1",
		Port:     5432,
		UserName: "postgres",
		Password: "123456",
		DBName:   "qipaidb",
	}
)
