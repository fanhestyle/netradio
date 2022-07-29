package medialib

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/time/rate"
)

const (
	CHANNEL_MENU_ID = 0
	MP3_BITRATE     = 320 * 1024
)

var idx int = 0

type channel_context_st struct {
	id      uint8
	desc    string
	mp3list []string
	pos     int
	fd      *os.File
	offset  uint64
	limiter *rate.Limiter
}

type Channel_menu_st struct {
	Id           uint8
	ChannelMenus []Channel_menu_item
}

type Channel_menu_item struct {
	Id   uint8
	Desc string
}

type Channel_Data struct {
	Id   uint8
	Data []byte
}

func GetChannelInfo() []*channel_context_st {

	var channel_info []*channel_context_st

	mediaFS := os.DirFS(Server_conf.Media_dir)
	subDirList, err := fs.ReadDir(mediaFS, ".")
	if err != nil {
		log.Fatal(err)
	}
	if len(subDirList) == 0 {
		fmt.Fprintf(os.Stderr, "media directory is empty!")
		os.Exit(1)
	}

	for _, ls := range subDirList {
		if ls.IsDir() {
			one_channel := readChannelInfo(mediaFS, ls)
			if one_channel != nil {
				channel_info = append(channel_info, one_channel)
			}
		}
	}
	if len(channel_info) <= 0 {
		fmt.Fprintf(os.Stderr, "No Media file found")
		os.Exit(1)
	}
	return channel_info
}

func readChannelInfo(mediaFS fs.FS, ls fs.DirEntry) *channel_context_st {

	mp3FS, err := fs.Sub(mediaFS, ls.Name())
	if err != nil {
		log.Fatal(err)
	}

	listOfMp3, err := fs.Glob(mp3FS, "*.mp3")
	if err != nil {
		log.Fatal(err)
	}
	if len(listOfMp3) == 0 {
		fmt.Fprintf(os.Stderr, "Waring: lost mp3 file!\n")
		return nil
	}

	listOfTxt, err := fs.Glob(mp3FS, "*.txt")
	if err != nil {
		log.Fatal(err)
	}
	if len(listOfTxt) < 1 {
		fmt.Fprintf(os.Stderr, "Warning: lost description file!\n")
		return nil
	} else if len(listOfTxt) > 1 {
		fmt.Fprintf(os.Stderr, "Warning: too many description files, ignore\n")
		return nil
	}

	descfile := filepath.Join(Server_conf.Media_dir, ls.Name(), listOfTxt[0])
	fd, err := os.Open(descfile)
	if err != nil {
		log.Fatal(err)
	}
	descContent, err := io.ReadAll(fd)
	if err != nil {
		log.Fatal(err)
	}

	var mp3listAbs []string
	for _, mp3Name := range listOfMp3 {
		absfile := filepath.Join(Server_conf.Media_dir, ls.Name(), mp3Name)
		mp3listAbs = append(mp3listAbs, absfile)
	}

	idx++
	var one_channel_context channel_context_st
	one_channel_context.id = uint8(idx)
	one_channel_context.desc = string(descContent)
	one_channel_context.mp3list = mp3listAbs
	one_channel_context.offset = 0
	one_channel_context.pos = 0

	curFD, err := os.Open(one_channel_context.mp3list[one_channel_context.pos])
	if err != nil {
		log.Fatal(err)
	}
	one_channel_context.fd = curFD
	one_channel_context.limiter = rate.NewLimiter(rate.Every(1*time.Second), 1)

	return &one_channel_context
}

/*
type channel_context_st struct {
	id      uint8
	desc    string
	mp3list []string
	pos     int
	fd      *os.File
	offset  uint64
	limiter *rate.Limiter
}
*/

func (chn *channel_context_st) ReadN(bytes int) []byte {

	b := make([]byte, bytes)
	//
	err := chn.limiter.Wait(context.TODO())
	if err != nil {
		fmt.Println("Error: ", err)
	}

	n, err := chn.fd.ReadAt(b, int64(chn.offset))
	if err != nil {
		log.Fatal(err)
	}

	chn.offset += uint64(n)

	if n == 0 {
		chn.pos++
		if chn.pos == len(chn.mp3list) {
			return nil
		} else {
			curFD, err := os.Open(chn.mp3list[chn.pos])
			if err != nil {
				log.Fatal(err)
			}
			chn.fd = curFD
			chn.offset = 0
		}
	}

	return b
}
