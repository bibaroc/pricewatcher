package amzn

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bibaroc/pricewatcher/pkg"
)

func FetchResults(
	ctx context.Context,
	domain string,
	cats Categories,
	blacklist map[string]*string,
	onOfferMatching func(string, string, string, Offer),
) func() error {
	return func() error {
		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		for {
			select {
			case <-ctx.Done():
				log.Printf("amzn.FetchResults contextCancelled for domain %s\n", domain)
				return nil
			case <-time.After(5 * time.Minute):
				for groupName, group := range cats {
					for _, item := range group.Items {
						time.Sleep(1000 * time.Millisecond)

						productPageRequest, err := makeRequest(domain, item.ASIN)
						if err != nil {
							return fmt.Errorf("could not prepare request: %w", err)
						}

						productPageBody, code, err := pkg.GetResBody(client, productPageRequest)
						if err != nil {
							return fmt.Errorf("failed to get response body for %s on %s: %w", item.ASIN, domain, err)
						}

						switch code {
						case http.StatusOK:
							offers, err := getOffers(bytes.NewReader(productPageBody))
							if err != nil {
								return fmt.Errorf("failed to read offers on product page %s on %s: %w", item.ASIN, domain, err)
							}

							for _, offer := range offers {
								if offer.Price < *item.MaxPrice {
									if _, ok := blacklist[offer.ShippingFrom]; !ok {
										onOfferMatching(domain, groupName, item.ASIN, offer)
									}
								}
							}
						default:
							log.Println(code, item.ASIN, "on", domain)
						}

					}
				}
			}
		}
	}
}
