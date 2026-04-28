package config

func AdminToken() string {
	return cfg.String("admin.token")
}
