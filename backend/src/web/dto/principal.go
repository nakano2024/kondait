package dto

const (
	ScopeCookingItemsRead = "cooking-items.read"
)

type Principal struct {
	UserCode string
	Scopes   []string
}
