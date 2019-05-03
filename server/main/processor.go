package main

import (
	"fmt"
	"io"
	"net"
	"serverClient/common/message"
	"serverClient/common/utils"
	"serverClient/server/process"
)

//定義主處理器
type Processor struct {
	Conn net.Conn
}


//主處理器的處理,將出錯的消息返回給上級打印
func (this *Processor) processing() (err error) {
	//處理器使用連接初始化一個傳輸實體
	for {
		tf :=&utils.Transfer{
			Conn:this.Conn,
		}
		//傳輸實體不斷從連接上讀取客戶端的數據
		mes, err :=tf.ReadPkg()
		if err != nil{
			if err ==io.EOF {
				fmt.Println("客戶端退出！")
			}
			return err
		}
		fmt.Println("讀取的數據為，",mes)
		//讀取客戶端消息後交給服務器消息處理函數
		if err = this.serverProcessMes(&mes); err != nil {
			return err
		}
	}
}

//服務器消息分類處理
func (this *Processor)serverProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	//客戶端發送的是登錄消息，實例化子處理器——用戶處理器調用登錄函數
	case message.LoginMesType:
		up :=&subprocess.UserProcess{
			Conn:this.Conn,
		}
		err = up.ServerProcessLogin(mes)
	//客戶端發送的是登錄消息，實例化子處理器——用戶處理器調用註冊函數
	case message.RegisterMesType:
		up :=&subprocess.UserProcess{
			Conn:this.Conn,
		}
		err = up.ServerProcessRegister(mes)
	//客戶端發送的是群聊消息，實例化子處理器——消息處理器調用群聊消息函數
	case message.GroupSmsMesType:
		smsProcess := &subprocess.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	case message.PrivateSmsMesType:
		smsProcess := &subprocess.SmsProcess{}
		smsProcess.SendPrivateMes(mes)
	case message.GetOnlineUserMesType:
		up :=&subprocess.UserProcess{
			Conn:this.Conn,
		}
		up.ServerProcessGetOnlienUser()
	case message.LogoutMesType:
		up :=&subprocess.UserProcess{
			Conn:this.Conn,
		}
		up.ServerProcessLogout(mes)
	//客戶端發送的是消息無法識別，提示服務器
	default:
		fmt.Println("消息類型不存在，無法處理")
	}
	return
}
