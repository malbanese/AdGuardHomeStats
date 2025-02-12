package proxy

import "net/http"

type AuthType interface {
	Authenticate(req *http.Request)
}

// Simple authentication type for plaintext username and password
type BasicAuthType struct {
	Username string
	Password string
}

func (a *BasicAuthType) Authenticate(req *http.Request) {
	req.SetBasicAuth(a.Username, a.Password)
}
