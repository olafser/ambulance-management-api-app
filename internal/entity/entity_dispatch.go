package entity

import "time"

type DispatchEntity struct {
	DispatchID          int64     `bson:"dispatchId"`
	IncidentNumber      string    `bson:"incidentNumber"`
	CallerName          string    `bson:"callerName"`
	PatientName         string    `bson:"patientName"`
	StreetAddress       string    `bson:"streetAddress"`
	City                string    `bson:"city"`
	DispatchReason      string    `bson:"dispatchReason"`
	Priority            string    `bson:"priority"`
	Status              string    `bson:"status"`
	AmbulanceCallSign   string    `bson:"ambulanceCallSign"`
	DestinationHospital string    `bson:"destinationHospital"`
	DispatcherName      string    `bson:"dispatcherName"`
	CreatedAt           time.Time `bson:"createdAt"`
	UpdatedAt           time.Time `bson:"updatedAt"`
	Notes               string    `bson:"notes,omitempty"`
}
