// Gostafa 2026.
// SPDX-License-Identifier: Apache-2.0.

package port

// Address returns the destination and path identifying the bus object.
func (object *Object) Address() (destination, path string) {
	return object.Destination, object.Path
}

// PropertyAddress returns the complete address used to read one or all properties.
func (query *PropertyQuery) PropertyAddress() PropertyAddress {
	return PropertyAddress{
		Destination: query.Object.Destination,
		Path:        query.Object.Path,
		Interface:   query.Interface,
		Name:        query.Name,
	}
}
