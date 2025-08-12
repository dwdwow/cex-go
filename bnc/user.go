package bnc

type User struct {
}

func (u *User) NewListenKey(url string) (string, error) {
	return "", nil
}

func (u *User) KeepListenKey(url string, listenKey string) (string, error) {
	return "", nil
}
