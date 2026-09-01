package pkgUtility

import "github.com/cornelk/hashmap"

// Hashable 是 cornelk/hashmap 接受的 key 型別集合（數值與字串）。
// 它自己的 key constraint 沒有匯出，這裡列一份等價的給 BiMultiMap 當型別參數約束用。
type Hashable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

type BiMultiMap[L, R Hashable] struct {
	leftToRights *hashmap.Map[L, *hashmap.Map[R, struct{}]]
	rightToLefts *hashmap.Map[R, *hashmap.Map[L, struct{}]]
}

func NewBiMultiMap[L, R Hashable]() *BiMultiMap[L, R] {
	return &BiMultiMap[L, R]{
		leftToRights: hashmap.New[L, *hashmap.Map[R, struct{}]](),
		rightToLefts: hashmap.New[R, *hashmap.Map[L, struct{}]](),
	}
}

// innerOf 取出 key 對應的內層 map，沒有就原子地建一張再放進去。
// 先 Get 一次是為了讓「key 已存在」的常見路徑不用白配一張新 map。
func innerOf[K Hashable, IK Hashable](oOuter *hashmap.Map[K, *hashmap.Map[IK, struct{}]], key K) *hashmap.Map[IK, struct{}] {
	oInner, bGotten := oOuter.Get(key)
	if bGotten {
		return oInner
	}

	oInner, _ = oOuter.GetOrInsert(key, hashmap.New[IK, struct{}]())
	return oInner
}

// Insert 建立 l/r 的雙向關聯，重複呼叫同一組 l/r 不會有副作用。
func (oSelf *BiMultiMap[L, R]) Insert(l L, r R) {
	innerOf(oSelf.leftToRights, l).Set(r, struct{}{})
	innerOf(oSelf.rightToLefts, r).Set(l, struct{}{})
}

// Left 回傳這個 left 目前關聯到的所有 right。
func (oSelf *BiMultiMap[L, R]) Left(l L) []R {
	oRights, bGotten := oSelf.leftToRights.Get(l)
	if !bGotten {
		return nil
	}

	aResult := make([]R, 0, oRights.Len())
	oRights.Range(func(r R, _ struct{}) bool {
		aResult = append(aResult, r)
		return true
	})

	return aResult
}

// Right 回傳這個 right 目前關聯到的所有 left。
func (oSelf *BiMultiMap[L, R]) Right(r R) []L {
	oLefts, bGotten := oSelf.rightToLefts.Get(r)
	if !bGotten {
		return nil
	}

	aResult := make([]L, 0, oLefts.Len())
	oLefts.Range(func(l L, _ struct{}) bool {
		aResult = append(aResult, l)
		return true
	})

	return aResult
}

// Remove 拆掉單一一組 l/r 的關聯，其他跟 l 或 r 相關的關聯不受影響。
func (oSelf *BiMultiMap[L, R]) Remove(l L, r R) {
	if oRights, bGotten := oSelf.leftToRights.Get(l); bGotten {
		oRights.Del(r)
		if oRights.Len() == 0 {
			oSelf.leftToRights.Del(l)
		}
	}

	if oLefts, bGotten := oSelf.rightToLefts.Get(r); bGotten {
		oLefts.Del(l)
		if oLefts.Len() == 0 {
			oSelf.rightToLefts.Del(r)
		}
	}
}

// RemoveLeft 拆掉這個 left 的全部關聯，連帶清掉對面每個 right 記著這個 left 的紀錄。
func (oSelf *BiMultiMap[L, R]) RemoveLeft(l L) {
	oRights, bGotten := oSelf.leftToRights.Get(l)
	if !bGotten {
		return
	}

	oRights.Range(func(r R, _ struct{}) bool {
		if oLefts, bLefts := oSelf.rightToLefts.Get(r); bLefts {
			oLefts.Del(l)
			if oLefts.Len() == 0 {
				oSelf.rightToLefts.Del(r)
			}
		}
		return true
	})

	oSelf.leftToRights.Del(l)
}

// RemoveRight 拆掉這個 right 的全部關聯，連帶清掉對面每個 left 記著這個 right 的紀錄。
func (oSelf *BiMultiMap[L, R]) RemoveRight(r R) {
	oLefts, bGotten := oSelf.rightToLefts.Get(r)
	if !bGotten {
		return
	}

	oLefts.Range(func(l L, _ struct{}) bool {
		if oRights, bRights := oSelf.leftToRights.Get(l); bRights {
			oRights.Del(r)
			if oRights.Len() == 0 {
				oSelf.leftToRights.Del(l)
			}
		}
		return true
	})

	oSelf.rightToLefts.Del(r)
}
