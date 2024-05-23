package initialise

// InitResource 初始化资源
func InitResource() error {
	err := InitInfra()
	if err != nil {
		return err
	}
	return nil
}

// InitInfra 初始化基础组件
func InitInfra() error {
	return nil
}
