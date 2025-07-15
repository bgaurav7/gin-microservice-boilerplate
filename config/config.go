package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/bgaurav7/gin-microservice-boilerplate/internal/infrastructure/logger"
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Logger   logger.Config  `mapstructure:"logger"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	RBAC     RBACConfig     `mapstructure:"rbac"`
}

// AppConfig represents the application configuration
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// AuthConfig represents the authentication configuration
type AuthConfig struct {
	JWTSecret      string `mapstructure:"jwt_secret"`
	JWTExpiryHours int    `mapstructure:"jwt_expiry_hours"`
	SuperAdminEmail string `mapstructure:"superadmin_email"`
}

// RBACConfig represents the RBAC configuration
type RBACConfig struct {
	ModelPath  string `mapstructure:"model_path"`
	PolicyPath string `mapstructure:"policy_path"`
}

// Load loads the configuration from the config file and environment variables
func Load() (*Config, error) {
	// Determine which config file to load based on environment
	env := os.Getenv("APP_ENVIRONMENT")
	if env == "" || (env != "dev" && env != "prod") {
		env = "dev" // Default to development environment
	}

	// First load the common config
	baseConfig := viper.New()
	baseConfig.SetConfigName("config")
	baseConfig.SetConfigType("yaml")
	baseConfig.AddConfigPath("./config")

	// Read common configuration file
	if err := baseConfig.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read common config file: %w", err)
	}

	// Now load the environment-specific config
	envConfig := viper.New()
	envConfig.SetConfigName(env)
	envConfig.SetConfigType("yaml")
	envConfig.AddConfigPath("./config")

	// Read environment-specific configuration file
	if err := envConfig.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read %s config file: %w", env, err)
	}

	// Merge the configurations
	for _, k := range envConfig.AllKeys() {
		baseConfig.Set(k, envConfig.Get(k))
	}

	// Set up environment variable handling
	baseConfig.SetEnvPrefix("")
	baseConfig.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	baseConfig.AutomaticEnv()

	// Bind environment variables to override config file values
	baseConfig.BindEnv("app.name", "APP_NAME")
	baseConfig.BindEnv("app.environment", "APP_ENVIRONMENT")
	baseConfig.BindEnv("server.host", "SERVER_HOST")
	baseConfig.BindEnv("server.port", "SERVER_PORT")
	baseConfig.BindEnv("server.read_timeout", "SERVER_READ_TIMEOUT")
	baseConfig.BindEnv("server.write_timeout", "SERVER_WRITE_TIMEOUT")
	baseConfig.BindEnv("logger.level", "LOGGER_LEVEL")
	baseConfig.BindEnv("database.driver", "DB_DRIVER")
	baseConfig.BindEnv("database.host", "DB_HOST")
	baseConfig.BindEnv("database.port", "DB_PORT")
	baseConfig.BindEnv("database.username", "DB_USERNAME")
	baseConfig.BindEnv("database.password", "DB_PASSWORD")
	baseConfig.BindEnv("database.name", "DB_NAME")
	baseConfig.BindEnv("database.sslmode", "DB_SSLMODE")
	baseConfig.BindEnv("database.max_idle_conns", "DB_MAX_IDLE_CONNS")
	baseConfig.BindEnv("database.max_open_conns", "DB_MAX_OPEN_CONNS")
	baseConfig.BindEnv("database.conn_max_lifetime", "DB_CONN_MAX_LIFETIME")
	baseConfig.BindEnv("auth.jwt_secret", "JWT_SECRET")
	baseConfig.BindEnv("auth.jwt_expiry_hours", "JWT_EXPIRY_HOURS")
	baseConfig.BindEnv("auth.superadmin_email", "SUPERADMIN_EMAIL")

	// Unmarshal configuration
	var config Config
	if err := baseConfig.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// DSN returns the database connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		c.Driver,
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}
