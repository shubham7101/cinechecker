package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

var providerFiles = [2]string{"anime-providers.json", "movies-providers.json"}

var client = http.Client{
	CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

type Provider struct {
	Url string `json:"url"`
}

func main() {
	for _, filePath := range providerFiles {
		fmt.Printf("============%s============\n", filePath)
		checkFile(filePath)
	}
	fmt.Println("Checking completed")
}

func checkFile(filePath string) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "not able to open file %s file: %v\n", filePath, err)
		os.Exit(1)
	}

	var providers map[string]*Provider
	if err := json.NewDecoder(file).Decode(&providers); err != nil {
		fmt.Fprintf(os.Stderr, "error decoding json: %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	changed := false
	errCh := make(chan error, len(providers))
	ctx := context.Background()

	for pn, p := range providers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			URL, err := checkUrl(ctx, p.Url)
			if err != nil {
				errCh <- err
				return
			}
			if !changed && URL != p.Url {
				changed = true
			}
			providers[pn].Url = URL
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		fmt.Fprintf(os.Stderr, "not able to check url: %v\n", err)
	}

	if changed {
		if err := file.Truncate(0); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		file.Seek(0, 0)
		if err := json.NewEncoder(file).Encode(providers); err != nil {
			fmt.Fprintf(os.Stderr, "error encoding json: %v\n", err)
		}
	}

	for pn, p := range providers {
		fmt.Printf("name=%s url=%s\n", pn, p.Url)
	}
}

func checkUrl(ctx context.Context, URL string) (string, error) {
	// could use method head but get is more reliable
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return URL, err
	}
	res, err := client.Do(req)
	if err != nil {
		return URL, err
	}
	defer res.Body.Close()
	location := res.Header.Get("Location")
	if location == "" {
		return URL, nil
	}
	newURL, err := url.Parse(location)
	if err != nil {
		return URL, err
	}

	return fmt.Sprintf("%s://%s", newURL.Scheme, newURL.Host), nil
}
