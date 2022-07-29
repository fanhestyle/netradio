package medialib

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type ChannelSender struct {
	channel_info *channel_context_st
}

func NewChannelSender(chninfo *channel_context_st) ChannelSender {
	channelSender := ChannelSender{}
	channelSender.channel_info = chninfo
	return channelSender
}

func (snder ChannelSender) Send(conn *net.UDPConn) {

	for {
		//mp3 rate //320 * 1024 / 8 // 320kbit/s
		const MP3_RATE = 320 * 1024 / 8
		data := snder.channel_info.ReadN(MP3_RATE)

		if data == nil {
			break
		}

		channel_data := Channel_Data{}
		channel_data.Id = snder.channel_info.id
		channel_data.Data = data

		var buf bytes.Buffer
		encoder := gob.NewEncoder(&buf)
		err := encoder.Encode(channel_data)
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("send channel %v-%v\n", channel_data.Id, len(data))
	}
}
