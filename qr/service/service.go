package logic

var Logic = NewLogic()

type ILogic interface {
}

func NewLogic() ILogic {
	return &logic{}
}

type logic struct {
}
