package core

import (
	"fmt"
	"github.com/dokidokikoi/my-zinx/ziface"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"mmo/pb"
	"sync"
)

var (
	PidGen int32      = 1 // 生成玩家 id 的计数器
	IdLock sync.Mutex     // 保护计算器的锁
)

// 玩家对象
type Player struct {
	Pid  int32              // 玩家id
	Conn ziface.IConnection // 当前玩家连接
	X    float32            // 平面 x 坐标
	Y    float32            // 高度
	Z    float32            // 平面 y 坐标
	V    float32            // 旋转角度（0~360）
}

func (p *Player) SendMsg(msgID uint32, data proto.Message) {
	// 将 proto Message 结构体格式化
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err:", err)
		return
	}

	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	// 调用 zinx 框架的 sendMsg 发包
	if err := p.Conn.SendMsg(msgID, msg); err != nil {
		fmt.Println("Player SendMsg error!")
		return
	}
}

// 告知客户端 pid，同步已经生成的玩家 id
func (p *Player) SyncPid() {
	// 组建 BroadCast 协议 Proto 数据
	data := &pb.SyncPid{
		Pid: p.Pid,
	}

	p.SendMsg(1, data)
}

// 广播玩家自己的地点
func (p *Player) BroadCastStartPosition() {
	// 组建 BroadCast 协议 Proto 数据
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			&pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.Y,
			},
		},
	}
	p.SendMsg(200, msg)
}

// 世界聊天
func (p *Player) Talk(content string) {
	// 组建 msgID=200 proto 数据
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	// 得到当前世界所有在线玩家
	players := WorldMgrObj.GetAllPlayers()

	// 向所有玩家发送 msgID=200 的消息
	for _, v := range players {
		v.SendMsg(200, msg)
	}
}

// 向当前玩家周边的（九宫格内）玩家广播自己的位置，让他们显示自己
func (p *Player) SyncSurrounding() {
	// 根据自己的位置，获取周围九宫格的玩家 pid
	pids := WorldMgrObj.AoiMgr.GetPIDsByPos(p.X, p.Z)
	// 获取 pid 对应玩家对象
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	// 组建 msgID=200 proto 数据
	msg := &pb.BroadCast{
		Pid: p.Pid,
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
	// 向当前玩家周边的（九宫格内）玩家发送自己的位置，出现在别人视野里
	for _, player := range players {
		player.SendMsg(200, msg)
	}

	// 让周围的玩家出现在自己的视野里
	playersData := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		playersData = append(playersData, &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		})
	}
	syncPlayersMsg := &pb.SyncPlayers{
		Ps: playersData[:],
	}
	p.SendMsg(202, syncPlayersMsg)
}

func (p *Player) GetSurroundingPlayers() []*Player {
	// 得到当前 aoi 区域的所有pid
	pids := WorldMgrObj.AoiMgr.GetPIDsByPos(p.X, p.Z)

	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	return players
}

// 更新玩家位置
func (p *Player) UpdatePos(x, y, z, v float32) {
	p.X, p.Y, p.Z, p.V = x, y, z, v
	// 组装 proto 消息
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	// 获取当前玩家周边全部玩家
	players := p.GetSurroundingPlayers()
	// 向周边的每个玩家发送 msgID=200 消息，更新移动位置消息
	for _, player := range players {
		player.SendMsg(200, msg)
	}
}

// 玩家下线
func (p *Player) LostConnection() {
	players := p.GetSurroundingPlayers()

	msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	for _, player := range players {
		player.SendMsg(201, msg)
	}

	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
	WorldMgrObj.RemovePlayer(p.Pid)
}

func NewPlayer(conn ziface.IConnection) *Player {
	// 生成一个 pid
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	p := &Player{
		Pid:  id,
		Conn: conn,
		// 随机在 160 坐标点，基于 x 轴偏移若干坐标
		X: float32(160 + rand.Intn(10)),
		Y: 0,
		// 随机在 134 坐标点，基于 y 轴偏移若干坐标
		Z: float32(134 + rand.Intn(17)),
		V: 0,
	}

	return p
}
