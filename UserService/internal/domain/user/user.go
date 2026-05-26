package domainuser

const (
	BasicRole   = "basic"
	PremiumRole = "premium"
)

type User struct {
	Name     string
	Email    string
	Password string
	Role     string
}
