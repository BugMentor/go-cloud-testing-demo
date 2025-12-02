package isolation

// CHANGE HERE: We use -source=user.go to avoid package compilation errors during generation
//go:generate mockgen -source=user.go -destination=mocks/mock_dbclient.go -package=mocks

// DBClient defines the interface for interacting with the database
type DBClient interface {
	GetUser(id string) (string, error)
}

// UserService depends on the interface, not the concrete implementation
type UserService struct {
	Db DBClient
}

// GetUserName retrieves a user by ID and formats their name
func (s *UserService) GetUserName(id string) (string, error) {
	name, err := s.Db.GetUser(id)
	if err != nil {
		return "", err
	}
	return "User: " + name, nil
}
