package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"

	"github.com/fanhestyle/netradio/medialib"
)

type client_conf_st struct {
	rcvport    string
	mgroup     string
	player_cmd string
}

var client_conf = &client_conf_st{
	rcvport:    medialib.DEFAULT_PORT,
	mgroup:     medialib.DEFAULT_MGROUP,
	player_cmd: medialib.DEFAULT_PLAYERCMD,
}

func init() {
	flag.StringVar(&client_conf.rcvport, "P", client_conf.rcvport, "specify receive port")
	flag.StringVar(&client_conf.mgroup, "M", client_conf.mgroup, "specify multicast group")
	flag.StringVar(&client_conf.player_cmd, "p", client_conf.player_cmd, "specify player")

	flag.StringVar(&client_conf.rcvport, "port", client_conf.rcvport, "specify receive port")
	flag.StringVar(&client_conf.mgroup, "mgroup", client_conf.mgroup, "specify multicast group")
	flag.StringVar(&client_conf.player_cmd, "player", client_conf.player_cmd, "specify player")
}

func main() {

	flag.Parse()

	udpAddrString := net.JoinHostPort(client_conf.mgroup, client_conf.rcvport)

	udpAddr, err := net.ResolveUDPAddr("udp", udpAddrString)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	var chosenChannelNum int

	for {
		buf := make([]byte, 65535)
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Fatal(err)
		}

		dec := gob.NewDecoder(bytes.NewReader(buf[:n]))
		q := medialib.Channel_menu_st{}

		err = dec.Decode(&q)
		if err != nil {
			log.Fatal(err)
		}

		if q.Id != medialib.CHANNEL_MENU_ID {
			continue
		}

		println("Please choose a channel:")

		channel_menus := q.ChannelMenus
		for _, oneItem := range channel_menus {
			fmt.Print(oneItem.Id, ": ", oneItem.Desc)
		}

		var channelID int
		_, err = fmt.Scanf("%d", &channelID)
		if err != nil {
			fmt.Println("Invalid Input, try again: ")
			continue
		}

		if channelID > len(channel_menus) || channelID <= 0 {
			fmt.Println("Invalid Chosen Number, try again: ")
			continue
		}

		chosenChannelNum = channelID
		break
	}

	fmt.Printf("You are now listening to channel %d\n", chosenChannelNum)

	cmd := exec.Command("mpg123", "-")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			buf := make([]byte, 65535)
			n, err := conn.Read(buf[:])
			if err != nil {
				log.Fatal(err)
			}

			dec := gob.NewDecoder(bytes.NewReader(buf[:n]))
			q := medialib.Channel_Data{}

			err = dec.Decode(&q)
			if err != nil {
				log.Fatal(err)
			}

			if q.Id != uint8(chosenChannelNum) {
				continue
			}

			buff := bytes.NewBuffer(q.Data)
			io.Copy(stdin, buff)
		}

	}()
	cmd.Run()
}
