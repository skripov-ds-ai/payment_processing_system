package relational

import "fmt"

// SQLConnectConfig determines config to connect to SQL database
type SQLConnectConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

// NewSQLConnectConfig creates new pg config instance
func NewSQLConnectConfig(username, password, host, port, database string) *SQLConnectConfig {
	return &SQLConnectConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}
}

// DSN creates dsn string from config
func (cfg SQLConnectConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password,
		cfg.Host, cfg.Port, cfg.Database,
	)
}
