package subprocess

import (
	"encoding/json"
	"fmt"
	"net"
	"serverClient/common/message"
	"serverClient/common/utils"
)

//消息處理器
type SmsProcess struct {
	UserId int
}
//服務器轉發消息
func (this *SmsProcess)SendGroupMes(mes *message.Message)  {
	//定義消息數據體
	var smsMes message.SmsMes
	//需要解碼獲取用戶ID以便之後不給自己轉發消息
	if err := json.Unmarshal([]byte(mes.Data),&smsMes); err != nil {
		fmt.Println("反序列化失敗")
		return
	}
	if len(myUserMgr.GetAllOnlineUser()) == 1 {
		fmt.Println("群裡沒有人啊")
		return
	}
	//查詢所有用戶並排除自身，發送消息給其他用戶
	for id, up := range myUserMgr.GetAllOnlineUser() {
		if id == smsMes.User.UserId {
			continue
		}
		this.SendMesToOneOnlineUser(mes, &smsMes, up.Conn)
	}
	return
}
func (this *SmsProcess)SendPrivateMes(mes *message.Message)  {
	//定義消息數據體
	var smsMes message.SmsMes
	//需要解碼查找對應的目標
	if err := json.Unmarshal([]byte(mes.Data),&smsMes); err != nil {
		fmt.Println("反序列化失敗")
		return
	}
	up, err := myUserMgr.GetOnlineUserById(smsMes.To)
	if err != nil {
		fmt.Println(err)
		return
	}
	this.SendMesToOneOnlineUser(mes,&smsMes,up.Conn)
	return
}
func (this *SmsProcess)SendMesToOneOnlineUser(mes *message.Message,smsMes *message.SmsMes, conn net.Conn) {
	tf := &utils.Transfer{
		Conn:conn,
	}
	if err := tf.WritePkg(mes,smsMes); err != nil {
		fmt.Println("消息發送失敗")
	}
}