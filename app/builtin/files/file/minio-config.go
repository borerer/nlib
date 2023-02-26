package file

type MinioConfig struct {
	Endpoint  string `yaml:"endpoint" mapstructure:"endpoint"`
	AccessKey string `yaml:"access-key" mapstructure:"access-key"`
	SecretKey string `yaml:"secret-key" mapstructure:"secret-key"`
	UseSSL    bool   `yaml:"use-ssl" mapstructure:"use-ssl"`
	Bucket    string `yaml:"bucket" mapstructure:"bucket"`
}
