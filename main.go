package main

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"log"
	"mmoGame/apis"
	"mmoGame/core"
)

func OnconnectionAdd(conn ziface.IConnection) {
	//创建一个玩家
	player := core.NewPlayer(conn)

	//同步当前的PlayerID给客户端， 走MsgID:1 消息
	player.SyncPID()

	//同步当前玩家的初始化坐标信息给客户端，走MsgID:200消息
	player.BroadCastStartPosition()

	//将当前新上线玩家添加到worldManager中
	core.WorldMgrObj.AddPlayer(player)

	//将该连接绑定属性PID
	conn.SetProperty("pID", player.PID)

	//同步周边玩家上线信息，与现实周边玩家信息
	player.SyncSurrounding()

	log.Println("=====> Player pIDID = ", player.PID, " arrived ====")

}

//当客户端断开连接的时候的hook函数
func OnConnectionLost(conn ziface.IConnection) {
	//获取当前连接的PID属性
	pID, err := conn.GetProperty("pID")
	if err != nil {
		log.Println("Pid is nil ")
		return
	}
	//根据pID获取对应的玩家对象
	player := core.WorldMgrObj.GetPlayerByPID(pID.(int32))

	//触发玩家下线业务
	if pID != nil {
		player.LostConnection()
	}

	log.Println("====> Player ", pID, " left =====")

}


func main() {
	//1 创建一个server句柄
	s := znet.NewServer()

	//2 配置路由
	//s.AddRouter(0, &PingRouter{})
	s.SetOnConnStart(OnconnectionAdd)
	s.SetOnConnStop(OnConnectionLost)
	s.AddRouter(2,&apis.WorldChatApi{})
	s.AddRouter(3,&apis.MoveApi{})
	//3 开启服务
	s.Serve()
}
