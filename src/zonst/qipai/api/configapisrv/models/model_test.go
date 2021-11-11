package models

import (
	"testing"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("postgres", "host=127.0.0.1 user=huangxin password=123 dbname=test sslmode=disable")
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
}

func TestGetAllPage(t *testing.T) {
	reply, err := GetAllPage(db)
	if err != nil {
		t.Errorf("err:%v", err)
	} else {
		t.Logf("request:%+v \n", reply)
	}

	//fmt.Println(fmt.Sprintf("%v平台刷新cdn失败，请稍后手动重试",[]string{"1","2"}))
}
