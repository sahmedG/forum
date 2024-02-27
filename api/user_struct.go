package api

type SignUp_form struct {
	Email    string
	Uname    string
	Password string
	EX_ID    string
}

type LogIn_form struct {
	Email    string
	Password string
	EX_ID    string
}

type Redirect struct {
	URL     string `json:"URL"`
	Success bool   `json:"success"`
}

type CallbackGH struct {
	URLS string `json:"URLS"`
}

var Code string

type Config struct {
	CertFile     string `json:"cert_file"`
	KeyFile      string `json:"key_file"`
	ServerAddr   string `json:"server_addr"`
	ReadTimeout  string `json:"read_timeout"`
	WriteTimeout string `json:"write_timeout"`
	IdleTimeout  string `json:"idle_timeout"`
}
