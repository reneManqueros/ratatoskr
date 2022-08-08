package resolver

import (
	"context"
	"errors"
	"github.com/miekg/dns"
	loggerpkg "github.com/reneManqueros/logger"
	"log"
	"net/url"
	"ratatosk/internal/models/config"
	"strings"
	"time"
)

func (resolver *Resolver) parseUpstream(m *dns.Msg) Response {
	startTime := time.Now().UnixMilli()
	upstreamResponse, err := resolver.getResponseFromUpstream(m, resolver.Upstreams)
	finishTime := time.Now().UnixMilli()
	if err == nil {
		query := strings.TrimSuffix(m.Question[0].Name, ".")
		resolver.SetCache(query, *upstreamResponse)
		*m = *upstreamResponse
		return Response{
			HasAnswer: true,
			Query:     query,
			Status:    RESPONSEUPSTREAMED,
			TimeTaken: int(finishTime - startTime),
		}
	}
	return Response{
		HasAnswer: false,
		Query:     "",
		Status:    RESPONSEEMPTY,
	}
}

func (resolver *Resolver) getResponseFromUpstream(msg *dns.Msg, upstreams []*url.URL) (*dns.Msg, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Settings.UpstreamTimeout)*time.Millisecond)
	defer cancel()

	resch := make(chan *dns.Msg, len(upstreams))
	for _, up := range upstreams {
		go func(u *url.URL) {
			m, _, err := resolver.Client.Exchange(msg, u.Host)
			if err == nil {
				resch <- m
				return
			}
			log.Println(u.String(), err)
		}(up)
	}

	var errmsg *dns.Msg

	for i := 0; i < len(upstreams); i++ {
		select {
		case <-ctx.Done():
			loggerpkg.Logger{}.Debug("getResponseFromUpstream TIMEOUT")
			return nil, errors.New("time out")
		case m := <-resch:
			if m.MsgHdr.Rcode == dns.RcodeSuccess {
				return m, nil
			}
			errmsg = m
		}
	}
	if errmsg != nil {
		return errmsg, nil
	}
	return nil, errors.New("empty result")
}
