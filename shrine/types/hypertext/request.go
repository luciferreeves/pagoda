package hypertext

type Param struct {
	Key   string
	Value string
}

type Request struct {
	Path        string
	Method      string
	Query       []Param
	Params      []Param
	Headers     []Param
	QueryString string
	IP          string
	URL         string
}