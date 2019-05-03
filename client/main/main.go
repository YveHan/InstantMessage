package main

import (
	"fmt"
	"os"
	"serverClient/client/process"
)

var userId int
var userPwd string
var userName string

func main() {
	var key int
	for true {//無限循環，用於輸入不合法數字時，繼續顯示菜單
		fmt.Println("——————————歡迎登陸多人聊天系統——————————")
		fmt.Println("\t\t\t 1 登錄聊天室")
		fmt.Println("\t\t\t 2 註冊用戶")
		fmt.Println("\t\t\t 3 退出系統")
		fmt.Println("\t\t\t 請選擇（1-3）")
		fmt.Scanf("%d\n",&key)
		switch key {
		case 1:
			fmt.Println("登錄聊天室")
			fmt.Println("請輸入用戶名")
			fmt.Scanf("%d\n",&userId)
			fmt.Println("請輸入用戶密碼")
			fmt.Scanf("%s\n",&userPwd)
			up := &process.UserProcess{}
			up.Login(userId,userPwd)
		case 2:
			fmt.Println("註冊用戶")
			fmt.Println("請輸入用戶id")
			fmt.Scanf("%d\n",&userId)
			fmt.Println("請輸入用戶密碼")
			fmt.Scanf("%s\n",&userPwd)
			fmt.Println("請輸入用戶名稱")
			fmt.Scanf("%s\n",&userName)
			up := &process.UserProcess{}
			up.Register(userId,userPwd,userName)
		case 3:
			fmt.Println("退出系統")
			os.Exit(0)
		default:
			fmt.Println("對不起，你的輸入有誤")
		}
	}
}
