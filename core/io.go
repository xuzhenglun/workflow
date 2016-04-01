package core

type ReponseWriter interface {
	Write([]byte) (int, error)
}

type Request struct {
	Name string
	Id   string
	Type string
	Args string
	Auth []byte
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
	ReloadConfig() string
	GetActivities([]byte, string) (string, error)
	ListNeedArgs([]byte, string) ([]string, error)
	GetPurviewActivities(string) ([]string, error)
}

type DataBase interface {
	AddRow(...map[string]string) error
	DeleteRow(...string) error
	FindRow(string, ...string) (string, error)
	ModifyRow(...map[string]string) error
	GetJustDone(string) (string, string, error)
	GetList(string, string) (string, error)
}
