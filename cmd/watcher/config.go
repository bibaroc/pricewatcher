package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bibaroc/pricewatcher/pkg/amzn"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Amazon amzn.Config `json:"amazon,omitempty" yaml:"amazon"`
}

func (c Config) String() string {
	return fmt.Sprintf("%#v", c)
}
func (c Config) ApplyDefaults() {
	for name, marketplace := range c.Amazon.Markerplaces {
		for gname, group := range marketplace {
			for i, item := range group.Items {
				if item.MaxPrice == nil {
					c.Amazon.Markerplaces[name][gname].Items[i].MaxPrice = group.MaxPrice
				}
			}
		}
	}
}

func readConfig() (*Config, error) {
	configFile := ""
	flag.StringVar(&configFile, "config-file", "cfg.yml", "--config-file ./cfg.yml")
	flag.Parse()

	configData, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(configData, config); err != nil {
		return nil, err
	}

	config.ApplyDefaults()

	return config, nil
}
