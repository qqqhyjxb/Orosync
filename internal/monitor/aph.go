package monitor

/*
	0：电量
	1: cpu利用率
	2: 内存利用率
	3：网络延迟
	4: 偏离距离
*/

type APH struct {
	JudgmentMatrix [][]float64 //判断矩阵
}

var GlobalAPH *APH

func InitGlobalAPH() {
	GlobalAPH = &APH{}
}

func (a *APH) Init() {
	a.JudgmentMatrix = [][]float64{
		[]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		[]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		[]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		[]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
		[]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
	}
}
