package test

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     int64
	UserName string
	Password string
	DBName   string
}

func (c DBConfig) ToString() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", c.Host, c.Port,
		c.UserName, c.Password, c.DBName,
	)
}

func (c DBConfig) Connect() (db *sqlx.DB) {
	// 链接数据库
	db, err := sqlx.Open("postgres", c.ToString())
	if err != nil {
		log.Fatalln(err)
	}

	db.SetMaxIdleConns(5)
	return db
}
