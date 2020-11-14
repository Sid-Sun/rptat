package internal

// Config is the father of stupids
type Config struct {
	App         AppConfig     `toml:"App"`
	ProxyConfig []ProxyConfig `toml:"Proxies"`
}

// APIConfig is stupid
type AppConfig struct {
	Env  string `toml:"env"`
	Host string `toml:"host"`
	Port uint   `toml:"port"`
}

// MetricsConfig is also stupid
type MetricsConfig struct {
	MinForSync           uint `toml:"max_pending"`
	PeriodicSyncInterval uint `toml:"periodic_sync_interval"`
}

// StoreConfig is also stupid
type StoreConfig struct {
	FileName  string `toml:"file_name"`
	FilePerms int    `toml:"file_perms"`
}

// ProxyConfig is very
type ProxyConfig struct {
	Protocol      string        `toml:"protocol"`
	Port          int           `toml:"port"`
	Host          string        `toml:"host"`
	Hostname      string        `toml:"hostname"`
	StoreConfig   StoreConfig   `toml:"Store"`
	MetricsConfig MetricsConfig `toml:"Metrics"`
}
