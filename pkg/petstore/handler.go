//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=cfg.yaml ./PetStore.yaml

package petstore

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ashwinrs/gos/internal/models"
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
func (p *PetStoreHandler) GetInsurances(w http.ResponseWriter, r *http.Request) {
	var queryResult []models.InsuranceEntity

	if result := p.db.Find(&queryResult); result.Error != nil {
		log.Println(result.Error)
	}

	responseResult := make([]Insurance, len(queryResult))
	for i, qr := range queryResult {
		responseResult[i] = Insurance{Id: int64(qr.ID), Name: qr.Name}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseResult)
}
