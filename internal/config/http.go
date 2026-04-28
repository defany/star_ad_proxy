package config

func HttpAddr() string {
	return cfg.String("http.addr")
}
