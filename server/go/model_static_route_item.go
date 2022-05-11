/*
 * Kusk Gateway API
 *
 * This is the Kusk Gateway Management API
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type StaticRouteItem struct {
	Name string `json:"name"`

	Namespace string `json:"namespace"`
}

// AssertStaticRouteItemRequired checks if the required fields are not zero-ed
func AssertStaticRouteItemRequired(obj StaticRouteItem) error {
	elements := map[string]interface{}{
		"name":      obj.Name,
		"namespace": obj.Namespace,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertRecurseStaticRouteItemRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of StaticRouteItem (e.g. [][]StaticRouteItem), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseStaticRouteItemRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aStaticRouteItem, ok := obj.(StaticRouteItem)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertStaticRouteItemRequired(aStaticRouteItem)
	})
}
