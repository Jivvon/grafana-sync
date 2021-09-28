package grafana

type Config struct {
	Host	string `mapstructure:"host"`
	Port    string `mapstructure:"port"`
	AdminID string `mapstructure:"adminID"`
	AdminPW string `mapstructure:"adminPW"`
}

func (c *Config) Address() string {
	return c.Host + ":" + c.Port
}

func (c *Config) Auth() string {
	return c.AdminID + ":" + c.AdminPW
}

