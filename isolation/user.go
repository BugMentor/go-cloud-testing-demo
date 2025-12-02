package isolation

//go:generate mockgen -source=user.go -destination=mocks/mock_dbclient.go -package=mocks

// DBClient define la interfaz para interactuar con la base de datos
type DBClient interface {
	GetUser(id string) (string, error)
}

// UserService depende de la interfaz, no de la implementación concreta
type UserService struct {
	Db DBClient
}

func (s *UserService) GetUserName(id string) (string, error) {
	name, err := s.Db.GetUser(id)
	if err != nil {
		return "", err
	}
	return "User: " + name, nil
}
