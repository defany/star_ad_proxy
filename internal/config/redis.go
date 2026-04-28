package config

func DragonflyAddr() string {
	return cfg.String("dragonfly.addr")
}

func DragonflyPassword() string {
	return cfg.String("dragonfly.password")
}

func DragonflyDB() int {
	return cfg.Int("dragonfly.db")
}
