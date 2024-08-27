package reply

type UnknownErrorReply struct {
}

var unknownBytes = []byte("-Err unknown\r\n")
var theUnknownErrorReply = new(UnknownErrorReply)

func (u *UnknownErrorReply) Error() string {
	return "Err unknown"
}

func (u *UnknownErrorReply) ToBytes() []byte {
	return unknownBytes
}

func MakeUnknownReply() *UnknownErrorReply {
	return theUnknownErrorReply
}

type ArgNumberErrorReply struct {
	Cmd string
}

func (a *ArgNumberErrorReply) Error() string {
	return "ERR wrong number of arguments for '" + a.Cmd + "' command\r\n"
}

func (a *ArgNumberErrorReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" + a.Cmd + "' command\r\n")

}

func MakeArgNumberErrorReply(cmd string) *ArgNumberErrorReply {
	return &ArgNumberErrorReply{Cmd: cmd}
}

type SyntaxErrorReply struct {
}

var syntaxBytes = []byte("-ERR syntax error\r\n")
var syntaxErrorReply = new(SyntaxErrorReply)

func (s *SyntaxErrorReply) Error() string {
	return "Syntax Error"
}

func (s *SyntaxErrorReply) ToBytes() []byte {
	return syntaxBytes
}

func MakeSyntaxErrorReply() *SyntaxErrorReply {
	return syntaxErrorReply
}

type WrongTypeErrorReply struct {
}

var wrongTypeBytes = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")
var wrongTypeErrorReply = new(WrongTypeErrorReply)

func (w *WrongTypeErrorReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

func (w *WrongTypeErrorReply) ToBytes() []byte {
	return wrongTypeBytes
}

func MakeWrongTypeErrorReply() *WrongTypeErrorReply {
	return wrongTypeErrorReply
}

type ProtocolErrorReply struct {
	Msg string
}

func (p *ProtocolErrorReply) Error() string {
	return "ERR Protocol error" + p.Msg
}

func (p *ProtocolErrorReply) ToBytes() []byte {
	return []byte("ERR Protocol error:'" + p.Msg + "'\r\n")
}
