package config

import (
    "fmt"

    "github.com/caarlos0/env/v11"
    "github.com/joho/godotenv"
)

type ServiceConfig struct {
    Env        string `env:"ENV" envDefault:"dev"`
    SrvPort    string `env:"SERVER_PORT,required" envDefault:"8070"`
    PgHost     string `env:"POSTGRES_HOST,required" envDefault:"0.0.0.0"`
    PgPort     uint   `env:"POSTGRES_PORT,required" envDefault:"5433"`
    PgUser     string `env:"POSTGRES_USER,required" envDefault:"segmentServiceUser"`
    PgPass     string `env:"POSTGRES_PASSWORD,required" envDefault:"verySecretPassword"`
    PgDB       string `env:"POSTGRES_DB" envDefault:"segments_service"`
    PgTestPort uint   `env:"POSTGRES_TEST_PORT" envDefault:"5435"`
}

// Load environment variables for the service configuration.
func Load() (*ServiceConfig, error) {
    _ = godotenv.Load()
    cfg, err := env.ParseAs[ServiceConfig]()
    if err != nil {
        return nil, err
    }
    return &cfg, nil
}

// GetPostgresDsn constructs and returns
// a data source name for our Postgres connection.
// If DEV variable is set to "testing", the dsn
// creates a connection URL for testing database.
func (c *ServiceConfig) GetPostgresDsn() string {
    port := c.PgPort
    dbname := c.PgDB
    if c.Env == "testing" {
        port = c.PgTestPort
        dbname = c.PgDB + "_test"
    }
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=disable",
        c.PgUser, c.PgPass, c.PgHost, port, dbname,
    )
}

func (c *ServiceConfig) ServerAddr() string {
    return fmt.Sprintf(":%s", c.SrvPort)
}
