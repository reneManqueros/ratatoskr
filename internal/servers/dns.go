package servers

import (
	"fmt"
	"github.com/miekg/dns"
	loggerpkg "github.com/reneManqueros/logger"
	"net/url"
	"ratatosk/internal/models/config"
	"ratatosk/internal/models/resolver"
	"time"
)

func DNS() {
	var upstreams []*url.URL

	for _, a := range config.Settings.UpstreamServers {
		u, err := url.Parse(a)
		if err != nil {
			loggerpkg.Logger{}.Fatal(err)
		}
		upstreams = append(upstreams, u)
	}

	u, err := url.Parse(fmt.Sprintf("udp://%s:%v", config.Settings.DNSAddress, config.Settings.DNSPort))
	if err != nil {
		loggerpkg.Logger{}.Fatal(err)
	}
	srv := &resolver.Resolver{
		URL:       u,
		Upstreams: upstreams,
		Client: &dns.Client{
			Net:     "udp",
			Timeout: time.Duration(config.Settings.UpstreamTimeout) * time.Millisecond,
		},
		Cache: map[string]dns.Msg{},
	}
	srv.Serve()
}
