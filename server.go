package main

import (
	"fmt"
	"github.com/dokidokikoi/my-zinx/ziface"
	"github.com/dokidokikoi/my-zinx/znet"
	"mmo/api"
	"mmo/core"
)

func main() {
	// 创建服务器句柄
	s := znet.NewServer()

	// 注册客户端连接建立和丢失函数
	s.SetOnConnStart(OnConnectionAdd)
	// 注册客户端连接丢失函数
	s.SetOnConnStop(OnConnectionLost)

	s.AddRouter(2, &api.WorldChatApi{})
	s.AddRouter(3, &api.MoveApi{})

	// 启动服务
	s.Serve()
}

func OnConnectionAdd(conn ziface.IConnection) {
	// 创建一个玩家
	player := core.NewPlayer(conn)

	// 同步当前的 PlayerID，使用 MsgID=1 的消息
	player.SyncPid()

	// 同步当前玩家的初始化坐标消息，使用 MsgID=200 的消息
	player.BroadCastStartPosition()

	// 将当前新上线玩家添加到 worldManager 中
	core.WorldMgrObj.AddPlayer(player)

	// 将连接绑定 pid
	conn.SetProperty("pid", player.Pid)

	// 同步周边玩家上线消息
	player.SyncSurrounding()

	fmt.Println("====> Player pid=", player.Pid, "arrived ====")
}

func OnConnectionLost(conn ziface.IConnection) {
	pid, _ := conn.GetProperty("pid")
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	if pid != nil {
		player.LostConnection()
	}

	fmt.Println("====> Player", pid, "left ====")
}
