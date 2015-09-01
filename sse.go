package sse

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

//SSE name constants
const (
	eName = "event"
	dName = "data"
)

//Client is the default client used for requests.
var Client = &http.Client{}

func liveReq(verb, uri string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(verb, uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/event-stream")

	return req, nil
}

//Event is a go representation of an http server-sent event
type Event struct {
	URI  string
	Type string
	Data io.Reader
}

//Notify takes a uri and channel, and will send an Event down the channel when
//recieved.
func Notify(uri string, evCh chan *Event) error {

	req, err := liveReq("GET", uri, nil)

	if err != nil {
		return fmt.Errorf("error getting sse request: %v", err)
	}

	res, err := Client.Do(req)

	if err != nil {
		return fmt.Errorf("error performing request for %s: %v", uri, err)
	}

	go func() {
		br := bufio.NewReader(res.Body)
		defer res.Body.Close()

		delim := []byte{':', ' '}

		var currEvent *Event

		for {
			bs, err := br.ReadBytes('\n')

			if err != nil {
				return
			}

			if len(bs) < 2 {
				continue
			}

			spl := bytes.Split(bs, delim)

			if len(spl) < 2 {
				continue
			}

			switch string(spl[0]) {
			case eName:
				currEvent = &Event{URI: uri}
				currEvent.Type = string(spl[1])
			case dName:
				currEvent.Data = bytes.NewBuffer(spl[1])
				evCh <- currEvent
			}
		}
	}()

	return nil
}
