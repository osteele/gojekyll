package main

import (
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func boolVar(name string, ptr **bool) kingpin.Action {
	return func(pc *kingpin.ParseContext) error {
		value := lookupKingpinValue(name, pc)
		option := false
		if value != nil && *value == "true" {
			option = true
		}
		*ptr = &option
		return nil
	}
}

func stringVar(name string, ptr **string) kingpin.Action {
	return func(pc *kingpin.ParseContext) error {
		value := lookupKingpinValue(name, pc)
		if value != nil {
			copy := *value
			*ptr = &copy
		}
		return nil
	}
}

// From https://github.com/alecthomas/kingpin/issues/184
// Replace w/ kingpin.v3 API when available
func lookupKingpinValue(name string, pc *kingpin.ParseContext) *string {
	for _, el := range pc.Elements {
		switch el.Clause.(type) {
		case *kingpin.ArgClause:
			if (el.Clause).(*kingpin.ArgClause).Model().Name == name {
				return el.Value
			}
		case *kingpin.FlagClause:
			if (el.Clause).(*kingpin.FlagClause).Model().Name == name {
				return el.Value
			}
		}
	}
	def := ""
	return &def
}
