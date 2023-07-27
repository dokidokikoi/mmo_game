package api

import (
	"fmt"
	"github.com/dokidokikoi/my-zinx/ziface"
	"github.com/dokidokikoi/my-zinx/znet"
	"google.golang.org/protobuf/proto"
	"mmo/core"
	"mmo/pb"
)

type WorldChatApi struct {
	znet.BaseRouter
}

func (*WorldChatApi) Handle(req ziface.IRequest) {
	// 1.将客户端传来的 proto 解码
	msg := &pb.Talk{}
	err := proto.Unmarshal(req.GetData(), msg)
	if err != nil {
		fmt.Println("Talk Unmarshall error:", err)
		return
	}

	// 2.得知当前的消息是从哪个玩家传递过来的，从连接属性 pid 获取
	pid, err := req.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error:", err)
		req.GetConnection().Stop()
		return
	}

	// 3.根据 pid 得到 player 对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	// 4.让 player 发起聊天广播请求
	player.Talk(msg.Content)
}
