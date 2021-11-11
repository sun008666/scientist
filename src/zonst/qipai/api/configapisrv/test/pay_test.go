package test

import (
	"testing"
	"time"
	"zonst/qipai/api/configapisrv/handlers"
	"zonst/qipai/api/configapisrv/utils"
	unitTest "zonst/qipai/gin-unittest-demo"
)

type CommonResp struct {
	Errno  string
	Errmsg string
}

// RoomCardProductList 充值列表
type RoomCardProductList struct {
	ID                  int    `gorm:"column:id" json:"id"`
	GameID              int    `gorm:"column:game_id" json:"game_id"`
	GameAreaID          string `gorm:"column:game_area_id" json:"game_area_id"`
	CardNum             int    `gorm:"column:card_num" json:"card_num"`
	Status              int    `gorm:"column:status" json:"status"`
	Price               int    `gorm:"column:price" json:"price"`
	AddTime             string `gorm:"column:add_time" json:"add_time"`
	Category            string `gorm:"column:category" json:"category"`
	PresentRoomCard     int    `gorm:"column:present_room_card" json:"present_room_card"`
	EndTime             string `gorm:"column:end_time" json:"end_time"`
	Represent           string `gorm:"column:represent" json:"represent"`
	AppleKey            string `gorm:"column:apple_key" json:"apple_key"`
	IsWhite             bool   `gorm:"column:is_white" json:"is_white"`
	IsBlack             bool   `gorm:"column:is_black" json:"is_black"`
	CardPackID          int32  `gorm:"card_pack_id" json:"card_pack_id"`
	IsSeniorAgent       int32  `gorm:"is_senior_agent" json:"is_senior_agent"`
	AddTimeStr          string `gorm:"-" json:"-"`
	EndTimeStr          string `gorm:"-" json:"-"`
	PurchaseLimitNumber int32  `gorm:"column:purchase_limit_number" json:"purchase_limit_number"`
}
type OnPayListListRequestResp struct {
	Data   []*RoomCardProductList `json:"data"`
	Errno  string
	Errmsg string
}

var (
	gameIDList = []int{65, 66, 65, 67, 68}
	gameID     = 66
	payMap     = make(map[int][]*RoomCardProductList, 0)
)

// TestOnPayListAddRequest 测试限购次数数据添加功能
func TestOnPayListAddRequest(t *testing.T) {
	uri := "/v1/pay/list/add"
	nowString := time.Now().Format(utils.LayoutTime)
	for _, val := range gameIDList {
		param := &handlers.OnPayListAddReq{
			GameID:              val,
			GameAreaID:          "-1",
			CardNum:             12,
			Status:              1,
			Price:               100,
			PresentRoomCard:     20,
			PurchaseLimitNumber: 3,
			BlackList:           []int{123, 123, 111},
			WhiteList:           []int{222, 333},
			AddTime:             nowString,
			EndTime:             time.Now().Format(utils.LayoutTime),
		}
		var resp CommonResp
		err := unitTest.TestHandlerUnMarshalResp("POST", uri, "json", param, &resp)
		if err != nil {
			t.Fatalf("TestOnPayListAddRequest: 请求出错，err:%v\n", err)
		}
		if resp.Errno != "0" {
			t.Fatalf("TestOnPayListAddRequest: 请求出错, 不符合预期，resp:%v\n", resp)
		}
	}
	return
}

func TestOnPayListListRequest(t *testing.T) {
	uri := "/v1/pay/list/list"
	var resp OnPayListListRequestResp
	param := &handlers.OnPayListListReq{
		GameID:     66,
		GameAreaID: -1,
	}
	err := unitTest.TestHandlerUnMarshalResp("POST", uri, "json", param, &resp)
	if err != nil {
		t.Fatalf("TestOnPayListListRequest: 请求出错，err:%v\n", err)
	}
	if resp.Errno != "0" {
		t.Fatalf("TestOnPayListListRequest: 请求出错, 不符合预期，resp:%v\n", resp)
	}
	for _, val := range resp.Data {
		if _, ok := payMap[val.GameID]; !ok {
			payMap[val.GameID] = make([]*RoomCardProductList, 0)
		}
		payMap[val.GameID] = append(payMap[val.GameID], val)
	}
	return
}

// TestOnPayListUpdateRequest 测试修改限购次数修改功能
func TestOnPayListUpdateRequest(t *testing.T) {
	id := payMap[gameID][0].ID
	purchaseLimitNumber := int32(5)
	uri := "/v1/pay/list/update"
	param := handlers.OnPayListUpdateReq{
		GameID:              gameID,
		ID:                  id,
		CardNum:             32,
		PresentRoomCard:     55,
		PurchaseLimitNumber: purchaseLimitNumber,
		WhiteList:           []int{},
		BlackList:           []int{},
		AddTime:             time.Now().Format(utils.LayoutTime),
		EndTime:             time.Now().Format(utils.LayoutTime),
	}
	var resp CommonResp
	err := unitTest.TestHandlerUnMarshalResp("POST", uri, "json", param, &resp)
	if err != nil {
		t.Fatalf("TestOnPayListUpdateRequest: 请求出错，err:%v\n", err)
	}
	if resp.Errno != "0" {
		t.Fatalf("TestOnPayListUpdateRequest: 请求出错, 不符合预期，resp:%v\n", resp)
	}

	uri2 := "/v1/pay/list/list"
	var resp2 OnPayListListRequestResp
	param2 := &handlers.OnPayListListReq{
		GameID:     gameID,
		GameAreaID: -1,
	}
	err = unitTest.TestHandlerUnMarshalResp("POST", uri2, "json", param2, &resp2)
	if err != nil {
		t.Fatalf("TestOnPayListUpdateRequest: 请求出错，err:%v\n", err)
	}
	if resp2.Errno != "0" {
		t.Fatalf("TestOnPayListUpdateRequest: 请求出错, 不符合预期，resp2:%v\n", resp2)
	}
	if len(resp2.Data) <= 0 || resp2.Data[0].ID != id || resp2.Data[0].PurchaseLimitNumber != purchaseLimitNumber {
		t.Fatalf("TestOnPayListUpdateRequest: 数据不符合预期，resp2:%v\n", resp2)
	}
	return
}

// TestOnPayListCopyRequest 测试复制功能
func TestOnPayListCopyRequest(t *testing.T) {
	uri := "/v1/pay/list/copy"
	param := handlers.PayListCopyReq{
		IDList:   []int{payMap[gameID][0].ID},
		ToGameID: int(70),
	}
	var resp CommonResp
	err := unitTest.TestHandlerUnMarshalResp("POST", uri, "json", param, &resp)
	if err != nil {
		t.Fatalf("TestOnPayListCopyRequest: 请求出错，err:%v\n", err)
	}
	if resp.Errno != "0" {
		t.Fatalf("TestOnPayListCopyRequest: 请求出错, 不符合预期，resp:%v\n", resp)
	}

	uri2 := "/v1/pay/list/list"
	var resp2 OnPayListListRequestResp
	param2 := &handlers.OnPayListListReq{
		GameID:     70,
		GameAreaID: -1,
	}
	err = unitTest.TestHandlerUnMarshalResp("POST", uri2, "json", param2, &resp2)
	if err != nil {
		t.Fatalf("TestOnPayListCopyRequest: 请求出错，err:%v\n", err)
	}
	if resp2.Errno != "0" {
		t.Fatalf("TestOnPayListCopyRequest: 请求出错, 不符合预期，resp2:%v\n", resp2)
	}
	if len(resp2.Data) <= 0 {
		t.Fatalf("TestOnPayListCopyRequest: 数据不符合预期，resp2:%v\n", resp2)
	}
	return
}

func TestTruncateInfo(t *testing.T) {
	return
	if _, err := clubdb.Exec("truncate room_card_product_list"); err != nil {
		t.Fatalf("TestTruncateInfo: 清空数据出错, err:%v\n", err.Error())
	}
	defer clubdb.Close()
}
