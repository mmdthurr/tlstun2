package tunn

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

type Cli struct {
	NodeName string
	AuthKey  string
	BckPort  string
	NCPort   string
}

func (c Cli) Ltocnc() {
	http.HandleFunc(fmt.Sprintf("/%s", c.AuthKey), func(w http.ResponseWriter, r *http.Request) {
		go c.StartConn(fmt.Sprintf("%s:4443", strings.Split(r.RemoteAddr, ":")[0]))
	})

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", c.NCPort), nil)
	if err != nil {
		log.Fatal(err)
	}

}

func (c Cli) StartConn(addr string) {

	conf := tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Print(err)
		return
	}

	tlsConn := tls.Client(conn, &conf)
	tlsConn.Write([]byte(fmt.Sprintf("%s_%s_", c.AuthKey, c.NodeName)))

	destConn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", c.BckPort))
	if err != nil {
		log.Printf("failed: %s", err)
	}
	go Proxy(destConn, tlsConn)

}
