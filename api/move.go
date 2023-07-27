package api

import (
	"fmt"
	"github.com/dokidokikoi/my-zinx/ziface"
	"github.com/dokidokikoi/my-zinx/znet"
	"google.golang.org/protobuf/proto"
	"mmo/core"
	"mmo/pb"
)

type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(req ziface.IRequest) {
	// 将客户端传来的消息解码
	msg := &pb.Position{}
	err := proto.Unmarshal(req.GetData(), msg)
	if err != nil {
		fmt.Println("Move: Position Unmarshal error:", err)
		return
	}

	// 获取玩家id
	pid, err := req.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error:", err)
		req.GetConnection().Stop()
		return
	}

	fmt.Printf("player id = %d, move(%f,%f,%f,%f)\n", pid, msg.X, msg.Y, msg.Z, msg.V)

	// 让 player 发起移动位置广播
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	player.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
}
