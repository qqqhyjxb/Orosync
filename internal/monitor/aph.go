package monitor

import (
	"fmt"
	"math"
)

// APH 结构体包含判断矩阵
type APH struct {
	JudgmentMatrix [][]float64 // 判断矩阵
}

// GlobalAPH 全局变量，用于存储全局的 AHP 实例
var GlobalAPH *APH

// InitGlobalAPH 初始化全局 AHP 实例
func InitGlobalAPH() {
	GlobalAPH = &APH{}
}

// Init 初始化判断矩阵
func (a *APH) Init() {
	a.JudgmentMatrix = [][]float64{
		//         0:电量  1:CPU利用率  2:内存利用率 3:网络延迟 4:偏离距离
		{1, 4, 4, 4, 5},              // 电量
		{0.25, 1, 2, 4, 5},           // CPU
		{0.25, 0.5, 1, 2, 4},         // 内存
		{0.25, 0.25, 0.5, 1, 4},      // 网络延迟
		{0.125, 0.25, 0.25, 0.25, 1}, // 偏离距离
	}
}

// 计算矩阵每列的和
func columnSums(matrix [][]float64) []float64 {
	numCols := len(matrix[0])
	columnSums := make([]float64, numCols)

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < numCols; j++ {
			columnSums[j] += matrix[i][j]
		}
	}

	return columnSums
}

// 归一化矩阵
func normalizeMatrix(matrix [][]float64, columnSums []float64) [][]float64 {
	numRows := len(matrix)
	numCols := len(matrix[0])
	normalizedMatrix := make([][]float64, numRows)

	for i := 0; i < numRows; i++ {
		normalizedMatrix[i] = make([]float64, numCols)
		for j := 0; j < numCols; j++ {
			normalizedMatrix[i][j] = matrix[i][j] / columnSums[j]
		}
	}

	return normalizedMatrix
}

// 计算矩阵每行的平均值，即权重
func calculateWeights(normalizedMatrix [][]float64) []float64 {
	numRows := len(normalizedMatrix)
	weights := make([]float64, numRows)

	for i := 0; i < numRows; i++ {
		rowSum := 0.0
		for j := 0; j < len(normalizedMatrix[i]); j++ {
			rowSum += normalizedMatrix[i][j]
		}
		weights[i] = rowSum / float64(len(normalizedMatrix[i]))
	}

	return weights
}

// calculateLambdaMax 计算判断矩阵的最大特征值
func calculateLambdaMax(matrix [][]float64, weights []float64) float64 {
	n := len(matrix)
	sum := 0.0

	for i := 0; i < n; i++ {
		rowSum := 0.0
		for j := 0; j < n; j++ {
			rowSum += matrix[i][j] * weights[j]
		}
		sum += rowSum / weights[i]
	}

	return sum / float64(n)
}

// getRI 获取随机一致性指标
func getRI(n int) float64 {
	riTable := map[int]float64{
		1:  0.0,
		2:  0.0,
		3:  0.58,
		4:  0.90,
		5:  1.12,
		6:  1.24,
		7:  1.32,
		8:  1.41,
		9:  1.45,
		10: 1.49,
	}
	return riTable[n]
}

// checkConsistency 执行一致性检验
func (a *APH) checkConsistency(weights []float64) (cr float64, valid bool) {
	n := len(a.JudgmentMatrix)
	if n <= 1 {
		return 0.0, true
	}

	lambdaMax := calculateLambdaMax(a.JudgmentMatrix, weights)
	ci := (lambdaMax - float64(n)) / float64(n-1)
	ri := getRI(n)
	cr = ci / ri

	// 保留4位小数
	cr = math.Round(cr*10000) / 10000
	return cr, cr < 0.1
}

// CalculateWeights 计算 AHP 权重（添加一致性检验）
func (a *APH) CalculateWeights() ([]float64, error) {
	c := columnSums(a.JudgmentMatrix)
	normalizedMatrix := normalizeMatrix(a.JudgmentMatrix, c)
	weights := calculateWeights(normalizedMatrix)

	if cr, valid := a.checkConsistency(weights); !valid {
		return weights, fmt.Errorf("consistency check failed, CR=%.4f ≥ 0.1, please adjust judgment matrix", cr)
	}

	return weights, nil
}

// UpdateJudgmentMatrix 更新判断矩阵
func (a *APH) UpdateJudgmentMatrix(matrix [][]float64) {
	a.JudgmentMatrix = matrix
	fmt.Println("\nJudgment Matrix updated successfully.")
}

// GetLatestWeights 获取最新的权重
func (a *APH) GetLatestWeights() ([]float64, error) {
	return a.CalculateWeights()
}

// Test 测试方法
func (a *APH) Test() {
	fmt.Println("=== Initial Weights Test ===")
	if weights, err := a.CalculateWeights(); err != nil {
		fmt.Println("Initial Weights Calculation Results:")
		printWeights(weights)
		fmt.Println("❌ Error:", err)
	} else {
		fmt.Println("Initial Weights Calculation Results (Consistency Passed):")
		printWeights(weights)
	}

	// 更新测试矩阵
	newMatrix := [][]float64{
		{1, 1.0 / 3, 1.0 / 2, 2, 1},
		{3, 1, 2, 4, 3},
		{2, 1.0 / 2, 1, 3, 2},
		{1.0 / 2, 1.0 / 4, 1.0 / 3, 1, 1.0 / 2},
		{1, 1.0 / 3, 1.0 / 2, 2, 1},
	}
	a.UpdateJudgmentMatrix(newMatrix)

	fmt.Println("\n=== Updated Weights Test ===")
	if weights, err := a.GetLatestWeights(); err != nil {
		fmt.Println("Updated Weights Calculation Results:")
		printWeights(weights)
		fmt.Println("❌ Error:", err)
	} else {
		fmt.Println("Updated Weights Calculation Results (Consistency Passed):")
		printWeights(weights)
	}
}

// printWeights 辅助打印函数
func printWeights(weights []float64) {
	for i, weight := range weights {
		fmt.Printf("Indicator %d Weight: %.4f\n", i+1, weight)
	}
}
