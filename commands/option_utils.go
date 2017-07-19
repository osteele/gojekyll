package commands

import (
	"strconv"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func boolVar(name string, ptr **bool) kingpin.Action {
	return func(pc *kingpin.ParseContext) error {
		value := lookupKingpinValue(name, pc)
		if value != nil {
			arg := *value == "true"
			*ptr = &arg
		}
		return nil
	}
}

func intVar(name string, ptr **int) kingpin.Action {
	return func(pc *kingpin.ParseContext) error {
		value := lookupKingpinValue(name, pc)
		if value != nil {
			n, err := strconv.Atoi(*value)
			if err != nil {
				panic(err)
			}
			*ptr = &n
		}
		return nil
	}
}

func stringVar(name string, ptr **string) kingpin.Action {
	return func(pc *kingpin.ParseContext) error {
		value := lookupKingpinValue(name, pc)
		if value != nil {
			arg := *value
			*ptr = &arg
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
