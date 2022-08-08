package resolver

import (
	"fmt"
	"github.com/miekg/dns"
	loggerpkg "github.com/reneManqueros/logger"
	"io/ioutil"
	"log"
	"ratatoskr/internal/models/config"
	"strings"
)

var LocalBlackList = map[string]string{}

func (resolver *Resolver) parseBlacklist(m *dns.Msg) Response {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
		case dns.TypeCNAME:
			query := strings.TrimSuffix(q.Name, ".")
			if _, ok := LocalBlackList[query]; ok {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, "255.255.255.0"))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				} else {
					loggerpkg.Logger{}.Debug("BLACKLIST ERROR", err)
				}

				return Response{
					HasAnswer: true,
					Query:     query,
					Status:    RESPONSEBLOCKED,
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

func BlacklistLoad() {
	files, err := ioutil.ReadDir(config.Settings.ContentRoot + "/data/")
	if err != nil {
		log.Fatal(err)
	}

	LocalBlackList = map[string]string{}
	for _, file := range files {
		if file.Name() == "_localentries.json" {
			continue
		}
		if file.IsDir() == false {
			lines, err := ioutil.ReadFile(config.Settings.ContentRoot + "/data/" + file.Name())
			if err != nil {
				log.Println(err)
			}
			for _, line := range strings.Split(string(lines), "\n") {
				if strings.HasPrefix(line, "#") == false && len(line) >= 9 {
					recordSplice := strings.Split(line, " ")
					if len(recordSplice) < 2 {
						log.Println(recordSplice)
					} else {
						ip, hostname := recordSplice[0], recordSplice[1]
						LocalBlackList[hostname] = ip
					}
				}
			}
		}
	}
	loggerpkg.Logger{}.Debug("LOADED to blacklist", len(LocalBlackList))
}
