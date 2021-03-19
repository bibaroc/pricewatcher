package amzn

import (
	"context"
	"fmt"
	"net/http"
)

func makeRequest(domain, asin string) (*http.Request, error) {
	request, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		fmt.Sprintf(`https://%s/gp/aod/ajax/ref=dp_aod_afts?asin=%s`, domain, asin),
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request for %s on %s: %w", asin, domain, err)
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:86.0) Gecko/20100101 Firefox/86.0s")
	request.Header.Set("Accept-Language", "en-US,en;q=0.9")

	return request, nil
}
