package process

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"serverClient/common/message"
	"serverClient/common/utils"
)

//登錄成功後的二級菜單
func ShowMenu(conn net.Conn,id int) {
	fmt.Println("————恭喜登錄成功————")
	fmt.Println("————1、顯示用戶列表————")
	fmt.Println("————2、群發消息————")
	fmt.Println("————3、私發消息————")
	fmt.Println("————4、退出系統————")
	fmt.Println("————請選擇1-4————")
	var key int
	var content string
	var toid int
	//初始化群發消息處理器
	smsProcess := &SmsProcess{
		Conn:conn,
		UserId:id,
	}
	fmt.Scanf("%d\n",&key)
	switch key {
	case 1:
		up := &UserProcess{
			Conn:conn,
		}
		up.GetOnlineUser()
	case 2:
		fmt.Println("請輸入你想對大家說的話")
		fmt.Scanf("%s\n",&content)
		//群發消息需要連接、ID以及內容信息
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("請輸入私聊對象id")
		fmt.Scanf("%d\n",&toid)
		fmt.Println("請輸入私聊內容")
		fmt.Scanf("%s\n",&content)
		//私聊消息需要連接、ID以及內容信息
		smsProcess.SendPrivateMes(content,toid)
	case 4:
		fmt.Println("退出系統")
		up := &UserProcess{}
		up.Logout(conn,id)
		os.Exit(0)
	default:
		fmt.Println("輸入選項不正確")
	}
}

//維持服務端的通信
func severProcessMes(Coon net.Conn) {
	tf := &utils.Transfer{
		Conn:Coon,
	}
	for {
		mes, err := tf.ReadPkg()
		if err != nil {
			return
		}
		switch mes.Type {
		//接收其他用戶在線通知
		case message.NotifyUserStatusMesType:
			//定義幾個通知數據體
			var notifyUserStatusMes message.NotifyUserStatusMes
			if err := json.Unmarshal([]byte(mes.Data),&notifyUserStatusMes);err !=nil {
				return
			}
			fmt.Printf("用戶id:%d 上線了！",notifyUserStatusMes.UserId)
		case message.GetOnlineUserMesResType:
			var getOnlineUserMes message.ResMes
			if err := json.Unmarshal([]byte(mes.Data),&getOnlineUserMes); err != nil {
				return
			}
			fmt.Println("當前在線用戶：")
			for _,val := range getOnlineUserMes.UsersId {
				fmt.Println("用戶id：",val)
			}
		case message.GroupSmsMesType:
			//接收服務端轉發消息
			outputGroupMes(&mes)
		case message.PrivateSmsMesType:
			outputPrivateMes(&mes)
		default:
			fmt.Println("未知的消息")
		}
	}
}
