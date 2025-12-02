package isolation

import (
	"errors"
	"testing"

	"github.com/BugMentor/go-cloud-testing-demo/isolation/mocks"
	"github.com/golang/mock/gomock"
)

func TestGetUserName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Inyectamos el Mock generado
	mockDB := mocks.NewMockDBClient(ctrl)
	service := UserService{Db: mockDB}

	// Escenario 1: Éxito
	mockDB.EXPECT().GetUser("123").Return("Matías", nil)

	result, err := service.GetUserName("123")
	if result != "User: Matías" || err != nil {
		t.Errorf("Error inesperado: %v", err)
	}

	// Escenario 2: Fallo (Edge case)
	mockDB.EXPECT().GetUser("999").Return("", errors.New("DB error"))

	_, err = service.GetUserName("999")
	if err == nil {
		t.Error("Se esperaba un error cuando la DB falla")
	}
}
