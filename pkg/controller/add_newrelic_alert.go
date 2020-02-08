package controller

import (
	"github.com/fpetkovski/newrelic-operator/pkg/controller/newrelic_alert"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, newrelic_alert.Add)
}
