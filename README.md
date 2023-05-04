# gos
Golang OpenAPI Server template


## commands

```
curl -X POST -H "Content-Type: application/json" --data '{"name":"oreo","tag":"dalmatian"}' localhost:8080/pets 
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