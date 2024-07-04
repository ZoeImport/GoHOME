package reply

type PongReply struct {
}

var pongBytes = []byte("+PING\r\n")

var thePongReply = new(PongReply)

func MakePongReply() *PongReply {
	return thePongReply
}

func (r *PongReply) ToBytes() []byte {
	return pongBytes
}

type OkReply struct {
}

var okBytes = []byte("+OK\r\n")

var theOkReply = new(OkReply)

func MakeOkReply() *OkReply {
	return theOkReply
}

func (r *OkReply) ToBytes() []byte {
	return okBytes
}

type NullBulkReply struct {
}

var nullBulkBytes = []byte("$-1\r\n")
var theNullBulkBytes = new(NullBulkReply)

func MakeNullBulkReply() *NullBulkReply {
	return theNullBulkBytes
}

func (n *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

type EmptyMultiBulkReply struct {
}

var emptyMultiBulkBytes = []byte("*0\r\n")
var theEmptyMultiBulkReply = new(EmptyMultiBulkReply)

func MakeEmptyMultiBulkBytes() *EmptyMultiBulkReply {
	return theEmptyMultiBulkReply
}

func (e *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

type NoReply struct {
}

var noBytes = []byte("")

var theNoReply = new(NoReply)

func MakeNoReply() *NoReply {
	return theNoReply
}

func (n *NoReply) ToBytes() []byte {
	return noBytes
}
