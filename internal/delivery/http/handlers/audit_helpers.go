package handlers

import (
	"net"
	"net/http"
	"strings"
)

func requestIP(r *http.Request) *string {
	headers := []string{
		"X-Forwarded-For",
		"X-Real-IP",
	}

	for _, header := range headers {
		value := strings.TrimSpace(r.Header.Get(header))
		if value == "" {
			continue
		}

		parts := strings.Split(value, ",")
		ip := strings.TrimSpace(parts[0])
		if ip != "" {
			return &ip
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return &host
	}

	if r.RemoteAddr != "" {
		ip := r.RemoteAddr
		return &ip
	}

	return nil
}

func requestUserAgent(r *http.Request) *string {
	ua := strings.TrimSpace(r.UserAgent())
	if ua == "" {
		return nil
	}

	return &ua
}
