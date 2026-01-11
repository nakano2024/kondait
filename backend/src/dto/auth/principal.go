package auth

const (
	ScopeCookingItemsRead = "cooking-items.read"
)

type Principal struct {
	UserCode string
	Scopes   []string
}
