package apis

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/golang/protobuf/proto"
	"log"
	"mmoGame/core"
	"mmoGame/pb"
)

type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi)Handle(request ziface.IRequest)  {
	//1. 将客户端传来的proto协议解码
	msg:=&pb.Position{}
	err:=proto.Unmarshal(request.GetData(),msg)
	if err != nil {
		log.Println("Move: Position Unmarshal error ", err)
		return
	}
	//2. 得知当前的消息是从哪个玩家传递来的,从连接属性pID中获取
	pid,err:=request.GetConnection().GetProperty("pID")
	if err != nil {
		log.Println("GetProperty pID error", err)
		request.GetConnection().Stop()
		return
	}

	//3. 根据pID得到player对象
	player:=core.WorldMgrObj.GetPlayerByPID(pid.(int32))
	//4. 让player对象发起移动位置信息广播
	//player.SendMsg()
	player.UpdatePos(msg)
}
