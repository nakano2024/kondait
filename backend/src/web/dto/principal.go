package dto

const PrincipalContextKeyName = "principal"

const (
	ScopeCookingItemsRead = "cooking-items.read"
)

type Principal struct {
	ActorCode string
	Scopes    []string
}
