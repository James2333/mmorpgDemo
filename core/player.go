package core

import (
	"github.com/aceld/zinx/ziface"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"mmoGame/pb"
	"sync"
)

//玩家对象
type Player struct {
	PID  int32              //玩家ID
	Conn ziface.IConnection //当前玩家的连接
	X    float32            //平面x坐标
	Y    float32            //高度
	Z    float32            //平面y坐标 (注意不是Y)
	V    float32            //旋转0-360度
}

/*
	Player ID 生成器
*/
var PIDGen int32 = 1  //用来生成玩家ID的计数器
var IDLock sync.Mutex //保护PIDGen的互斥机制

//创建一个玩家对象
func NewPlayer(conn ziface.IConnection) *Player {
	//生成一个PID
	IDLock.Lock()
	ID := PIDGen
	PIDGen++
	IDLock.Unlock()

	p := &Player{
		PID:  ID,
		Conn: conn,
		X:    float32(160 + rand.Intn(50)), //随机在160坐标点 基于X轴偏移若干坐标
		Y:    0,                            //高度为0
		Z:    float32(134 + rand.Intn(50)), //随机在134坐标点 基于Y轴偏移若干坐标
		V:    0,                            //角度为0，尚未实现
	}

	return p
}

//告知客户端pID,同步已经生成的玩家ID给客户端
func (p *Player) SyncPID() {
	//组建MsgID0 proto数据
	data := &pb.SyncPID{
		PID: p.PID,
	}

	//发送数据给客户端
	p.SendMsg(1, data)
}

//广播玩家自己的出生地点
func (p *Player) BroadCastStartPosition() {
	//bc:=&pb.BroadCast{
	//	PID:  0,
	//	Tp:   0,
	//	Data: &pb.BroadCast_P{P:&pb.Position{
	//		X: 0,
	//		Y: 0,
	//		Z: 0,
	//		V: 0,
	//	}},
	//}
	//组建MsgID200 proto数据
	msg := &pb.BroadCast{
		PID: p.PID,
		Tp:  2, //TP2 代表广播坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//发送数据给客户端
	p.SendMsg(200, msg)
}



func (p *Player)SendMsg(msgId uint32,data proto.Message){
	//将proto message序列化
	msg,err:=proto.Marshal(data)
	if err != nil {
		return
	}
	//将二进制文件通过zinx链接里面的send方法发送给客户端。
	if p.Conn==nil{
		return
	}
	if err=p.Conn.SendMsg(msgId,msg);err!=nil{
		return
	}

}
