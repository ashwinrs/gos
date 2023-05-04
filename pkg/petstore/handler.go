//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=cfg.yaml ./PetStore.yaml

package petstore

import (
	"net/http"

	"gorm.io/gorm"
)

type PetStoreHandler struct {
	db *gorm.DB
}

// Make sure we conform to ServerInterface
var _ ServerInterface = (*PetStoreHandler)(nil)

func NewPetStoreHandler(db *gorm.DB) *PetStoreHandler {
	return &PetStoreHandler{
		db: db,
	}
}

// GetPets implements all the handlers in the ServerInterface
func (p *PetStoreHandler) GetPets(w http.ResponseWriter, r *http.Request) {
}
