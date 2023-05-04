//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=cfg.yaml ./PetStore.yaml

package petstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/ashwinrs/gos/internal/models"
	"gorm.io/gorm"
)

type PetStoreHandler struct {
	NextId int64
	Lock   sync.Mutex
	db     *gorm.DB
}

// Make sure we conform to ServerInterface

var _ ServerInterface = (*PetStoreHandler)(nil)

func NewPetStoreHandler(db *gorm.DB) *PetStoreHandler {
	return &PetStoreHandler{
		NextId: 1000,
		db:     db,
	}
}

// GetPets implements all the handlers in the ServerInterface
func (p *PetStoreHandler) GetPets(w http.ResponseWriter, r *http.Request, params GetPetsParams) {
	var queryResult []models.PetEntity
	limit := 10
	if params.Limit != nil {
		limit = int(*params.Limit)
	}
	if params.Tags != nil {
		if result := p.db.Limit(limit).Where("tag IN ?", *params.Tags).Find(&queryResult); result.Error != nil {
			log.Println(result.Error)
		}
	} else {
		if result := p.db.Limit(limit).Find(&queryResult); result.Error != nil {
			log.Println(result.Error)
			return
		}
	}

	response := make([]Pet, len(queryResult))
	for i, pet := range queryResult {
		response[i] = Pet{
			Id:   pet.Id,
			Name: pet.Name,
			Tag:  pet.Tag,
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (p *PetStoreHandler) AddPet(w http.ResponseWriter, r *http.Request) {
	// We expect a NewPet object in the request body.
	var newPet NewPet
	if err := json.NewDecoder(r.Body).Decode(&newPet); err != nil {
		sendPetStoreHandlerError(w, http.StatusBadRequest, "Invalid format for NewPet")
		return
	}

	// We handle pets, not NewPets, which have an additional ID field
	var pet models.PetEntity
	pet.Name = newPet.Name
	pet.Tag = newPet.Tag
	pet.Id = p.NextId
	p.NextId = p.NextId + 1

	if result := p.db.Create(&pet); result.Error != nil {
		log.Println(result.Error)
	}

	response := Pet{
		Id:   pet.Id,
		Name: pet.Name,
		Tag:  pet.Tag,
	}

	// Now, we have to return the NewPet
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (p *PetStoreHandler) FindPetByID(w http.ResponseWriter, r *http.Request, id int64) {
	var pet models.PetEntity

	// check if it exists
	result := p.db.First(&pet, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		sendPetStoreHandlerError(w, http.StatusNotFound, fmt.Sprintf("Could not find pet with ID %d", id))
		return
	}

	response := Pet{
		Id:   pet.Id,
		Name: pet.Name,
		Tag:  pet.Tag,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (p *PetStoreHandler) DeletePet(w http.ResponseWriter, r *http.Request, id int64) {
	var pet models.PetEntity

	// check if it exists
	result := p.db.First(&pet, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		sendPetStoreHandlerError(w, http.StatusNotFound, fmt.Sprintf("Could not find pet with ID %d", id))
		return
	}

	// delete if it does exist
	if result = p.db.Delete(&pet, id); result.Error != nil {
		log.Println(result.Error)
	}

	w.WriteHeader(http.StatusNoContent)
}

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendPetStoreHandlerError(w http.ResponseWriter, code int, message string) {
	petErr := Error{
		Code:    int32(code),
		Message: message,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(petErr)
}
