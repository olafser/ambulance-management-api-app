package entity

type VehicleEntity struct {
	VehicleID       int64  `bson:"vehicleId"`
	CallSign        string `bson:"callSign"`
	VehicleType     string `bson:"vehicleType"`
	PlateNumber     string `bson:"plateNumber"`
	Station         string `bson:"station"`
	AssignedCrew    string `bson:"assignedCrew,omitempty"`
	Status          string `bson:"status"`
	MileageKm       int32  `bson:"mileageKm"`
	LastServiceDate string `bson:"lastServiceDate"`
	Notes           string `bson:"notes,omitempty"`
}
