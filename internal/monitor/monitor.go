package monitor

func InitMonitor() {
	InitGlobalAPH()
	GlobalAPH.Init()
	InitZScoreCalculator(50)
	//InitTriggerMechanism()
}
