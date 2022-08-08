package resolver

import (
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	loggerpkg "github.com/reneManqueros/logger"
	"io/ioutil"
	"ratatosk/internal/models/config"
	"strings"
)

var LocalEntries = map[string]string{}

func LocalLoad() {
	b, err := ioutil.ReadFile(config.Settings.ContentRoot + "/data/_localentries.json")
	if err == nil {
		err := json.Unmarshal(b, &LocalEntries)
		if err != nil {
			loggerpkg.Logger{}.Debug("LocalLoad", err)
		}
	}

}

func (resolver *Resolver) parseLocal(m *dns.Msg) Response {
	for _, q := range m.Question {
		query := strings.TrimSuffix(q.Name, ".")
		ip := LocalEntries[query]

		if ip != "" {
			rrText := fmt.Sprintf("%s A %s", q.Name, ip)
			rr, err := dns.NewRR(rrText)
			if err == nil {
				m.Answer = append(m.Answer, rr)
				return Response{
					HasAnswer: true,
					Query:     query,
					Status:    RESPONSELOCAL,
				}
			} else {
				loggerpkg.Logger{}.Debug("parseLocal", err)
			}
		}

	}
	return Response{
		HasAnswer: false,
		Query:     "",
		Status:    RESPONSEEMPTY,
	}
}
