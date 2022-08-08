package resolver

import (
	"github.com/miekg/dns"
	"strings"
)

func (resolver *Resolver) SetCache(query string, msg dns.Msg) {
	resolver.CacheMutex.Lock()
	resolver.Cache[query] = msg
	resolver.CacheMutex.Unlock()
}

func (resolver *Resolver) GetCache(query string) dns.Msg {
	resolver.CacheMutex.Lock()
	msg := resolver.Cache[query]
	resolver.CacheMutex.Unlock()
	return msg
}

func (resolver *Resolver) InCache(query string) (dns.Msg, bool) {
	resolver.CacheMutex.Lock()
	value, ok := resolver.Cache[query]
	resolver.CacheMutex.Unlock()
	return value, ok
}

func (resolver *Resolver) parseCache(m *dns.Msg) Response {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
		case dns.TypeCNAME:
			query := strings.TrimSuffix(q.Name, ".")
			if value, ok := resolver.InCache(query); ok {
				m.Answer = value.Answer
				return Response{
					HasAnswer: true,
					Query:     query,
					Status:    RESPONSECACHED,
				}
			}
		}
	}
	return Response{
		HasAnswer: false,
		Query:     "",
		Status:    RESPONSEEMPTY,
	}
}
