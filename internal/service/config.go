package service

type Config struct {
	DataFolder string
	BaseFolder string
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
