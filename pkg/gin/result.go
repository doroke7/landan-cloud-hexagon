package pkgGin

// Result 同時容納「單筆」與「多筆」兩種回應形狀，靠 omitempty 只序列化有設值的那個。
// any 欄位的 omitempty 判斷的是「interface 為 nil」，沒設的那個維持 nil interface 就會被省略。
type Result struct {
	One  any `json:"one,omitempty"`
	Ones any `json:"ones,omitempty"`
	Tree any `json:"tree,omitempty"`
}

// NewResult 三個參數對應 One / Ones / Tree，用不到的傳 nil。
func NewResult(oOne any, aOnes any, aTree any) *Result {
	return &Result{
		One:  oOne,
		Ones: aOnes,
		Tree: aTree,
	}
}
