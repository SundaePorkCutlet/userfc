package config

type Config struct {
	App      AppConfig      `yaml:"app" validate:"required"`
	Database DatabaseConfig `yaml:"database" validate:"required"`
	Redis    RedisConfig    `yaml:"redis" validate:"required"`
	Secret   SecretConfig   `yaml:"secret"`
	Vault    VaultConfig    `yaml:"vault"`
	Tracing  TracingConfig  `yaml:"tracing"`
}

type TracingConfig struct {
	Endpoint    string `yaml:"endpoint" mapstructure:"endpoint"`
	ServiceName string `yaml:"service_name" mapstructure:"service_name"`
	Enabled     bool   `yaml:"enabled" mapstructure:"enabled"`
}

type VaultConfig struct {
	Host  string `yaml:"host" validate:"required"`
	Token string `yaml:"token" validate:"required"`
	Path  string `yaml:"path" validate:"required"`
}

type SecretVaultConfig struct {
	DatabaseSecret DatabaseSecretConfig `json:"database"`
	RedisSecret    RedisSecretConfig    `json:"redis"`
	JWTSecret      string               `json:"jwt_secret"`
}

type DatabaseSecretConfig struct {
	Password string `json:"password"`
}

type RedisSecretConfig struct {
	Password string `json:"password"`
}

type AppConfig struct {
	Port     string `yaml:"port" validate:"required"`
	GRPCPort string `yaml:"grpc_port" mapstructure:"grpc_port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Name     string `yaml:"name" validate:"required"`
}

type RedisConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     string `yaml:"port" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

type SecretConfig struct {
	JwtSecret string `yaml:"jwt_secret" validate:"required"`
}
