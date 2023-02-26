package database

type MongoConfig struct {
	URI      string `yaml:"uri" mapstructure:"uri"`
	Database string `yaml:"database" mapstructure:"database"`
}
