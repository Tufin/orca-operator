package controller

import (
	"github.com/tufin/orca-operator/pkg/controller/policy"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, policy.Add)
}
