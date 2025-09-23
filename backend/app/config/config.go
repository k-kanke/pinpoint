package config

import "os"

type Config struct {
    APIPort         string
    FrontendOrigin  string
    DBHost          string
    DBPort          string
    DBUser          string
    DBPassword      string
    DBName          string
}

func Load() Config {
    cfg := Config{
        APIPort:        getenv("API_PORT", "8080"),
        FrontendOrigin: getenv("FRONTEND_ORIGIN", "http://localhost:5173"),
        DBHost:         getenv("DB_HOST", "localhost"),
        DBPort:         getenv("DB_PORT", "5432"),
        DBUser:         getenv("DB_USER", "postgres"),
        DBPassword:     getenv("DB_PASSWORD", "postgres"),
        DBName:         getenv("DB_NAME", "pinpoint"),
    }
    return cfg
}

func getenv(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}

