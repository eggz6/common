package ret

type Code int

const (
	Success      Code = 200
	InvalidParam Code = 400
	NoFound      Code = 404
	ServerError  Code = 500
)

var Msg map[Code]string = map[Code]string{
	Success:      "success",
	InvalidParam: "invalid params",
	NoFound:      "no found resource",
	ServerError:  "server internal error",
}

func (c Code) Stat() int {
	return int(c)
}

func (c Code) String() string {
	res, ok := Msg[c]
	if !ok {
		return "unknown"
	}

	return res
}
