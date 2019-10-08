package webserver

// Config ...
type Config struct {
	BindAddr                 string `toml:"bind_addr"`
	LogLevel                 string `toml:"log_level"`
	DatabaseConnectionString string `toml:"dsn_url"`
	URLShema                 string `toml:"shema"`
	Hostname                 string `toml:"hostname"`
	SwaggerPath              string `toml:"docs_url"`
	SwaggerFile              string `toml:"docs_source"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8888",
		LogLevel: "debug",

		URLShema:    "http",
		Hostname:    "localhost",
		SwaggerPath: "/api/v1/docs/swagger.json",
		SwaggerFile: "docs/swagger.json",
	}
}
