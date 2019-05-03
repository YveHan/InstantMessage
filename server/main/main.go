package main

import (
	"fmt"
	"net"
	"serverClient/server/model"
	"time"
)

//協程將連接傳遞給一個處理器實體，在協程退出時關閉客戶端連接
func process(conn net.Conn){
	defer conn.Close()
	//使用連接初始化一個處理器
	processor := &Processor{
		Conn:conn,
	}
	//處理器開始處理
	if err := processor.processing(); err !=nil {
		return
	}
}

func initUserDao()  {
	model.MyUserDao = model.NewUserDao(pool)
}

//主函數監聽端口並不斷接受客戶端到來的連接
func main() {
	//初始化redis連接池以提高效率
	initPool("localhost:6379",16,0,300*time.Second)
	//用初始化後的redis連接池初始化數據訪問對象
	initUserDao()
	fmt.Println("服務器在8888端口監聽……")
	listen, err := net.Listen("tcp","0.0.0.0:8888")
	if err != nil {
		fmt.Println("監聽端口出錯，",err)
		return
	}
	defer listen.Close()
	for {
		fmt.Println("等待客戶端連接……")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("接收客戶端請求出錯……")
		}
		//每個客戶端的連接啟動一個協程處理
		go process(conn)
	}
}