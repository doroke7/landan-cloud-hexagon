package pkg

func NewResultOnes(aRows any) *ResultOnes {
	return &ResultOnes{
		Ones: aRows,
	}
}

func NewResultOne(oRow any) *ResultOne {
	return &ResultOne{
		One: oRow,
	}
}

type ResultOnes struct {
	Ones any `json:"ones"`
}

type ResultOne struct {
	One any `json:"one"`
}
