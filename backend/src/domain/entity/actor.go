package entity

type Actor struct {
	Code string
	Sub  string
	// TODO: add more fields if needed by domain rules.
}

func NewActor(code string, sub string) *Actor {
	return &Actor{
		Code: code,
		Sub:  sub,
	}
}
