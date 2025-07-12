package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type DNSRecord struct {
	Subdomain string `json:"subdomain"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

type Response struct {
	Data []DNSRecord `json:"data"`
}

func tldSort(subdomains []string) []string {
	reversed := make([]string, len(subdomains))
	for i, sub := range subdomains {
		parts := strings.Split(sub, ".")
		for j, k := 0, len(parts)-1; j < k; j, k = j+1, k-1 {
			parts[j], parts[k] = parts[k], parts[j]
		}
		reversed[i] = strings.Join(parts, ".")
	}
	sort.Strings(reversed)

	sorted := make([]string, len(reversed))
	for i, rev := range reversed {
		parts := strings.Split(rev, ".")
		for j, k := 0, len(parts)-1; j < k; j, k = j+1, k-1 {
			parts[j], parts[k] = parts[k], parts[j]
		}
		sorted[i] = strings.Join(parts, ".")
	}
	return sorted
}

func fetchSubdomains(apiKey, domain string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://api.shodan.io/dns/domain/%s?key=%s", domain, apiKey)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	var subdomains []string
	seen := make(map[string]bool)

	for _, record := range result.Data {
		full := record.Subdomain + "." + domain
		if !seen[full] {
			seen[full] = true
			subdomains = append(subdomains, full)
		}
	}
	return subdomains, nil
}

func main() {
	domain := flag.String("d", "", "Target domain (e.g. example.com)")
	apiKey := flag.String("k", "", "Shodan API key")
	flag.Parse()

	if *domain == "" || *apiKey == "" {
		fmt.Println("Usage: ./shodan-go -k <API_KEY> -d <DOMAIN>")
		os.Exit(1)
	}

	subs, err := fetchSubdomains(*apiKey, *domain)
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("[+] Final sorted subdomains:")
	for _, sub := range tldSort(subs) {
		fmt.Println(sub)
	}
}
