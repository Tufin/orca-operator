package orca_test

import (
	"testing"

	"github.com/tufin/orca-operator/pkg/controller/orca"
	"github.com/stretchr/testify/require"
)

func TestGetLabels(t *testing.T) {

	labels := orca.GetLabels("app=kite", "env=dev")

	require.Equal(t, "kite", labels["app"])
	require.Equal(t, "dev", labels["env"])
}
