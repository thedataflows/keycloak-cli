// Package-level volatile/write-only field state and the round-trip strip
// helpers, moved verbatim from pkg/manifest/manifest.go so that pkg/kcapi owns
// the mutable state (a manifest-side snapshot copy would go stale on reinstall).
// CloneResource and NormalizeValueForRoundTrip are exported (renamed from
// cloneResource / normalizeValueForRoundTrip) because pkg/manifest's remaining
// round-trip paths read them across the package boundary.

package kcapi

import (
	"sort"
	"strings"
)

var volatileRoundTripFields = map[string]struct{}{
	"id":               {},
	"internalid":       {},
	"containerid":      {},
	"access":           {},
	"origin":           {},
	"createdtimestamp": {},
}

var writeOnlyResourceFields = map[string]map[string]struct{}{
	"user":   {"credentials": {}},
	"client": {"clientSecret": {}},
}

// InstallVolatileFields replaces the global set of volatile round-trip field
// names. It is used by the catalog package after loading field overrides.
func InstallVolatileFields(fields map[string]map[string]struct{}) {
	merged := make(map[string]struct{})
	for name := range volatileRoundTripFields {
		merged[name] = struct{}{}
	}
	for _, set := range fields {
		for name := range set {
			merged[strings.ToLower(name)] = struct{}{}
		}
	}
	volatileRoundTripFields = merged
}

// InstallWriteOnlyFields replaces the per-resource-type write-only field sets.
// It is used by the catalog package after loading field overrides.
func InstallWriteOnlyFields(fields map[string]map[string]struct{}) {
	writeOnlyResourceFields = fields
}

// StripVolatileFields returns a copy of the resource with server-managed and
// write-only fields removed, suitable for validation before apply.
func StripVolatileFields(resource Resource) Resource {
	return CloneResource(resource, nil)
}

// CloneResource returns a copy of the resource with server-managed and volatile
// fields removed. If preserveIDs is non-nil, any resource whose "id" appears in
// the set keeps its id so that inline references from other resources can be
// remapped during apply.
func CloneResource(resource Resource, preserveIDs map[string]struct{}) Resource {
	clone := Resource{
		Type:       resource.Type,
		Realm:      strings.TrimSpace(resource.Realm),
		Delete:     resource.Delete,
		ParentType: resource.ParentType,
	}
	if len(resource.Data) > 0 {
		clone.Data = normalizeMapForRoundTrip(resource.Data, true)
	}
	if clone.Realm == "" {
		clone.Realm = rawStringField(clone.Data, "realm")
	}
	if clone.Data != nil {
		if fields, ok := writeOnlyResourceFields[clone.Type]; ok {
			deleteWriteOnlyFields(clone.Data, fields)
		}
		if preserveIDs != nil {
			if id := rawStringField(resource.Data, "id"); id != "" {
				if _, ok := preserveIDs[id]; ok {
					clone.Data["id"] = id
				}
			}
		}
	}
	return clone
}

func deleteWriteOnlyFields(data map[string]interface{}, fields map[string]struct{}) {
	for field := range fields {
		delete(data, field)
	}
	for _, value := range data {
		switch typed := value.(type) {
		case map[string]interface{}:
			deleteWriteOnlyFields(typed, fields)
		case []interface{}:
			for _, item := range typed {
				if m, ok := item.(map[string]interface{}); ok {
					deleteWriteOnlyFields(m, fields)
				}
			}
		}
	}
}

func normalizeMapForRoundTrip(values map[string]interface{}, stripVolatile bool) map[string]interface{} {
	if len(values) == 0 {
		return nil
	}

	result := make(map[string]interface{}, len(values))
	for key, value := range values {
		if stripVolatile && isVolatileField(key) {
			continue
		}
		result[key] = NormalizeValueForRoundTrip(value, stripVolatile)
	}
	return result
}

func isVolatileField(name string) bool {
	_, ok := volatileRoundTripFields[strings.ToLower(strings.TrimSpace(name))]
	return ok
}

// NormalizeValueForRoundTrip recursively normalizes a parsed JSON value,
// optionally stripping volatile fields; exported for pkg/manifest's
// relationship round-trip path (renamed from normalizeValueForRoundTrip).
func NormalizeValueForRoundTrip(value interface{}, stripVolatile bool) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		return normalizeMapForRoundTrip(typed, stripVolatile)
	case []interface{}:
		normalized := make([]interface{}, len(typed))
		keys := make([]string, len(typed))
		canSort := true
		for idx, item := range typed {
			normalized[idx] = NormalizeValueForRoundTrip(item, stripVolatile)
			key, ok := sortIdentity(normalized[idx])
			if !ok {
				canSort = false
				continue
			}
			keys[idx] = key
		}
		if canSort {
			sort.SliceStable(normalized, func(left, right int) bool {
				return keys[left] < keys[right]
			})
		}
		return normalized
	default:
		return value
	}
}

func sortIdentity(value interface{}) (string, bool) {
	mapped, ok := value.(map[string]interface{})
	if !ok {
		return "", false
	}
	for _, key := range []string{"name", "username", "clientId", "alias", "realm", "provider", "identityProvider", "path", "type"} {
		if s, ok := mapped[key].(string); ok && strings.TrimSpace(s) != "" {
			return key + ":" + strings.TrimSpace(s), true
		}
	}
	return "", false
}
