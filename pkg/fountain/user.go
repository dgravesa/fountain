package fountain

// User contains information about a Fountain user.
type User struct {
	ID       string `datastore:"userName"`
	FullName string `datastore:"fullName"`
	Email    string `datastore:"email"`
}
