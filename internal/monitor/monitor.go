package monitor

import "fmt"

func InitMonitor() {
	InitGlobalAPH()
	GlobalAPH.Init()

	weights, err := GlobalAPH.CalculateWeights()
	if err != nil {
		fmt.Printf("Error calculating weights: %v\n", err)
	}

	InitTriggerMechanism(weights, []float64{0.8, 0.8, 0.8, 150, 50}, 0.8, 3)

	InitZScoreCalculator(50)
}
