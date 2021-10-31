package employe

import "context"

//Represents an employe
type Employe struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Job        string    `json:"job"`
	EmployedAt string    `json:"employe_at"`
	Employes   []Employe `json:"employes,omitempty"`
}

// Repository describes the persistence on employe
type Repository interface {
	GetEmployeByID(ctx context.Context, id string) (Employe, error)
}
