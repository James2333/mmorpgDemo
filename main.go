package main

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"mmoGame/core"
)

func OnconnectionAdd(conn ziface.IConnection) {
	//创建一个玩家
	player:=core.NewPlayer(conn)
	//将玩家同步到客户端
	player.SyncPID()
	player.BroadCastStartPosition()

}
func main() {
	//1 创建一个server句柄
	s := znet.NewServer()

	//2 配置路由
	//s.AddRouter(0, &PingRouter{})
	s.SetOnConnStart(OnconnectionAdd)
	//3 开启服务
	s.Serve()
}
