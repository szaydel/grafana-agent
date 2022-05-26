// Package apache_http embeds https://github.com/Lusitaniae/apache_exporter
package apache_http //nolint:golint

import (
	"fmt"
	"net/url"

	ae "github.com/Lusitaniae/apache_exporter/collector"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/agent/pkg/integrations"
)

// DefaultConfig holds the default settings for the apache_http integration
var DefaultConfig = Config{
	ApacheAddr:         "http://localhost/server-status?auto",
	ApacheHostOverride: "",
	ApacheInsecure:     false,
}

// Config controls the apache_http integration.
type Config struct {
	ApacheAddr         string `yaml:"scrape_uri,omitempty"`
	ApacheHostOverride string `yaml:"host_override,omitempty"`
	ApacheInsecure     bool   `yaml:"insecure,omitempty"`
}

// UnmarshalYAML implements yaml.Unmarshaler for Config
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultConfig

	type plain Config
	return unmarshal((*plain)(c))
}

// Name returns the name of the integration this config is for.
func (c *Config) Name() string {
	return "apache_http"
}

// InstanceKey returns the addr of the apache server.
func (c *Config) InstanceKey(agentKey string) (string, error) {
	u, err := url.Parse(c.ApacheAddr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", u.Hostname(), u.Port()), nil
}

// NewIntegration converts the config into an integration instance.
func (c *Config) NewIntegration(logger log.Logger) (integrations.Integration, error) {
	return New(logger, c)
}

func init() {
	integrations.RegisterIntegration(&Config{})
}

func New(logger log.Logger, c *Config) (integrations.Integration, error) {
	conf := &ae.Config{
		ScrapeURI:    c.ApacheAddr,
		HostOverride: c.ApacheHostOverride,
		Insecure:     c.ApacheInsecure,
	}

	//check scrape URI
	_, err := url.ParseRequestURI(conf.ScrapeURI)
	if err != nil {
		level.Error(logger).Log("msg", "scrape_uri is invalid", "err", err)
		return nil, err
	}
	aeExporter := ae.NewExporter(logger, conf)

	return integrations.NewCollectorIntegration(
		c.Name(),
		integrations.WithCollectors(aeExporter),
	), nil
}
