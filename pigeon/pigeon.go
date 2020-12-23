package pigeon

type M map[string]string

type Request struct {
	URL    string
	Query  M
	Header M
	Body   []byte
}

type Response struct {
	Status int
	Text   string
}

func Send(request Request) (Response, error) {
	l := GetLogger()
}
