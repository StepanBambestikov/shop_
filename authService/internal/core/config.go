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
	Endpoint      string `config:"endpoint" validate:"required,default=/metrics" yaml:"endpoint"`
	ExportDefault bool   `config:"export_default" validate:"required,default=true" yaml:"export_default"`
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

type IntegrationsConfig struct {
	Keycloak KeycloakConfig `config:"keycloak" validate:"required" yaml:"keycloak"`
}

type Config struct {
	Server       ServerConfig       `config:"server,required" yaml:"server"`
	Health       HealthConfig       `config:"health,required" yaml:"health"`
	Metrics      MetricsConfig      `config:"metrics,required" yaml:"metrics"`
	Swagger      SwaggerConfig      `config:"swagger,required" yaml:"swagger"`
	Integrations IntegrationsConfig `config:"integrations,required" yaml:"integrations"`
}
