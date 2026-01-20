package dto

const PrincipalContextKeyName = "principal"

const (
	ScopeCookingItem       = "cooking-item"
	ScopeCookingItemRead   = "cooking-item.read"
	ScopeCookingItemWrite  = "cooking-item.write"
	ScopeCookingItemDelete = "cooking-item.delete"
)

type Principal struct {
	ActorCode string
	Scopes    []string
}
