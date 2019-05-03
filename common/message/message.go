package message

const (
	LoginMesType = "LoginMes"
	LoginResMesType = "LoginResMes"
	LogoutMesType  = "LogoutMes"
	RegisterMesType = "RegisterMes"
	RegisterResMesType = "RegisterResMes"
	GetOnlineUserMesType = "GetOnlineUserMes"
	GetOnlineUserMesResType = "GetOnlineUserResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	GroupSmsMesType = "GroupSmsMes"
	PrivateSmsMesType = "PrivateSmsMesType"
)

const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}
type UserMes struct {//客戶端註冊、登錄或註銷消息
	User User `json:"user"`
}

type ResMes struct {//服務器響應消息
	Code int `json:"code"`
	Error string `json:"error"`
	UsersId []int
}

type NotifyUserStatusMes struct {//通知在線的消息
	UserId int `json:"userId"`
	Status int `json:"status"`
}

type SmsMes struct {//群發、私聊的消息
	Content string `json:"content"`
	To int
	User User `json:"user"`
}