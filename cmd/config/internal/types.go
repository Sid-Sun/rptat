package internal

// Config is the father of stupids
type Config struct {
	API APIConfig `toml:"API"`
	App struct {
		Env string `toml:"env"`
	} `toml:"App"`
	ProxyConfig   ProxyConfig   `toml:"Proxy"`
	MetricsConfig MetricsConfig `toml:"Metrics"`
}

// APIConfig is stupid
type APIConfig struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
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

type (
	// ProxyConfig is very
	ProxyConfig struct {
		Listen ListenCfg `toml:"listen"`
		Serve  ServeCfg  `toml:"serve"`
		Store StoreConfig `toml:"store"`
	}
	ListenCfg struct {
		Port int    `toml:"port"`
		Host string `toml:"host"`
	}
	ServeCfg struct {
		Protocol string `toml:"protocol"`
		Port     int    `toml:"port"`
		Host     string `toml:"host"`
	}
)
