package controller

import (
	"github.com/tufin/orca-operator/pkg/controller/orca"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, orca.Add)
}
