# gos
Golang OpenAPI Server template


## commands

```
curl -X POST -H "Content-Type: application/json" --data '{"name":"oreo","tag":"dalmatian"}' localhost:8080/pets 

curl -X POST -H "Content-Type: application/json" --data '{"name":"ashwin","email":"testemail@gmail.com", "insuranceId":1, "phone":"7161234567", "dob":"2021-02-18T21:54:42.123Z"}' localhost:8080/patients

```


# Requirements:
APIs needed:
- POST /login
- GET /patients
  - pagination
  - search
- POST /patients
- PUT /patients/{ID}
- GET /insurances

- Nogginways Users and Admins can be created manually

Datastore entities:
- NogginwayEmployeeEntity
  - ID
  - Type (user/admin)
  - Name
  - Password (ideally salted and hashed)
- InsuranceEntity
  - ID
  - Name
- PatientEntity
  - ID
  - Name
  - Email
  - Phone
  - InsuranceID
  - DoB