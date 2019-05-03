package process

import (
	"encoding/json"
	"fmt"
	"net"
	"serverClient/common/message"
	"serverClient/common/utils"
	"time"
)

type UserProcess struct {
	Conn net.Conn
	UserId int
}

//用戶處理器的註冊函數
func (this *UserProcess) Register(userId int,userPwd, userName string) (err error) {
	//連接服務器
	conn, err := net.Dial("tcp","localhost:8888")
	if err != nil {
		fmt.Println("連接服務器錯誤")
		return
	}
	defer conn.Close()
	//定義消息
	var mes message.Message
	mes.Type = message.RegisterMesType
	//定義消息數據體
	var registerMes message.UserMes
	registerMes.User.UserId =userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName
	tf := &utils.Transfer{
		Conn:conn,
	}
	//發送給服務器註冊信息
	err = tf.WritePkg(&mes,&registerMes)
	mes, err = tf.ReadPkg()
	if err != nil {
		return
	}
	//取出數據體中的狀態碼並判斷
	var registerResMes message.ResMes
	if err = json.Unmarshal([]byte(mes.Data),&registerResMes); err != nil{
		fmt.Println("註冊響應數據體反序列化失敗")
		return
	}
	if registerResMes.Code == 200 {
		fmt.Println("註冊成功，重新登錄吧")
		//顯示提示信息一秒後結束函數返回到主界面
		time.Sleep(time.Second)
		return
	}else {
		fmt.Println(registerResMes.Error)
		//顯示提示信息一秒後結束函數返回到主界面
		time.Sleep(time.Second)
		return
	}
	return
}

//用戶處理器的登錄函數
func (this *UserProcess)Login(id int,pd string) (err error) {
	//連接客戶端
	conn, err := net.Dial("tcp","localhost:8888")
	if err != nil {
		fmt.Println("連接服務器失敗")
	}
	defer conn.Close()
	//定義消息
	var mes message.Message
	mes.Type = message.LoginMesType
	//定義數據體
	var loginMes message.UserMes
	loginMes.User.UserId = id
	loginMes.User.UserPwd = pd
	tf := &utils.Transfer{
		Conn:conn,
	}
	err = tf.WritePkg(&mes,&loginMes)
	mes, err = tf.ReadPkg()
	if err != nil {
		return
	}
	//定義響應消息，直接省略類型判斷
	var loginResMes message.ResMes
	if err = json.Unmarshal([]byte(mes.Data),&loginResMes); err != nil {
		fmt.Println("登錄響應數據體反序列化失敗")
		return
	}
	if loginResMes.Code == 200 {
		//保持連接
		this.Conn = conn
		this.UserId = id
		//將客戶端與服務器的連接傳遞給協程，讀取連接上的數據
		go severProcessMes(conn)
		//登錄成功後進入二級菜單,傳遞連接與自己的ID
		//每當一次二級菜單選擇結束，又會重新顯示
		for {
			ShowMenu(conn,id)
		}
	}else{
		fmt.Println(loginResMes.Error)
	}
	return
}

//用戶處理器的顯示在線用戶函數
func (this *UserProcess)GetOnlineUser() (err error) {
	var mes message.Message
	mes.Type = message.GetOnlineUserMesType
	//定義數據體
	var getOnlineUserMes message.UserMes
	//	//將處理器保存的連接傳遞給傳輸實體
	tf := &utils.Transfer{
		Conn:this.Conn,
	}
	err = tf.WritePkg(&mes,&getOnlineUserMes)
	//退出程序後沒有必要操作本地用戶在線列表了
	return
}

//用戶處理器的登出函數
func (this *UserProcess)Logout(conn net.Conn,id int) (err error) {
	var mes message.Message
	mes.Type = message.LogoutMesType
	//定義數據體
	var logoutMes message.UserMes
	//將處理器保存的用戶ID賦值給響應消息
	logoutMes.User.UserId = id
	//	//將處理器保存的連接傳遞給傳輸實體
	tf := &utils.Transfer{
		Conn:conn,
	}
	err = tf.WritePkg(&mes,&logoutMes)
	//退出程序後沒有必要操作本地用戶在線列表了
	return
}
