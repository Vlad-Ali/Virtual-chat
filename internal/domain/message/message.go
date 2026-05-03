package message

type Type string

var (
	UndefinedType Type = ""
	UserType      Type = "user"
	ModelType     Type = "model"
	CallBackType  Type = "callback"
)

type Message struct {
	Msg  string
	Type Type
}
