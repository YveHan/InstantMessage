package subprocess

import (
	"encoding/json"
	"fmt"
	"net"
	"serverClient/common/message"
	"serverClient/server/model"
	"serverClient/common/utils"
)

//用戶處理器，包括從主處理器傳遞的連接與用於識別用戶的ID
type UserProcess struct {
	Conn net.Conn
	UserId int
}

//服務器處理註冊的函數
func (this *UserProcess)ServerProcessRegister(mes *message.Message) (err error) {
	//定義一個註冊消息數據體
	var registerMes message.UserMes
	//將客戶端註冊數據賦值給定義的數據體
	if err = json.Unmarshal([]byte(mes.Data),&registerMes); err != nil {
		fmt.Println("反序列化失敗")
		return
	}
	//定義響應消息
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	//定義響應消息的數據體
	var registerResMes message.ResMes
	if err = model.MyUserDao.Register(&registerMes.User); err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 401
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		}else {
			registerResMes.Code = 500
			registerResMes.Error = "服務器內部錯誤-註冊"
		}
	}else {
		registerResMes.Code = 200
	}
	//用戶處理器的連接傳遞給傳輸實體
	tf := &utils.Transfer{
		Conn:this.Conn,
	}
	err = tf.WritePkg(&resMes,&registerResMes)
	return
}

//服務器處理登錄的函數，具體操作的錯誤可以在這裡打印出來並返回給上級，調用函數的錯誤直接返回
func (this *UserProcess)ServerProcessLogin(mes *message.Message) (err error) {
	var loginMes message.UserMes
	//客戶端數據第二反序列化，將消息數據主體賦值給具體消息類型
	if err = json.Unmarshal([]byte(mes.Data),&loginMes); err != nil {
		fmt.Println("客戶端登錄數據第二反序列化失敗")
		return
	}
	//定義響應消息
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	//定義響應消息的數據體
	var loginResMes message.ResMes
	//與數據庫中的儲存對象比對
	_, err = model.MyUserDao.Login(loginMes.User.UserId,loginMes.User.UserPwd)
	if err != nil{
		//用戶未授權，用戶ID或密碼錯誤
		loginResMes.Code = 401
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Error = "該用戶不存在"
		}else if err == model.ERROR_USER_PWD {
			loginResMes.Error = "密碼錯誤"
		}else {
			//服務器內部遇到了一個錯誤
			loginResMes.Code = 500
			loginResMes.Error = "服務器內部錯誤-登錄"
		}
	}else {
		//成功處理請求
		loginResMes.Code = 200
		//登錄成功後添加自己在在線用戶列表中
		this.UserId = loginMes.User.UserId
		myUserMgr.AddOnlineUser(this)
		//通知其他用戶自己上線了
		this.NotifyOthersOnlineUser(loginMes.User.UserId)

		fmt.Println(this.UserId,"登錄成功")
	}
	//用戶處理器的連接傳遞給傳輸實體
	tf := &utils.Transfer{
		Conn:this.Conn,
	}
	err = tf.WritePkg(&resMes,&loginResMes)
	return
}


//服務器處理獲取在線用戶的函數
func (this *UserProcess)ServerProcessGetOnlienUser() (err error) {
	//定義響應消息
	var resMes message.Message
	resMes.Type = message.GetOnlineUserMesResType
	//定義響應消息的數據體
	var getOnlineUserResMes message.ResMes
	getOnlineUserResMes.Code = 200
	//將當前在線用戶的ID他添加到成功登錄後的響應消息中
	for id, _ := range myUserMgr.GetAllOnlineUser(){
		getOnlineUserResMes.UsersId = append(getOnlineUserResMes.UsersId,id)
	}
	//用戶處理器的連接傳遞給傳輸實體
	tf := &utils.Transfer{
		Conn:this.Conn,
	}
	err = tf.WritePkg(&resMes,&getOnlineUserResMes)
	return
}

//服務器處理用戶登出的的函數
func (this *UserProcess)ServerProcessLogout(mes *message.Message) (err error) {
	//反序列化消息
	var logoutMes message.UserMes
	err = json.Unmarshal([]byte(mes.Data),&logoutMes)
	//從在線用戶列表中刪除登出用戶
	myUserMgr.DelOnlineUser(logoutMes.User.UserId)
	return
}

func (this *UserProcess)NotifyOthersOnlineUser(userId int)  {
	//注意up是地址，指向userprocess
	for id, up := range myUserMgr.GetAllOnlineUser() {
		//自己不通知
		if id == userId {
			continue
		}
		//用查詢出的用戶連接調用通知方法告知上線的客戶端ID
		up.NotifyMeOnline(userId)
	}
}

//回傳一個消息告訴其他客戶端
func (this *UserProcess)NotifyMeOnline(userId int)  {
	//定義消息類型
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	//定義數據體
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline
	tf := &utils.Transfer{
		Conn:this.Conn,
	}
	if err := tf.WritePkg(&mes,&notifyUserStatusMes); err != nil {
		fmt.Println("通知在線消息出錯了")
		return
	}
}