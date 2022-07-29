package medialib

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"time"
)

type MenuSender struct {
	menu Channel_menu_st
}

func NewMenuSender(infos []*channel_context_st) MenuSender {

	m := MenuSender{}
	m.menu.Id = CHANNEL_MENU_ID

	for i := 0; i < len(infos); i++ {
		oneItem := Channel_menu_item{}
		oneItem.Id = infos[i].id
		oneItem.Desc = infos[i].desc
		m.menu.ChannelMenus = append(m.menu.ChannelMenus, oneItem)
	}

	return m
}

func (snder MenuSender) Send(conn *net.UDPConn) {
	timer := time.NewTimer(1 * time.Second)

	for {
		timer.Reset(1 * time.Second)
		select {
		case <-timer.C:

			var buf bytes.Buffer
			encoder := gob.NewEncoder(&buf)
			err := encoder.Encode(snder.menu)
			if err != nil {
				log.Fatal(err)
			}
			_, err = conn.Write(buf.Bytes())
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
