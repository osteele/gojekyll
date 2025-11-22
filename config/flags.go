package config

import (
	"reflect"
)

// Flags are applied after the configuration file is loaded.
// They are pointers to represent optional types, to tell whether they have been set.
type Flags struct {
	// these are pointers so we can tell whether they've been set, and leave
	// the config file alone if not
	Destination, Host, BaseURL  *string
	Drafts, Future, Unpublished *bool
	Incremental, Verbose        *bool
	Port                        *int

	// these aren't in the config file, so make them actual values
	DryRun, ForcePolling, Watch bool
	ConfigFile                  string
}

// ApplyFlags overwrites the configuration with values from flags.
func (c *Config) ApplyFlags(f Flags) {
	// anything you can do I can do meta
	rs, rd := reflect.ValueOf(f), reflect.ValueOf(c).Elem()
	rt := rs.Type()
	for i, n := 0, rs.NumField(); i < n; i++ {
		field := rt.Field(i)
		val := rs.Field(i)
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				continue
			}
			val = val.Elem()
		}
		rd.FieldByName(field.Name).Set(val)

		// Update the liquid variables map for fields that should be exposed to templates
		// Convert field name to lowercase for YAML/liquid compatibility
		key := toLiquidKey(field.Name)
		if key != "" {
			c.Set(key, val.Interface())
		}
	}
}

// toLiquidKey converts a struct field name to its corresponding liquid variable key
func toLiquidKey(fieldName string) string {
	// Map struct field names to their liquid variable names
	// Only map fields that are exposed to liquid templates
	switch fieldName {
	case "BaseURL":
		return "baseurl"
	case "Host":
		return "host"
	case "Port":
		return "port"
	case "Drafts":
		return "show_drafts"
	case "Future":
		return "future"
	case "Unpublished":
		return "unpublished"
	case "Verbose":
		return "verbose"
	case "Incremental":
		return "incremental"
	case "Destination":
		return "destination"
	default:
		return ""
	}
}
