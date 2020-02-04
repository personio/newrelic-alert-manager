package controller

import (
	"github.com/fpetkovski/newrelic-operator/pkg/alert_policies/controller"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, controller.Add)
}
