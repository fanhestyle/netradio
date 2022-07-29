package main

import (
	"flag"
	"log"
	"net"
	"sync"

	"github.com/fanhestyle/netradio/medialib"
)

func init() {
	flag.StringVar(&medialib.Server_conf.Rcvport, "P", medialib.Server_conf.Rcvport, "specify port")
	flag.StringVar(&medialib.Server_conf.Mgroup, "M", medialib.Server_conf.Mgroup, "specify multicast group")
	flag.BoolVar(&medialib.Server_conf.Runmode, "F", medialib.Server_conf.Runmode, "specify foreground runmode")
	flag.StringVar(&medialib.Server_conf.Media_dir, "D", medialib.Server_conf.Media_dir, "specify media directory")
	flag.StringVar(&medialib.Server_conf.Ifname, "I", medialib.Server_conf.Ifname, "specify interface name")
}

func main() {

	flag.Parse()

	serverAddrString := net.JoinHostPort(medialib.Server_conf.Mgroup, medialib.Server_conf.Rcvport)
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddrString)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	channel_context_infos := medialib.GetChannelInfo()

	menu_sender := medialib.NewMenuSender(channel_context_infos)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		menu_sender.Send(conn)
	}()

	wg.Add(len(channel_context_infos))
	for i := 0; i < len(channel_context_infos); i++ {
		channel_sender := medialib.NewChannelSender(channel_context_infos[i])

		go func() {
			defer wg.Done()
			channel_sender.Send(conn)
		}()
	}

	wg.Wait()
}
