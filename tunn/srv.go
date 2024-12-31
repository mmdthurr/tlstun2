package tunn

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"strings"
)

type Srv struct {
	LSrvAddr string
	LCliAddr string
	AuthKey  string
	Tlskey   string
	Tlscert  string
}

type ChanConn struct {
	Conn net.Conn
	Name string
}

var Connq = make(chan ChanConn)
var SrvMap = make(map[string]string)

func HandleLSrv(conn net.Conn, passwd string) {
	AuthBuff := make([]byte, 4096)
	_, _ = conn.Read(AuthBuff)
	hello_buff := strings.Split(string(AuthBuff), "_")
	if hello_buff[0] == passwd {
		Connq <- ChanConn{Conn: conn, Name: hello_buff[1]}

	} else {
		log.Printf("%s\n", "HELLO THIS IS IT\n")
		rsp := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\nContent-Length: 45\r\n\r\nfaghat heyvona ba ghanon jangal okht migiran!"
		_, _ = conn.Write([]byte(rsp))

	}

}

func (s Srv) LSrv() {

	go LCli(s.LCliAddr)

	cert, err := tls.LoadX509KeyPair(s.Tlscert, s.Tlskey)
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	conf := tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	l, err := net.Listen("tcp", s.LSrvAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
		}
		go HandleLSrv(tls.Server(conn, &conf), s.AuthKey)
	}

}

func LCli(addr string) {

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	for {

		conn, err := l.Accept()
		if err != nil {
			log.Printf("Lcli err %s", err)
		}

		buff := make([]byte, 4096)
		rn, _ := conn.Read(buff)

		// tst mode
		//		buff = append(buff, []byte("\r\nHost: kdk.l.ir\r\n")...)

		spd := strings.Split(string(buff), "\r\n")
		// log.Printf("%s", spd[10](buff), "\r\n"))
		for i := 0; i < len(spd); i++ {
			if strings.HasPrefix(spd[i], "Host: ") {
				rhost := strings.TrimPrefix(spd[i], "Host: ")

				// customize it based on your domain since my domain be something like this kkdfs.usa.choskosh.cfd then [1] would result in usa
				pk := strings.Split(rhost, ".")[1]

				log.Printf("%s", pk)

				raddr, ok := SrvMap[pk]
				if ok {
					http.Get(raddr)
					var srvconn net.Conn
					for {
						conns := <-Connq

						if conns.Name == pk {
							srvconn = conns.Conn
							srvconn.Write(buff[:rn])
							go Proxy(conn, srvconn)
							break
						} else {
							Connq <- conns
						}
					}
				}
				break
			}

		}
	}

}
