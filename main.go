package main

import (
	"flag"
	"fmt"
	"mmd/tlstun2/tunn"
	"strings"
)

func main() {

	mode := flag.String("m", "c", "c client , s server mode")

	v2P := flag.String("v2p", "1086", "v2ray port")
	nodename := flag.String("name", "usa", "tunnel node name")
	ncport := flag.String("ncp", "6058", "http port server to be called and  make new conn")

	tlscert := flag.String("cert", "tls.cert", "tls certificate")
	tlskey := flag.String("key", "tls.key", "tls key")
	srvaddr := flag.String("sr", "0.0.0.0:4443", "addr")
	cliaddr := flag.String("clir", "0.0.0.0:80", "cli addr")
	lofsrv := flag.String("lofs", "l_127.0.0.1:6058_,", "list of srv to call and initiate newconnfrom")

	passwd := flag.String("passwd", "123456", "tunnel passwd")

	flag.Parse()

	if *mode == "s" {
		rhost := strings.Split(*lofsrv, ",")

		for i := 0; i < len(rhost)-1; i++ {

			spd := strings.Split(rhost[i], "_")
			tunn.SrvMap[spd[0]] = fmt.Sprintf("http://%s/%s", spd[1], *passwd)
		}
		fmt.Print("iam here\n")
		s := tunn.Srv{
			LSrvAddr: *srvaddr,
			LCliAddr: *cliaddr,
			AuthKey:  *passwd,
			Tlskey:   *tlskey,
			Tlscert:  *tlscert,
		}
		s.LSrv()

	} else if *mode == "c" {

		c := tunn.Cli{
			NodeName: *nodename,
			AuthKey:  *passwd,
			BckPort:  *v2P,
			NCPort:   *ncport,
		}
		c.Ltocnc()

	}
}
