package process

import (
	"encoding/json"
	"fmt"
	"serverClient/common/message"
)

func outputGroupMes(mes *message.Message) {
	var smsMes message.SmsMes
	if err := json.Unmarshal([]byte(mes.Data),&smsMes); err != nil {
		fmt.Println("解碼失敗")
		return
	}
	fmt.Printf("用戶id:\t%d 對大家說：\t%s",smsMes.User.UserId,smsMes.Content)
	fmt.Println()
}

func outputPrivateMes(mes *message.Message) {
	var smsMes message.SmsMes
	if err := json.Unmarshal([]byte(mes.Data),&smsMes); err != nil {
		fmt.Println("解碼失敗")
		return
	}
	fmt.Printf("用戶id:\t%d 對你說：\t%s",smsMes.User.UserId,smsMes.Content)
	fmt.Println()
}