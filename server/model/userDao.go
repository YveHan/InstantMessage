package model

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"serverClient/common/message"
)

//初始化使用的全局變量
var MyUserDao *userDao

//私有化
type userDao struct {
	pool *redis.Pool
}

//單例模式，保證全局唯一性
func NewUserDao(pool *redis.Pool) (UserDao *userDao) {
	UserDao = &userDao{
		pool:pool,
	}
	return
}

//根據ID通過redis連接池中的一個連接查詢用戶
func getUserById(conn redis.Conn,id int)(user *message.User, err error){
	//注意redis語法 command + key + filed + value
	res, err := redis.String(conn.Do("HGet","users",id))
	if err != nil {
		//判斷查詢的用戶是否存在，不存在就賦一個錯誤值
		if err == redis.ErrNil {
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	//實例化一個用戶
	user = &message.User{}
	//將查詢的用戶數據賦值給用戶實例
	if err = json.Unmarshal([]byte(res),user);err != nil {
		fmt.Println("解碼錯誤")
		return
	}
	return
}

//與連接池交互的登錄操作
func (this *userDao)Login(userId int, userPwd string) (user *message.User,err error) {
	//獲取一個連接
	conn := this.pool.Get()
	defer conn.Close()
	//查詢對應ID的用戶
	user, err = getUserById(conn,userId)
	if err != nil {
		return
	}
	//將用戶輸入的密碼與從數據庫獲取用戶的密碼比較，如果輸錯就賦值一個錯誤
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

//與連接池交互的註冊操作
func (this *userDao)Register(user *message.User) (err error ) {
	//獲取一個連接
	conn := this.pool.Get()
	defer conn.Close()
	//查詢對應的ID，防止用戶ID重複
	_, err = getUserById(conn,user.UserId)
	if err !=  nil {
		//如果錯誤等於用戶不存在，數據庫將可以操作
		if err == ERROR_USER_NOTEXISTS {
			//將註冊信息數據序列化
			//注意這裡為了避免外面的err覆蓋裡面的err，所以作err1
			data, err1 := json.Marshal(user)
			if err1 != nil {
				fmt.Println("序列化註冊信息失敗",err)
				return
			}
			//將序列化後的信息通過redis連接保存到數據庫中
			if _, err = conn.Do("HSet","users",user.UserId,string(data)); err != nil {
				fmt.Println("保存用戶信息錯誤",err)
				return
			}
			return
		}else {
			fmt.Println("註冊操作發生了未知錯誤")
			return
		}
	}
	//查詢操作錯誤為空表示用戶已經存在了
	err = ERROR_USER_EXISTS
	return
}