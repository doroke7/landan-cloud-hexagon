package pkg

import "sync"

// BiMultiMap 是一個雙向、多對多的關聯容器：一個 left 可以對應多個 right，
// 一個 right 也可以對應多個 left，兩個方向的資料互為鏡像、保證一致。
// 所有操作共用同一把鎖，呼叫端不用自己處理並發，也不用自己記得「兩個
// 方向要一起改」——這是取代手動維護兩個各自獨立的 map、外面再包一層
// mutex 的寫法。
type BiMultiMap[L, R comparable] struct {
	mu           sync.RWMutex
	leftToRights map[L]map[R]struct{}
	rightToLefts map[R]map[L]struct{}
}

func NewBiMultiMap[L, R comparable]() *BiMultiMap[L, R] {
	return &BiMultiMap[L, R]{
		leftToRights: make(map[L]map[R]struct{}),
		rightToLefts: make(map[R]map[L]struct{}),
	}
}

// Insert 建立 l/r 的雙向關聯，重複呼叫同一組 l/r 不會有副作用。
func (oSelf *BiMultiMap[L, R]) Insert(l L, r R) {
	oSelf.mu.Lock()
	defer oSelf.mu.Unlock()

	if oSelf.leftToRights[l] == nil {
		oSelf.leftToRights[l] = make(map[R]struct{})
	}
	oSelf.leftToRights[l][r] = struct{}{}

	if oSelf.rightToLefts[r] == nil {
		oSelf.rightToLefts[r] = make(map[L]struct{})
	}
	oSelf.rightToLefts[r][l] = struct{}{}
}

// Left 回傳這個 left 目前關聯到的所有 right。
func (oSelf *BiMultiMap[L, R]) Left(l L) []R {
	oSelf.mu.RLock()
	defer oSelf.mu.RUnlock()

	aResult := make([]R, 0, len(oSelf.leftToRights[l]))
	for r := range oSelf.leftToRights[l] {
		aResult = append(aResult, r)
	}

	return aResult
}

// Right 回傳這個 right 目前關聯到的所有 left。
func (oSelf *BiMultiMap[L, R]) Right(r R) []L {
	oSelf.mu.RLock()
	defer oSelf.mu.RUnlock()

	aResult := make([]L, 0, len(oSelf.rightToLefts[r]))
	for l := range oSelf.rightToLefts[r] {
		aResult = append(aResult, l)
	}

	return aResult
}

// Remove 拆掉單一一組 l/r 的關聯，其他跟 l 或 r 相關的關聯不受影響。
func (oSelf *BiMultiMap[L, R]) Remove(l L, r R) {
	oSelf.mu.Lock()
	defer oSelf.mu.Unlock()

	delete(oSelf.leftToRights[l], r)
	if len(oSelf.leftToRights[l]) == 0 {
		delete(oSelf.leftToRights, l)
	}

	delete(oSelf.rightToLefts[r], l)
	if len(oSelf.rightToLefts[r]) == 0 {
		delete(oSelf.rightToLefts, r)
	}
}

// RemoveLeft 拆掉這個 left 的全部關聯，連帶清掉對面每個 right 記著這個 left 的紀錄。
func (oSelf *BiMultiMap[L, R]) RemoveLeft(l L) {
	oSelf.mu.Lock()
	defer oSelf.mu.Unlock()

	for r := range oSelf.leftToRights[l] {
		delete(oSelf.rightToLefts[r], l)
		if len(oSelf.rightToLefts[r]) == 0 {
			delete(oSelf.rightToLefts, r)
		}
	}

	delete(oSelf.leftToRights, l)
}

// RemoveRight 拆掉這個 right 的全部關聯，連帶清掉對面每個 left 記著這個 right 的紀錄。
func (oSelf *BiMultiMap[L, R]) RemoveRight(r R) {
	oSelf.mu.Lock()
	defer oSelf.mu.Unlock()

	for l := range oSelf.rightToLefts[r] {
		delete(oSelf.leftToRights[l], r)
		if len(oSelf.leftToRights[l]) == 0 {
			delete(oSelf.leftToRights, l)
		}
	}

	delete(oSelf.rightToLefts, r)
}
