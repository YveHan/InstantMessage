package process

import (
	"fmt"
	"net"
	"serverClient/common/message"
	"serverClient/common/utils"
)

//消息處理器
type SmsProcess struct {
	Conn net.Conn
	UserId int
}

//發送給其他客戶端的函數
func (this *SmsProcess)SendGroupMes(content string) (err error) {
	//定義消息
	var mes message.Message
	mes.Type = message.GroupSmsMesType
	//定義數據體
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.User.UserId = this.UserId
	tf := &utils.Transfer{
		Conn:this.Conn,
	}
	if err = tf.WritePkg(&mes,&smsMes);err != nil {
		fmt.Println("發送消息出錯了")
		return
	}
	return
}

func (this *SmsProcess)SendPrivateMes(content string,id int) (err error) {
	//定義消息
	var mes message.Message
	mes.Type = message.PrivateSmsMesType
	//定義數據體
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.User.UserId = this.UserId
	smsMes.To = id
	tf := &utils.Transfer{
		Conn:this.Conn,
	}
	if err = tf.WritePkg(&mes,&smsMes);err != nil {
		fmt.Println("發送消息出錯了")
		return
	}
	return
}