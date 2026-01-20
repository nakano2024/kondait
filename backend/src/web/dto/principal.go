package dto

const PrincipalContextKeyName = "principal"

const (
	ScopeCookingItems       = "cooking-items"
	ScopeCookingItemsRead   = "cooking-items.read"
	ScopeCookingItemsWrite  = "cooking-items.write"
	ScopeCookingItemsDelete = "cooking-items.delete"
)

type Principal struct {
	ActorCode string
	Scopes    []string
}
