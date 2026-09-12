package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Security SecurityConfig `yaml:"security"`
	TLS      TLSConfig      `yaml:"tls"`
}

type AppConfig struct {
	Name string `yaml:"name"`
	Env  string `yaml:"env"`
}

type ServerConfig struct {
	Host                string `yaml:"host"`
	Port                int    `yaml:"port"`
	WSPath              string `yaml:"wsPath"`
	ReadBuffer          int    `yaml:"readBuffer"`
	WriteBuffer         int    `yaml:"writeBuffer"`
	PingIntervalSeconds int    `yaml:"pingIntervalSeconds"`
}

func (s ServerConfig) Addr() string { return fmt.Sprintf("%s:%d", s.Host, s.Port) }

type DatabaseConfig struct{ URL string `yaml:"url"` }

type SecurityConfig struct {
	JWTSecret     string `yaml:"jwtSecret"`
	JWTTTLMinutes int    `yaml:"jwtTtlMinutes"`
}

func (s SecurityConfig) JWTTTL() time.Duration {
	return time.Duration(s.JWTTTLMinutes) * time.Minute
}

type TLSConfig struct {
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

func (t TLSConfig) Enabled() bool { return t.CertFile != "" && t.KeyFile != "" }

// Load reads a yaml file and expands ${VAR} / ${VAR:-default} from the env.
func Load(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal([]byte(interpolate(string(raw))), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func interpolate(s string) string {
	for {
		i := strings.Index(s, "${")
		if i < 0 {
			return s
		}
		j := strings.Index(s[i:], "}")
		if j < 0 {
			return s
		}
		expr := s[i+2 : i+j]
		name, def := expr, ""
		if k := strings.Index(expr, ":-"); k >= 0 {
			name, def = expr[:k], expr[k+2:]
		}
		v := os.Getenv(name)
		if v == "" {
			v = def
		}
		s = s[:i] + v + s[i+j+1:]
	}
}