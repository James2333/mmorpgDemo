package core

import (
	"github.com/aceld/zinx/ziface"
	"github.com/golang/protobuf/proto"
	"log"
	"math/rand"
	"mmoGame/pb"
	"sync"
	"time"
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



func (p *Player)SendMsg(msgID uint32,data proto.Message){
	//log.Printf("before Marshal data = %+v\n", data)
	//将proto Message结构体序列化
	msg, err := proto.Marshal(data)
	if err != nil {
		log.Println("marshal msg err: ", err)
		return
	}
	//log.Printf("after Marshal data = %+v\n", msg)

	if p.Conn == nil {
		log.Println("connection in player is nil")
		return
	}

	//调用Zinx框架的SendMsg发包
	if err := p.Conn.SendMsg(msgID, msg); err != nil {
		log.Println("Player SendMsg error !")
		return
	}

	return
}
func (p *Player)Talk(content string){
	//1. 组建MsgID200 proto数据
	msg:=&pb.BroadCast{
		PID:  p.PID,
		Tp:   1,
		Data: &pb.BroadCast_Content{Content: content},
	}
	//2. 得到当前世界所有的在线玩家
	players:=WorldMgrObj.Players
	//3. 向所有的玩家发送MsgID:200消息
	for _,v:=range players{
		v.SendMsg(200,msg)
	}
}

//给当前玩家周边的(九宫格内)玩家广播自己的位置，让他们显示自己
func (p *Player) SyncSurrounding() {
	//1 根据自己的位置，获取周围九宫格内的玩家pID
	pids:=WorldMgrObj.AoiMgr.GetPIDsByPos(p.X,p.Z)
	players:=make([]*Player,0,len(pids))
	for _,pid:=range pids{
		players=append(players,WorldMgrObj.GetPlayerByPID(int32(pid)))
	}
	//2 根据pID得到所有玩家对象
	//3 给这些玩家发送MsgID:200消息，让自己出现在对方视野中
	proto_msg:=&pb.BroadCast{
		PID:  p.PID,
		Tp:   2,
		Data: &pb.BroadCast_P{P: &pb.Position{
			X: p.X,
			Y: p.Y,
			Z: p.Z,
			V: p.V,
		}},
	}
	//3.1 组建MsgID200 proto数据
	//3.2 每个玩家分别给对应的客户端发送200消息，显示人物
	for _,player:=range players{
		player.SendMsg(200,proto_msg)
	}
	//4 让周围九宫格内的玩家出现在自己的视野中
	playerDate:=make([]*pb.Player,0,len(players))
	//4.1 制作Message SyncPlayers 数据
	for _,player:=range players{
		p:=&pb.Player{
			PID: player.PID,
			P:   &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		playerDate=append(playerDate,p)
	}
	//4.2 封装SyncPlayer protobuf数据
	syncplayers:=&pb.SyncPlayers{Ps: playerDate}
	//4.3 给当前玩家发送需要显示周围的全部玩家数据
	p.SendMsg(202,syncplayers)
}

func (p *Player) UpdatePos(ps *pb.Position) {
	oldGID := WorldMgrObj.AoiMgr.GetGIDByPos(p.X, p.Z)
	//计算新格子gID
	newGID := WorldMgrObj.AoiMgr.GetGIDByPos(ps.X, ps.Z)

	p.X = ps.X
	p.Y = ps.Y
	p.Z = ps.Z
	p.V = ps.V
	if oldGID != newGID {
		//触发gird切换
		//把pID从就的aoi格子中删除
		WorldMgrObj.AoiMgr.RemovePIDFromGrID(int(p.PID), oldGID)
		//把pID添加到新的aoi格子中去
		WorldMgrObj.AoiMgr.AddPIDToGrID(int(p.PID), newGID)

		_ = p.OnExchangeAoiGrID(oldGID, newGID)
	}
	//组装protobuf协议，发送位置给周围玩家
	msg := &pb.BroadCast{
		PID: p.PID,
		Tp:  4, //4- 移动之后的坐标信息
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	//获取当前玩家周边全部玩家
	players := p.GetSurroundingPlayers()
	//向周边的每个玩家发送MsgID:200消息，移动位置更新消息
	for _, player := range players {
		player.SendMsg(200, msg)
	}


}

//获得当前玩家的AOI周边玩家信息
func (p *Player) GetSurroundingPlayers() []*Player {
	//得到当前AOI区域的所有pID
	pIDs := WorldMgrObj.AoiMgr.GetPIDsByPos(p.X, p.Z)

	//将所有pID对应的Player放到Player切片中
	players := make([]*Player, 0, len(pIDs))
	for _, pID := range pIDs {
		players = append(players, WorldMgrObj.GetPlayerByPID(int32(pID)))
	}

	return players
}


func (p *Player) OnExchangeAoiGrID(oldGID, newGID int) error {
	//获取就的九宫格成员
	oldGrIDs := WorldMgrObj.AoiMgr.GetSurroundGrIDsByGID(oldGID)

	//为旧的九宫格成员建立哈希表,用来快速查找
	oldGrIDsMap := make(map[int]bool, len(oldGrIDs))
	for _, grID := range oldGrIDs {
		oldGrIDsMap[grID.GID] = true
	}

	//获取新的九宫格成员
	newGrIDs := WorldMgrObj.AoiMgr.GetSurroundGrIDsByGID(newGID)
	//为新的九宫格成员建立哈希表,用来快速查找
	newGrIDsMap := make(map[int]bool, len(newGrIDs))
	for _, grID := range newGrIDs {
		newGrIDsMap[grID.GID] = true
	}

	//------ > 处理视野消失 <-------
	offlineMsg := &pb.SyncPID{
		PID: p.PID,
	}

	//找到在旧的九宫格中出现,但是在新的九宫格中没有出现的格子
	leavingGrIDs := make([]*GrID, 0)
	for _, grID := range oldGrIDs {
		if _, ok := newGrIDsMap[grID.GID]; !ok {
			leavingGrIDs = append(leavingGrIDs, grID)
		}
	}

	//获取需要消失的格子中的全部玩家
	for _, grID := range leavingGrIDs {
		players := WorldMgrObj.GetPlayersByGID(grID.GID)
		for _, player := range players {
			//让自己在其他玩家的客户端中消失
			player.SendMsg(201, offlineMsg)

			//将其他玩家信息 在自己的客户端中消失
			anotherOfflineMsg := &pb.SyncPID{
				PID: player.PID,
			}
			p.SendMsg(201, anotherOfflineMsg)
			time.Sleep(200 * time.Millisecond)
		}
	}

	//------ > 处理视野出现 <-------

	//找到在新的九宫格内出现,但是没有在就的九宫格内出现的格子
	enteringGrIDs := make([]*GrID, 0)
	for _, grID := range newGrIDs {
		if _, ok := oldGrIDsMap[grID.GID]; !ok {
			enteringGrIDs = append(enteringGrIDs, grID)
		}
	}

	onlineMsg := &pb.BroadCast{
		PID: p.PID,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	//获取需要显示格子的全部玩家
	for _, grID := range enteringGrIDs {
		players := WorldMgrObj.GetPlayersByGID(grID.GID)

		for _, player := range players {
			//让自己出现在其他人视野中
			player.SendMsg(200, onlineMsg)

			//让其他人出现在自己的视野中
			anotherOnlineMsg := &pb.BroadCast{
				PID: player.PID,
				Tp:  2,
				Data: &pb.BroadCast_P{
					P: &pb.Position{
						X: player.X,
						Y: player.Y,
						Z: player.Z,
						V: player.V,
					},
				},
			}

			time.Sleep(200 * time.Millisecond)
			p.SendMsg(200, anotherOnlineMsg)
		}
	}

	return nil
}


func (p *Player) LostConnection() {
	//获取周边玩家
	players:=p.GetSurroundingPlayers()
	msg:=&pb.SyncPID{PID: p.PID}
	for _,player:=range players{
		player.SendMsg(201,msg)
	}

	//从世界管理器中删除玩家
	WorldMgrObj.AoiMgr.RemoveFromGrIDByPos(int(p.PID),p.X,p.Z)
	WorldMgrObj.RemovePlayerByPID(p.PID)
}

