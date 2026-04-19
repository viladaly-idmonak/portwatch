package geoip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Info holds geolocation data for an IP address.
type Info struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	City    string `json:"city"`
	Org     string `json:"org"`
}

// Lookup resolves geolocation information for a given IP address.
type Lookup interface {
	Lookup(ip string) (Info, error)
}

// Client is an HTTP-based geolocation client.
type Client struct {
	baseURL    string
	httpClient *http.Client
	cache      map[string]Info
}

// New returns a new Client with optional base URL override.
func New(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://ipapi.co"
	}
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
		cache:      make(map[string]Info),
	}
}

// Lookup returns geolocation info for the given IP, using a local cache.
func (c *Client) Lookup(ip string) (Info, error) {
	if !isPublicIP(ip) {
		return Info{IP: ip, Country: "private", City: "", Org: ""}, nil
	}
	if info, ok := c.cache[ip]; ok {
		return info, nil
	}
	url := fmt.Sprintf("%s/%s/json/", c.baseURL, ip)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return Info{}, fmt.Errorf("geoip lookup %s: %w", ip, err)
	}
	defer resp.Body.Close()
	var info Info
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return Info{}, fmt.Errorf("geoip decode %s: %w", ip, err)
	}
	c.cache[ip] = info
	return info, nil
}

func isPublicIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	privateRanges := []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.0/8", "::1/128"}
	for _, cidr := range privateRanges {
		_, network, _ := net.ParseCIDR(cidr)
		if network.Contains(ip) {
			return false
		}
	}
	return true
}
