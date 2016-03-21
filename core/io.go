package core

type ReponseWriter interface {
	Write([]byte) (int, error)
}

type Request struct {
	Name string
	Id   string
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
	RawLuaHandler(string) string
	GetMapper() map[string]string
}

type DataBase interface {
	AddRow([]string, map[string]interface{}) error
	DeleteRow(int) error
	FindRow(int) (string, error)
	ModifyRow([]string, map[string]interface{}) error
	GetFather(string) (string, error)
}
