// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package port_test

import (
	"testing"

	"github.com/gostafa/linuxdesktop/internal/port"
)

func TestRequestAddresses(t *testing.T) {
	t.Parallel()
	object := port.Object{Destination: "org.example.Service", Path: "/org/example/Object"}
	destination, path := object.Address()
	if destination != "org.example.Service" || path != "/org/example/Object" {
		t.Fatal("object address components were reordered", destination, path)
	}
	query := port.PropertyQuery{Object: object, Interface: "org.example.Interface", Name: "Version"}
	address := query.PropertyAddress()
	if address != (port.PropertyAddress{
		Destination: "org.example.Service", Path: "/org/example/Object",
		Interface: "org.example.Interface", Name: "Version",
	}) {
		t.Fatal("incorrect property request address", address)
	}
	query.Name = ""
	address = query.PropertyAddress()
	if address.Destination != destination || address.Path != path ||
		address.Interface != "org.example.Interface" || address.Name != "" {
		t.Fatal("reading all properties changed the object or interface", address)
	}
}
