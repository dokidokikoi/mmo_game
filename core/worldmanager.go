package core

import "sync"

type WorldManager struct {
	AoiMgr  *AOIManager
	Players map[int32]*Player
	pLock   sync.RWMutex
}

var WorldMgrObj *WorldManager

func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()

	// 将玩家添加到 aoi 网格规划中
	wm.AoiMgr.Add2GridByPos(int(player.Pid), player.X, player.Z)
}

func (wm *WorldManager) RemovePlayer(pid int32) {
	wm.pLock.Lock()
	delete(wm.Players, pid)
	wm.pLock.Unlock()
}

func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pid]
}

func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player, 0)
	for _, v := range wm.Players {
		players = append(players, v)
	}

	return players
}

func init() {
	WorldMgrObj = &WorldManager{
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		Players: make(map[int32]*Player),
	}
}
