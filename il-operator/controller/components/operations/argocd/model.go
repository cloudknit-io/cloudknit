package argocd

type (
	EntryIdentifier = string
	Permission      = string
)

const (
	Policy EntryIdentifier = "p"
	Group  EntryIdentifier = "g"
)

const (
	Allow Permission = "allow"
	Deny  Permission = "deny"
)

type RbacPolicy struct {
	Identifier EntryIdentifier
	Subject    string
	Resource   string
	Action     string
	Object     string
	Permission Permission
}

type RbacGroup struct {
	Identifier EntryIdentifier
	Group      string
	Role       string
}

type RbacMap struct {
	Policies map[string][]*RbacPolicy
	Groups   map[string][]*RbacGroup
}
