//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=cfg.yaml ./PetStore.yaml

package petstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ashwinrs/gos/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// todo:
// figure out DoB conversion

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

// GetInsurances returns insurances
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

// GetPatients implements all the handlers in the ServerInterface
func (p *PetStoreHandler) GetPatients(w http.ResponseWriter, r *http.Request) {
	var queryResult []models.PatientEntity

	if result := p.db.Find(&queryResult); result.Error != nil {
		log.Println(result.Error)
	}
	responseResult := make([]Patient, len(queryResult))
	for i, qr := range queryResult {

		responseResult[i] = Patient{
			Id:        int64(qr.ID),
			Name:      qr.Name,
			Email:     qr.Email,
			Insurance: qr.Insurance.Name,
			Phone:     qr.Phone,
			Dob:       time.Time(qr.DateOfBirth),
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseResult)
}

func (p *PetStoreHandler) AddPatient(w http.ResponseWriter, r *http.Request) {
	var newPatient NewPatient
	if err := json.NewDecoder(r.Body).Decode(&newPatient); err != nil {
		log.Println(err)
		sendPetStoreHandlerError(w, http.StatusBadRequest, "Invalid format for NewPatient")
		return
	}

	var dbEntity models.PatientEntity
	dbEntity.Name = newPatient.Name
	dbEntity.Phone = newPatient.Phone
	dbEntity.Email = newPatient.Email
	dbEntity.InsuranceID = int(newPatient.InsuranceId)
	dbEntity.Name = newPatient.Name
	dbEntity.DateOfBirth = datatypes.Date(newPatient.Dob)

	if result := p.db.Create(&dbEntity); result.Error != nil {
		log.Println(result.Error)
	}

	response := Patient{
		Id:        int64(dbEntity.ID),
		Name:      dbEntity.Name,
		Email:     dbEntity.Email,
		Phone:     dbEntity.Phone,
		Insurance: dbEntity.Insurance.Name,
		Dob:       time.Time(dbEntity.DateOfBirth),
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// PUT /patient/{id}
func (p *PetStoreHandler) UpdatePatientByID(w http.ResponseWriter, r *http.Request, id int64) {
	var newPatient NewPatient
	if err := json.NewDecoder(r.Body).Decode(&newPatient); err != nil {
		log.Println(err)
		sendPetStoreHandlerError(w, http.StatusBadRequest, "Invalid format for NewPatient")
		return
	}

	// check if it exists
	var existingEntry models.PatientEntity
	result := p.db.First(&existingEntry, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		sendPetStoreHandlerError(w, http.StatusNotFound, fmt.Sprintf("Could not find patient with ID %d", id))
		return
	}

	var dbEntity models.PatientEntity
	dbEntity.ID = uint(id)
	dbEntity.Name = newPatient.Name
	dbEntity.Phone = newPatient.Phone
	dbEntity.Email = newPatient.Email
	dbEntity.InsuranceID = int(newPatient.InsuranceId)
	dbEntity.Name = newPatient.Name
	dbEntity.DateOfBirth = datatypes.Date(newPatient.Dob)

	if result := p.db.Save(&dbEntity); result.Error != nil {
		log.Println(result.Error)
		return
	}

	response := Patient{
		Id:        int64(dbEntity.ID),
		Name:      dbEntity.Name,
		Email:     dbEntity.Email,
		Phone:     dbEntity.Phone,
		Insurance: dbEntity.Insurance.Name,
		Dob:       time.Time(dbEntity.DateOfBirth),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
