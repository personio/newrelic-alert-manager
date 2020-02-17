package controller

import (
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sync"
)

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs []func(manager.Manager, *sync.Mutex) error

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager, mutex *sync.Mutex) error {
	for _, f := range AddToManagerFuncs {
		if err := f(m, mutex); err != nil {
			return err
		}
	}
	return nil
}
