package organisation_test

import (
	"importcleaner/internal/organisation"
	"importcleaner/internal/parsing"
	"reflect"
	"testing"
)

func TestOrganiseImports(t *testing.T) {

	testImports := []parsing.ImportStatement{
		{
			ModuleName:  "somecoolmodule",
			Dependecies: []parsing.Dependecy{},
			External:    false,
		},
		{
			ModuleName: "dataclasses",
			Dependecies: []parsing.Dependecy{
				{Name: "dataclass", Alias: ""},
				{Name: "field", Alias: ""},
			},
			External: true,
		},
		{
			ModuleName: "zzz",
			Dependecies: []parsing.Dependecy{
				{Name: "zzzz", Alias: ""},
			},
			External: false,
		},
		{
			ModuleName: "anokdependencyiguess",
			Dependecies: []parsing.Dependecy{
				{Name: "thing", Alias: ""},
				{Name: "doodad", Alias: ""},
			},
			External: false,
		},
		{
			ModuleName:  "abc",
			Dependecies: []parsing.Dependecy{},
			External:    true,
		},
		{
			ModuleName:  "argparse",
			Dependecies: []parsing.Dependecy{},
			External:    true,
		},
	}

	expectedResult := [][]parsing.ImportStatement{
		{
			{
				ModuleName:  "abc",
				Dependecies: []parsing.Dependecy{},
				External:    true,
			},
			{
				ModuleName:  "argparse",
				Dependecies: []parsing.Dependecy{},
				External:    true,
			},
			{
				ModuleName: "dataclasses",
				Dependecies: []parsing.Dependecy{
					{Name: "dataclass", Alias: ""},
					{Name: "field", Alias: ""},
				},
				External: true,
			},
		},
		{
			{
				ModuleName: "anokdependencyiguess",
				Dependecies: []parsing.Dependecy{
					{Name: "thing", Alias: ""},
					{Name: "doodad", Alias: ""},
				},
				External: false,
			},
			{
				ModuleName:  "somecoolmodule",
				Dependecies: []parsing.Dependecy{},
				External:    false,
			},
			{
				ModuleName: "zzz",
				Dependecies: []parsing.Dependecy{
					{Name: "zzzz", Alias: ""},
				},
				External: false,
			},
		},
	}

	result := organisation.OrganiseImportStatements(testImports)

	if !reflect.DeepEqual(expectedResult, result) {
		t.Errorf("Expected %+v, by recieved %+v\n", expectedResult, result)
	}
}
