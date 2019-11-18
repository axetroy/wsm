package user

type SignInParams struct {
	Account  string `json:"account" valid:"required~请输入登陆账号"`
	Password string `json:"password" valid:"required~请输入密码"`
}

var Core *Service

type Service struct {
}

func New() *Service {
	return &Service{}
}

func init() {
	Core = New()
}
