package apis

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/golang/protobuf/proto"
	"log"
	"mmoGame/core"
	"mmoGame/pb"
)

type WorldChatApi struct {
	znet.BaseRouter
}


func (*WorldChatApi) Handle(request ziface.IRequest) {
	//1. 将客户端传来的proto协议解码
	proto_msg :=&pb.Talk{}
	err:=proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		log.Println("Talk Unmarshal error ", err)
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
	//4. 让player对象发起聊天广播请求
	player.Talk(proto_msg.Content)
}