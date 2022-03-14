package config

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"

	"github.com/tokopedia/tdk/go/env"
	"github.com/tokopedia/tdk/go/redis"
	yaml "gopkg.in/yaml.v2"
)

func New(repoName string) (*Config, error) {
	filename := getConfigFile(repoName)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

type Config struct {
	Server Server       `yaml:"server"`
	Redis  redis.Config `yaml:"redis"`
}

type Server struct {
	HTTP HTTP `yaml:"http"`
}

// HTTP defines server config for http server
type HTTP struct {
	Address        string `yaml:"address"`
	WriteTimeout   string `yaml:"write_timeout"`
	ReadTimeout    string `yaml:"read_timeout"`
	MaxHeaderBytes int    `yaml:"max_header_bytes"`
}

// getConfigFile get  config file name
// - files/etc/affiliate/affiliate.development.yaml in dev
// - otherwise /etc/affiliate/affiliate.{TKPENV}.yaml
func getConfigFile(repoName string) string {
	var (
		tkpEnv   = env.ServiceEnv()
		filename = fmt.Sprintf("%s.%s.yaml", repoName, tkpEnv)
	)
	// for non dev env, use config in /etc
	if tkpEnv != env.DevelopmentEnv {
		return fmt.Sprintf("/etc/%s/%s", repoName, filename)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	// use local files in dev
	repoPath := filepath.Join(gopath, "src/github.com/faithol1024", repoName)
	return filepath.Join(repoPath, "files/etc", repoName, filename)
}
