package resolver

import (
	"github.com/miekg/dns"
	loggerpkg "github.com/reneManqueros/logger"
	"log"
	"net"
	"net/url"
	"strconv"
	"sync"
)

type Resolver struct {
	Client     *dns.Client
	URL        *url.URL
	CacheMutex sync.Mutex
	Cache      map[string]dns.Msg
	Logger     chan Response
	Upstreams  []*url.URL
}

const (
	RESPONSECACHED     = "cached"
	RESPONSEBLOCKED    = "blocked"
	RESPONSEUPSTREAMED = "upstreamed"
	RESPONSELOCAL      = "local"
	RESPONSEEMPTY      = ""
)

type Response struct {
	Query     string
	Status    string
	Client    string
	TimeTaken int
	HasAnswer bool
}

func (resolver *Resolver) Serve() {

	ip, port, _ := net.SplitHostPort(resolver.URL.Host)
	_ip := net.ParseIP(ip)
	_port, _ := strconv.Atoi(port)

	udpconn, err := net.ListenUDP("udp", &net.UDPAddr{IP: _ip, Port: _port})
	if err != nil {
		log.Fatalf("listen udp error %s", err)
	}

	defer udpconn.Close()

	buf := make([]byte, 4096)
	for {
		n, addr, err := udpconn.ReadFrom(buf)
		if err != nil {
			log.Println(err)
			break
		}
		buf1 := make([]byte, n)
		copy(buf1, buf[:n])

		go resolver.handle(buf1, addr, udpconn)
	}
}

func (resolver *Resolver) handle(buf []byte, addr net.Addr, conn *net.UDPConn) {
	msg := new(dns.Msg)
	if err := msg.Unpack(buf); err != nil {
		loggerpkg.Logger{}.Debug("resolver.handle unpack", err)
		return
	}

	handlers := []func(m *dns.Msg) Response{
		resolver.parseCache,
		resolver.parseLocal,
		resolver.parseBlacklist,
		resolver.parseUpstream,
	}

	for _, handler := range handlers {
		response := handler(msg)
		if response.HasAnswer == true {
			resolver.SetCache(response.Query, *msg)
			response.Client = addr.String()
			break
		}
	}

	d, _ := msg.Pack()
	_, err := conn.WriteTo(d, addr)
	if err != nil {
		loggerpkg.Logger{}.Debug("resolver.handle pack", err)
	}
}
