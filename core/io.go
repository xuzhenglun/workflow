package core

type ReponseWriter interface {
	Write([]byte) (int, error)
}

type Request struct {
	Name string
	Type string
	Args string
}

type Response struct {
	Status int
	Body   string
}

func (this *Response) Write(b []byte) (int, error) {
	this.Body = string(b)
	return len(b), nil
}

type CoreIoBus interface {
	RequestHandler(ReponseWriter, *Request) error
	GetMapper() map[string]string
}
