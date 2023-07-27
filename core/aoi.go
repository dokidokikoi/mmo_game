package core

import "fmt"

const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)

type AOIManager struct {
	MinX  int           // 区域左边界坐标
	MaxX  int           // 区域右边界坐标
	CntsX int           // x 方向格子的数量
	MinY  int           // 区域上边界坐标
	MaxY  int           // 区域下边界坐标
	CntsY int           // y 方向格子的数量
	grids map[int]*Grid // 当前区域中都有那些格子，key=格子id，val=格子对象
}

// 根据格子的 gID 得到当前周边的九宫格信息
func (m *AOIManager) GetSurroundGridsByGid(gid int) (grids []*Grid) {
	// 判断 gid 是否存在
	if _, ok := m.grids[gid]; !ok {
		return
	}

	// 将当前 gid 添加到九宫格中
	grids = append(grids, m.grids[gid])

	// 根据 gid 得到当前格子所在的 x 轴编号
	idx := gid % m.CntsX

	// 判断当前 idx 的左边是否还有格子
	if idx > 0 {
		grids = append(grids, m.grids[gid-1])
	}
	// 判断当前 idx 的右边是否还有格子
	if idx < m.CntsX-1 {
		grids = append(grids, m.grids[gid+1])
	}

	// 将 x 轴的格子都取出，进行遍历，再分别得到每个格子的上下是否有格子
	gridsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gridsX = append(gridsX, v.GID)
	}

	// 遍历 x 轴格子
	for _, v := range gridsX {
		// 计算该格子处于第几列
		idy := v / m.CntsX
		// 判断当前 idy 的上面是否还有格子
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CntsX])
		}
		// 判断当前 idy 的下面是否还有格子
		if idy < m.CntsY-1 {
			grids = append(grids, m.grids[v+m.CntsX])
		}
	}
	return
}

// 通过横纵坐标获取格子的 id
func (m *AOIManager) GetGIDByPos(x, y float32) int {
	gx := (int(x) - m.MinX) / m.gridWidth()
	gy := (int(y) - m.MinY) / m.gridLength()

	return gy*m.CntsX + gx
}

// 通过横纵坐标得到周边九宫格内全部 playerIDs
func (m *AOIManager) GetPIDsByPos(x, y float32) (playerIDs []int) {
	gID := m.GetGIDByPos(x, y)

	grids := m.GetSurroundGridsByGid(gID)
	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
		fmt.Printf("===> grid ID: %d, pids: %v ====", v.GID, v.GetPlayerIDs())
	}

	return
}

// 根据 gid 获取该格子所有玩家 id
func (m *AOIManager) GetPidsByGid(gid int) (playerIDs []int) {
	playerIDs = m.grids[gid].GetPlayerIDs()
	return
}

// 移除一个格子中的玩家 id
func (m *AOIManager) RemovePidFromGrid(pid, gid int) {
	m.grids[gid].Remove(pid)
}

// 将玩家 id 添加到一个格子中
func (m *AOIManager) AddPid2Grid(pid, gid int) {
	m.grids[gid].Add(pid)
}

// 通过横纵坐标将玩家 id 添加到一个格子中
func (m *AOIManager) Add2GridByPos(pid int, x, y float32) {
	gid := m.GetGIDByPos(x, y)
	grid := m.grids[gid]
	grid.Add(pid)
}

// 通过横纵坐标移除一个格子中的玩家 id
func (m *AOIManager) RemoveFromGridByPos(pid int, x, y float32) {
	gid := m.GetGIDByPos(x, y)
	grid := m.grids[gid]
	grid.Remove(pid)
}

// 每个格子在 x 轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// 每个格子在 x 轴方向的长度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// 打印信息方法
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager: \nminX: %d, maxX: %d, cntsX: %d, minY: %d, maxY: %d, cnts: %d\n Grids in AOI Manager:\n",
		m.MinX, m.MaxY, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	// 给 aoi 初始化区域中所有的格子
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// 计算格子 ID
			// 格子编号：id = idy * nx + idx
			gid := y*cntsX + x

			// 初始化一个格子，然后放在 AOI 中的 map 里。key是当前格子的 ID
			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinY+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY*y*aoiMgr.gridLength(),
				aoiMgr.MinY*(y+1)*aoiMgr.gridLength(),
			)
		}
	}
	return aoiMgr
}
