package order

import "github.com/theoremoon/kosenctfx/scoreserver/model"

type Order struct {
	Deployments []*model.Deployment `json:"deployments"`
	Retires     []*model.Deployment `json:"retires"`
}
