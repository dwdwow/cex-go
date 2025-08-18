package bnc

type UserCfg struct {
	APIKey       string
	APISecretKey string
}

type User struct {
	cfg UserCfg
}

func NewUser(cfg UserCfg) *User {
	return &User{
		cfg: cfg,
	}
}

func (u *User) NewListenKey(url string) (string, error) {
	return "", nil
}

func (u *User) KeepListenKey(url string, listenKey string) (string, error) {
	return "", nil
}
