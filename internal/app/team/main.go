package team

var Core *Service

type Service struct {
}

func New() *Service {
	return &Service{}
}

func init() {
	Core = New()
}
