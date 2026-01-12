package dto

const (
	ScopeCookingItemsRead = "cooking-items.read"
)

type Principal struct {
	ActorCode string
	Scopes    []string
}
