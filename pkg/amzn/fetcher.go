package amzn

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

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
		client := &http.Client{}

		for {
			select {
			case <-ctx.Done():
				log.Printf("amzn.FetchResults contextCancelled for domain %s\n", domain)
				return nil
			default:
				for groupName, group := range cats {
					for _, item := range group.Items {
						productPageRequest, err := makeRequest(domain, item.ASIN)
						if err != nil {
							return fmt.Errorf("could not prepare request: %w", err)
						}

						productPageBody, err := pkg.Get200ResBody(client, productPageRequest)
						if err != nil {
							return fmt.Errorf("failed to get response body for %s on %s: %w", item.ASIN, domain, err)
						}

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
					}
				}
			}
		}
	}
}
