package organisation_test

import (
	"importcleaner/internal/organisation"
	"importcleaner/internal/parsing"
	"reflect"
	"testing"
)

func generateTextImports() []parsing.ImportStatement {
	return []parsing.ImportStatement{
		{
			ModulePathParts: []string{"somecoolmodule"},
			Dependecies:     []parsing.Dependecy{},
			External:        false,
		},
		{
			ModulePathParts: []string{"dataclasses"},
			Dependecies: []parsing.Dependecy{
				{Name: "dataclass", Alias: ""},
				{Name: "field", Alias: ""},
			},
			External: true,
		},
		{
			ModulePathParts: []string{"zzz"},
			Dependecies: []parsing.Dependecy{
				{Name: "zzzz", Alias: ""},
			},
			External: false,
		},
		{
			ModulePathParts: []string{"anokdependencyiguess"},
			Dependecies: []parsing.Dependecy{
				{Name: "thing", Alias: ""},
				{Name: "doodad", Alias: ""},
			},
			External: false,
		},
		{
			ModulePathParts: []string{"abc"},
			Dependecies:     []parsing.Dependecy{},
			External:        true,
		},
		{
			ModulePathParts: []string{"argparse"},
			Dependecies:     []parsing.Dependecy{},
			External:        true,
		},
	}
}

func TestOrganiseImports_NoGroupsConfigured(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: []organisation.ImportStatementGroupRule{},
	}
	testImports := generateTextImports()

	expectedResult := [][]parsing.ImportStatement{
		{
			{
				ModulePathParts: []string{"abc"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"anokdependencyiguess"},
				Dependecies: []parsing.Dependecy{
					{Name: "thing", Alias: ""},
					{Name: "doodad", Alias: ""},
				},
				External: false,
			},
			{
				ModulePathParts: []string{"argparse"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"dataclasses"},
				Dependecies: []parsing.Dependecy{
					{Name: "dataclass", Alias: ""},
					{Name: "field", Alias: ""},
				},
				External: true,
			},
			{
				ModulePathParts: []string{"somecoolmodule"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"zzz"},
				Dependecies: []parsing.Dependecy{
					{Name: "zzzz", Alias: ""},
				},
				External: false,
			},
		},
	}

	result := testOrganiser.OrganiseImportStatements(testImports)

	if !reflect.DeepEqual(expectedResult, result) {
		t.Errorf("Expected %+v, by recieved %+v\n", expectedResult, result)
	}
}

func TestOrganiseImports_ExternalGroup(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: []organisation.ImportStatementGroupRule{
			{
				DependencySourceType: organisation.DependencySourceTypeExternal,
			},
		},
	}
	testImports := generateTextImports()

	expectedResult := [][]parsing.ImportStatement{
		{
			{
				ModulePathParts: []string{"abc"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"argparse"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"dataclasses"},
				Dependecies: []parsing.Dependecy{
					{Name: "dataclass", Alias: ""},
					{Name: "field", Alias: ""},
				},
				External: true,
			},
		},
		{
			{
				ModulePathParts: []string{"anokdependencyiguess"},
				Dependecies: []parsing.Dependecy{
					{Name: "thing", Alias: ""},
					{Name: "doodad", Alias: ""},
				},
				External: false,
			},
			{
				ModulePathParts: []string{"somecoolmodule"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"zzz"},
				Dependecies: []parsing.Dependecy{
					{Name: "zzzz", Alias: ""},
				},
				External: false,
			},
		},
	}

	result := testOrganiser.OrganiseImportStatements(testImports)

	if !reflect.DeepEqual(expectedResult, result) {
		t.Errorf("Expected %+v, by recieved %+v\n", expectedResult, result)
	}
}

func TestOrganiseImports_InternalGroup(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: []organisation.ImportStatementGroupRule{
			{
				DependencySourceType: organisation.DependencySourceTypeInternal,
			},
		},
	}

	testImports := generateTextImports()

	expectedResult := [][]parsing.ImportStatement{
		{
			{
				ModulePathParts: []string{"anokdependencyiguess"},
				Dependecies: []parsing.Dependecy{
					{Name: "thing", Alias: ""},
					{Name: "doodad", Alias: ""},
				},
				External: false,
			},
			{
				ModulePathParts: []string{"somecoolmodule"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"zzz"},
				Dependecies: []parsing.Dependecy{
					{Name: "zzzz", Alias: ""},
				},
				External: false,
			},
		},
		{
			{
				ModulePathParts: []string{"abc"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"argparse"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"dataclasses"},
				Dependecies: []parsing.Dependecy{
					{Name: "dataclass", Alias: ""},
					{Name: "field", Alias: ""},
				},
				External: true,
			},
		},
	}

	result := testOrganiser.OrganiseImportStatements(testImports)

	if !reflect.DeepEqual(expectedResult, result) {
		t.Errorf("Expected %+v, by recieved %+v\n", expectedResult, result)
	}
}
