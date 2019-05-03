package subprocess

import "fmt"

var myUserMgr *userMgr

//用戶管理，使用存在的連接來標識在線用戶
type userMgr struct {
	onlineUsers map[int]*UserProcess
}

//初始化支持10個用戶同時在線，注意map可以動態增長
func init()  {
	myUserMgr = &userMgr{
		onlineUsers:make(map[int]*UserProcess,10),
	}
}

//使用客戶端連接的ID作為在線用戶的map的key值，客戶端連接ID源於用戶ID
func (this *userMgr) AddOnlineUser(up *UserProcess)  {
	this.onlineUsers[up.UserId] = up
}

//根據ID刪除某個在線用戶
func (this *userMgr)DelOnlineUser(userId int)  {
	delete(this.onlineUsers,userId)
}

//獲取所有在線用戶
func (this *userMgr)GetAllOnlineUser() map[int]*UserProcess  {
	return this.onlineUsers
}

//獲取指定ID的在線用戶
func (this *userMgr)GetOnlineUserById(userId int) (up *UserProcess,err error) {
	up, ok := this.onlineUsers[userId]
	if !ok {
		err = fmt.Errorf("用戶%d不存在或不在線",userId)
		return
	}
	return
}
