package core

type ServerConfig struct {
	Host string `config:"host" validate:"required,default=0.0.0.0" yaml:"host"`
	Port uint16 `config:"port" validate:"required,default=8080" yaml:"port"`
}

type SwaggerConfig struct {
	Enabled  bool   `config:"enabled" validate:"required,default=true" yaml:"enabled"`
	Endpoint string `config:"endpoint" validate:"required,default=/swagger" yaml:"endpoint"`
}

type HealthConfig struct {
	Enabled  bool   `config:"enabled" validate:"required,default=true" yaml:"enabled"`
	Endpoint string `config:"endpoint" validate:"required,default=/healthz" yaml:"endpoint"`
}

type MetricsConfig struct {
	Enabled       bool   `config:"enabled" validate:"required,default=true" yaml:"enabled"`
	Endpoint      string `config:"endpoint" validate:"required,default=/healthz" yaml:"endpoint"`
	ExportDefault bool   `config:"export_default" validate:"required,default=true" yaml:"export_default"`
}

type RabbitConfig struct {
	DSN        string `config:"dsn" validate:"required" yaml:"dsn"`
	Exchange   string `config:"exchange" validate:"required" yaml:"exchange"`
	MaxRetries uint64 `config:"max_retries" validate:"required" yaml:"max_retries"`
}

type PostgresConfig struct {
	Host     string `config:"host" validate:"required" yaml:"host"`
	Port     string `config:"port" validate:"required" yaml:"port"`
	Password string `config:"password" validate:"omitempty,default=''" yaml:"password"`
	User     string `config:"user" validate:"required" yaml:"user"`
	Dbname   string `config:"dbname" validate:"required" yaml:"dbname"`
}

type IntegrationsConfig struct {
	Rabbit   RabbitConfig   `config:"rabbitmq" validate:"required" yaml:"rabbitmq"`
	Postgres PostgresConfig `config:"postgres" validate:"required" yaml:"postgres"`
}

type Config struct {
	Server       ServerConfig       `config:"server,required" yaml:"server"`
	Health       HealthConfig       `config:"health,required" yaml:"health"`
	Metrics      MetricsConfig      `config:"metrics,required" yaml:"metrics"`
	Swagger      SwaggerConfig      `config:"swagger,required" yaml:"swagger"`
	Integrations IntegrationsConfig `config:"integrations,required" yaml:"integrations"`
}

type MultipleConfig struct {
	GateConfig Config `config:"gateapp,required" yaml:"gateapp"`
	CoreConfig Config `config:"coreapp,required" yaml:"coreapp"`
}

type KeycloakConfig struct {
	Uri                  string `config:"uri,required" yaml:"uri"`
	TokenRefreshInterval uint64 `config:"token_refresh_interval,required" yaml:"token_refresh_interval"`
	Realm                string `config:"realm,required" yaml:"realm"`

	Client struct {
		ID     string `config:"id,required" yaml:"id"`
		Secret string `config:"secret,required" yaml:"secret"`
	} `config:"client,required" yaml:"client"`
	Admin struct {
		Username string `config:"username,required" yaml:"username"`
		Password string `config:"password,required" yaml:"password"`
		Realm    string `config:"realm,required" yaml:"realm"`
	} `config:"admin,required" yaml:"admin"`
}
