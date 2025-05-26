package service

type Config struct {
	TracePort       string
	DebuggerPort    string
	DebugChartPort  string
	DataFolder      string
	BaseFolder      string
	K8sSharedFolder string
	Editor          string
	Locale          string
}

func (c *Config) CombinationFolder() string {
	return c.DataFolder + "/combination"
}

func (c *Config) DefinitionFolder() string {
	return c.DataFolder + "/definition"
}

func (c *Config) PlanFolder() string {
	return c.DataFolder + "/plan"
}

func (c *Config) StorageFolder() string {
	return c.DataFolder + "/storage"
}
