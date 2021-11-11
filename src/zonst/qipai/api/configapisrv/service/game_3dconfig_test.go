package service

///*func Test_On3DConfigListConfig(t *testing.T) {
//	cat := ashe.New(t).Use(ashe.Redis)
//	cat.UnitTesting(`TestFindByGameID`, func() {
//
//		cat.MiniUnitTesting(`key存在`, func() {
//			var c cache.ThreeDLobbyGameDisplay
//			conn := cat.GetRedisConn()
//			data := cat.GetRedisDataController()
//			defer data.FlushDB() // 清除所有数据，避免有其它单元测试影响
//
//			_ = data.Set(`three:lobby:display:1`, "22")
//			username, _, err := c.FindByGameID(conn, 1)
//			cat.Nil(err)
//			cat.Equal(username, `22`)
//		})
//		cat.MiniUnitTesting(`key不存在`, func() {
//
//		})
//	})
//}
//*/
