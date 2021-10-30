package employe

import "context"

//Employe represents an employe
type Employe struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Job        string    `json:"job"`
	EmployedAt string    `json:"employe_at"`
	Employes   []Employe `json:"employes,omitempty"`
}

// Repository describes the persistence on employe model
type Repository interface {
	GetEmployeByID(ctx context.Context, id string) (Employe, error)
}
