package g2d

import (
	"github.com/kercylan98/minotaur/utils/hash"
	"github.com/kercylan98/minotaur/utils/synchronization"
	"sync"
)

func NewRadiationPattern[ItemType comparable, Item RadiationPatternItem[ItemType]](matrix [][]Item) *RadiationPattern[ItemType, Item] {
	var clone = make([][]Item, len(matrix))
	for x := 0; x < len(matrix); x++ {
		ys := make([]Item, len(matrix[0]))
		for y := 0; y < len(matrix[0]); y++ {
			ys[y] = matrix[x][y]
		}
		clone[x] = ys
	}
	rp := &RadiationPattern[ItemType, Item]{
		matrix:    clone,
		links:     synchronization.NewMap[int64, map[int64]bool](),
		positions: map[int64][2]int{},
	}
	for x := 0; x < len(matrix); x++ {
		for y := 0; y < len(matrix[0]); y++ {
			item := matrix[x][y]
			rp.positions[item.GetGuid()] = PositionToArray(x, y)
			rp.searchNeighbour(x, y, synchronization.NewMap[int64, bool](), synchronization.NewMap[int64, bool]())
		}
	}
	return rp
}

// RadiationPattern 辐射图数据结构
//   - 辐射图用于将一个二维数组里相邻的所有类型相同的成员进行标注
type RadiationPattern[ItemType comparable, Item RadiationPatternItem[ItemType]] struct {
	matrix    [][]Item
	links     *synchronization.Map[int64, map[int64]bool] // 成员类型相同且相连的链接
	positions map[int64][2]int                            // 根据成员guid记录的成员位置
}

// GetLinks 获取特定成员能够辐射到的所有成员
func (slf *RadiationPattern[ItemType, Item]) GetLinks(guid int64) []int64 {
	return hash.KeyToSlice(slf.links.Get(guid))
}

// GetLinkPositions 获取特定成员能够辐射到的所有成员位置
func (slf *RadiationPattern[ItemType, Item]) GetLinkPositions(guid int64) [][2]int {
	links := slf.links.Get(guid)
	var result = make([][2]int, 0, len(links))
	for g := range links {
		result = append(result, slf.positions[g])
	}
	return result
}

// GetPosition 获取特定成员的位置
func (slf *RadiationPattern[ItemType, Item]) GetPosition(guid int64) [2]int {
	return slf.positions[guid]
}

// Refresh 刷新特定位置成员并且更新其辐射信息
func (slf *RadiationPattern[ItemType, Item]) Refresh(x, y int, item Item) {
	old := slf.matrix[x][y]
	oldGuid := old.GetGuid()
	for linkGuid := range slf.links.Get(oldGuid) {
		xy := slf.positions[linkGuid]
		slf.searchNeighbour(xy[0], xy[1], synchronization.NewMap[int64, bool](), synchronization.NewMap[int64, bool]())
	}
	slf.links.Delete(oldGuid)
	delete(slf.positions, oldGuid)

	slf.matrix[x][y] = item
	slf.positions[item.GetGuid()] = PositionToArray(x, y)
	slf.searchNeighbour(x, y, synchronization.NewMap[int64, bool](), synchronization.NewMap[int64, bool]())
}

// RefreshBySwap 通过交换的方式刷新两个成员的辐射信息
func (slf *RadiationPattern[ItemType, Item]) RefreshBySwap(x1, y1, x2, y2 int, item1, item2 Item) {
	var xys = [][2]int{PositionToArray(x1, y1), PositionToArray(x2, y2)}
	for _, xy := range xys {
		x, y := PositionArrayToXY(xy)
		old := slf.matrix[x][y]
		oldGuid := old.GetGuid()
		for linkGuid := range slf.links.Get(oldGuid) {
			xy := slf.positions[linkGuid]
			slf.searchNeighbour(xy[0], xy[1], synchronization.NewMap[int64, bool](), synchronization.NewMap[int64, bool]())
		}
		slf.links.Delete(oldGuid)
		delete(slf.positions, oldGuid)
	}
	for i, item := range []Item{item1, item2} {
		x, y := PositionArrayToXY(xys[i])
		slf.matrix[x][y] = item
		slf.positions[item.GetGuid()] = PositionToArray(x, y)
		slf.searchNeighbour(x, y, synchronization.NewMap[int64, bool](), synchronization.NewMap[int64, bool]())
	}
}

func (slf *RadiationPattern[ItemType, Item]) searchNeighbour(x, y int, filter *synchronization.Map[int64, bool], childrenLinks *synchronization.Map[int64, bool]) {
	var (
		item           = slf.matrix[x][y]
		neighboursLock sync.Mutex
		neighbours     = map[int64]bool{}
		itemType       = item.GetType()
		wait           sync.WaitGroup
		itemGuid       = item.GetGuid()
		handle         = func(x, y int) bool {
			neighbour := slf.matrix[x][y]
			if neighbour.GetType() != itemType {
				return false
			}
			neighbourGuid := neighbour.GetGuid()
			neighboursLock.Lock()
			neighbours[neighbourGuid] = true
			neighboursLock.Unlock()
			childrenLinks.Set(neighbourGuid, true)
			slf.searchNeighbour(x, y, filter, childrenLinks)
			return true
		}
	)
	if filter.Get(itemGuid) {
		return
	}
	filter.Set(itemGuid, true)
	wait.Add(4)
	go func() {
		for sy := y - 1; sy >= 0; sy-- {
			if !handle(x, sy) {
				break
			}
		}
		wait.Done()
	}()
	go func() {
		for sy := y + 1; sy < len(slf.matrix[0]); sy++ {
			if !handle(x, sy) {
				break
			}
		}
		wait.Done()
	}()
	go func() {
		for sx := x - 1; sx >= 0; sx-- {
			if !handle(sx, y) {
				break
			}
		}
		wait.Done()
	}()
	go func() {
		for sx := x + 1; sx < len(slf.matrix); sx++ {
			if !handle(sx, y) {
				break
			}
		}
		wait.Done()
	}()
	wait.Wait()
	childrenLinks.Range(func(key int64, value bool) {
		neighbours[key] = value
	})
	slf.links.Set(itemGuid, neighbours)
}
