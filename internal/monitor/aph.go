package monitor

import "fmt"

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
	// 根据具体的业务需求完善判断矩阵
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
		weights[i] = rowSum / float64(len(normalizedMatrix[i])) // 计算每行的平均值
	}

	return weights
}

// CalculateWeights 计算 AHP 权重
func (a *APH) CalculateWeights() []float64 {
	// 1. 计算每列的和
	c := columnSums(a.JudgmentMatrix)
	fmt.Println("Column Sums:", c)

	// 2. 归一化矩阵
	normalizedMatrix := normalizeMatrix(a.JudgmentMatrix, c)
	fmt.Println("\nNormalized Matrix:")
	for _, row := range normalizedMatrix {
		fmt.Println(row)
	}

	// 3. 计算权重
	weights := calculateWeights(normalizedMatrix)
	return weights
}

// UpdateJudgmentMatrix 更新判断矩阵
func (a *APH) UpdateJudgmentMatrix(matrix [][]float64) {
	a.JudgmentMatrix = matrix
	fmt.Println("\nJudgment Matrix updated successfully.")
}

// GetLatestWeights 获取最新的权重
func (a *APH) GetLatestWeights() []float64 {
	// 获取最新的权重值
	weights := a.CalculateWeights()
	return weights
}

// Test 测试ahp
func (a *APH) Test() {
	// 计算权重
	weights := GlobalAPH.CalculateWeights()

	// 输出计算结果
	for i, weight := range weights {
		fmt.Printf("Indicator %d Weight: %f\n", i+1, weight)
	}

	// 更新判断矩阵
	newMatrix := [][]float64{
		// 更新后的矩阵
		{1, 5, 4, 3, 6},
		{0.2, 1, 2.5, 3.5, 4},
		{0.25, 0.4, 1, 2, 3},
		{0.3, 0.3, 0.5, 1, 3},
		{0.125, 0.2, 0.3, 0.4, 1},
	}
	GlobalAPH.UpdateJudgmentMatrix(newMatrix)

	// 获取更新后的最新权重
	latestWeights := GlobalAPH.GetLatestWeights()

	// 输出更新后的权重
	fmt.Println("\nUpdated Weights:")
	for i, weight := range latestWeights {
		fmt.Printf("Updated Indicator %d Weight: %f\n", i+1, weight)
	}
}
