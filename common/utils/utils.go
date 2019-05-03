package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"serverClient/common/message"
)

//傳輸實體，包括客戶端連接與數據長度
type Transfer struct {
	Conn net.Conn
	Buf [8096]byte
}

//讀取客戶端的數據
func (this *Transfer)ReadPkg()(mes message.Message, err error){
	//先讀取數據包長度，4個字節長,
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil{
		fmt.Println("讀取數據出錯了，",err)
		return
	}
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4])
	fmt.Println(pkgLen)
	//讀取客戶端傳過來的特定長度數據
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	//如果數據長度不等於預期也結束函數
	if n != int(pkgLen)||err != nil {
		fmt.Println("讀取的數據包長度不匹配或",err)
		return
	}
	//首次反序列化客戶端消息並賦值給消息實例
	if err = json.Unmarshal(this.Buf[:pkgLen],&mes); err != nil {
		fmt.Println("數據首次反序列化失敗",err)
		return
	}
	return
}
func (this *Transfer)WritePkg(mes *message.Message,resmes interface{}) (err error) {
	//先對數據體序列化
	data, err := json.Marshal(resmes)
	if err != nil {
		fmt.Println("數據體序列化失敗")
		return
	}
	mes.Data = string(data)
	//對消息序列化
	data, err = json.Marshal(mes)
	if err !=  nil {
		fmt.Println("消息序列化失敗")
		return
	}
	var pkgLen uint32
	//計算消息長度並先發送數據長度
	pkgLen = uint32(len(data))
	binary.BigEndian.PutUint32(this.Buf[:4],pkgLen)
	n,err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("發送錯誤")
		return
	}
	//發送消息
	n,err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("寫入數據失敗")
		return
	}
	return
}