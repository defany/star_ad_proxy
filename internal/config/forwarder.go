package config

import "time"

func ForwarderTimeout() time.Duration {
	return cfg.Duration("forwarder.timeout")
}
