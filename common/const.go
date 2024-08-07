package common

var EmptyData = make(map[string]interface{})

// 存储玩家session的结构体
type LoginInfo struct {
	UserId   int
	UserName string
}
