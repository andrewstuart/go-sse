# sse
--
    import "github.com/andrewstuart/go-sse"


## Usage

```go
var Client = &http.Client{}
```
Client is the default client used for requests.

```go
var (
	//ErrNilChan will be returned by Notify if it is passed a nil channel
	ErrNilChan = fmt.Errorf("nil channel given")
)
```

```go
var GetReq = func(verb, uri string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(verb, uri, body)
}
```
GetReq is a function to return a single request. It will be used by notify to
get a request and can be replaces if additional configuration is desired on the
request. The "Accept" header will necessarily be overwritten.

#### func  Notify

```go
func Notify(uri string, evCh chan *Event) error
```
Notify takes a uri and channel, and will send an Event down the channel when
recieved.

#### type Event

```go
type Event struct {
	URI  string
	Type string
	Data io.Reader
}
```

Event is a go representation of an http server-sent event
