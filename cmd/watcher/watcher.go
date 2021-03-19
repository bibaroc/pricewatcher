package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/bibaroc/pricewatcher/pkg/amzn"
	"gopkg.in/yaml.v2"
)

var httpc *http.Client = &http.Client{}

func main() {
	cfg, err := readConfig()
	if err != nil {
		log.Println(err)
		return
	}

	wg := sync.WaitGroup{}

	for marketPlaceURL, groups := range cfg.AmazonMarkerplaces {
		wg.Add(1)

		go func(url string, c amzn.Category, blacklist map[string]*string) {
			defer wg.Done()

			client := &http.Client{}

			for {
				for groupName, group := range c {
					for _, item := range group.Items {
						req, err := makeAMZNRequest(url, item.ASIN)
						if err != nil {
							log.Println(err)
							return
						}

						resBody, err := get200ResBody(client, req)
						if err != nil {
							log.Printf("failed to get response body for %s on %s: %s\n", item.ASIN, url, err.Error())
							return
						}

						offers, err := getAmazonOffers(bytes.NewReader(resBody))
						if err != nil {
							log.Printf("failed to get offers for %s: %s\n", item.ASIN, err.Error())
							return
						}

						for _, offer := range offers {
							if offer.Price < *item.MaxPrice {
								if _, ok := blacklist[offer.ShippingFrom]; !ok {
									fmt.Printf("found offer for product %s:%q on %s\n", groupName, item.ASIN, url)
									fmt.Printf("price=%d shipper=%s condition=%q notes=%q\n", offer.Price, offer.ShippingFrom, offer.Condition, offer.Notes)
								}
							}
						}
					}
				}
			}
		}(marketPlaceURL, groups, cfg.AmazonSellerBlacklist)
	}

	wg.Wait()
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

type (
	Config struct {
		amzn.Config
	}
)

func (c Config) String() string {
	return fmt.Sprintf("%#v", c)
}
func (c Config) ApplyDefaults() {
	for name, marketplace := range c.AmazonMarkerplaces {
		for gname, group := range marketplace {
			for i, item := range group.Items {
				if item.MaxPrice == nil {
					c.AmazonMarkerplaces[name][gname].Items[i].MaxPrice = group.MaxPrice
				}
			}
		}
	}
}
