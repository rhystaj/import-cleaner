package organisation_test

import (
	"importcleaner/internal/organisation"
	"importcleaner/internal/parsing"
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateTestImports() []parsing.ImportStatement {
	return []parsing.ImportStatement{
		{
			ModulePathParts: []string{"helpers", "database"},
			Dependecies:     []parsing.Dependecy{},
			External:        false,
		},
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
			ModulePathParts: []string{"tools", "saw"},
			Dependecies:     []parsing.Dependecy{},
			External:        false,
		},
		{
			ModulePathParts: []string{"pydantic", "types", "numeric"},
			Dependecies:     []parsing.Dependecy{},
			External:        true,
		},
		{
			ModulePathParts: []string{"zzz"},
			Dependecies: []parsing.Dependecy{
				{Name: "zzzz", Alias: ""},
			},
			External: false,
		},
		{
			ModulePathParts: []string{"pydantic", "validators"},
			Dependecies:     []parsing.Dependecy{},
			External:        true,
		},
		{
			ModulePathParts: []string{"helpers", "commandline"},
			Dependecies:     []parsing.Dependecy{},
			External:        false,
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
			ModulePathParts: []string{"pydantic", "types", "string"},
			Dependecies:     []parsing.Dependecy{},
			External:        true,
		},
		{
			ModulePathParts: []string{"helpers", "imported"},
			Dependecies:     []parsing.Dependecy{},
			External:        true,
		},
		{
			ModulePathParts: []string{"tools", "hammer"},
			Dependecies:     []parsing.Dependecy{},
			External:        false,
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
	testImports := generateTestImports()

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
				ModulePathParts: []string{"helpers", "commandline"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"helpers", "database"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"helpers", "imported"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"pydantic", "types", "numeric"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"pydantic", "types", "string"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"pydantic", "validators"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"somecoolmodule"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"tools", "hammer"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"tools", "saw"},
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

	assert.Equal(t, expectedResult, result)
}

func TestOrganiseImports_ExternalGroup(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: []organisation.ImportStatementGroupRule{
			{
				DependencySourceType: organisation.DependencySourceTypeExternal,
			},
		},
	}
	testImports := generateTestImports()

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
			{
				ModulePathParts: []string{"helpers", "imported"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"pydantic", "types", "numeric"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"pydantic", "types", "string"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"pydantic", "validators"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
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
				ModulePathParts: []string{"helpers", "commandline"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"helpers", "database"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"somecoolmodule"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"tools", "hammer"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"tools", "saw"},
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

	assert.Equal(t, expectedResult, result)
}

func TestOrganiseImports_InternalGroup(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: []organisation.ImportStatementGroupRule{
			{
				DependencySourceType: organisation.DependencySourceTypeInternal,
			},
		},
	}

	testImports := generateTestImports()

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
				ModulePathParts: []string{"helpers", "commandline"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"helpers", "database"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"somecoolmodule"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"tools", "hammer"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"tools", "saw"},
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
			{
				ModulePathParts: []string{"helpers", "imported"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"pydantic", "types", "numeric"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"pydantic", "types", "string"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"pydantic", "validators"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
		},
	}

	result := testOrganiser.OrganiseImportStatements(testImports)

	assert.Equal(t, expectedResult, result)
}

func TestOrganiseImports_GroupByModulePaths(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: []organisation.ImportStatementGroupRule{
			{
				ModulePaths: [][]string{{"dataclasses"}, {"somecoolmodule"}, {"abc"}},
			},
			{
				ModulePaths: [][]string{{"tools"}, {"helpers", "database"}, {"argparse"}},
			},
			{
				ModulePaths: [][]string{{"tools", "hammer"}, {"pydantic", "types"}},
			},
		},
	}

	testImports := generateTestImports()

	expectedResult := [][]parsing.ImportStatement{
		{
			{
				ModulePathParts: []string{"abc"},
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
		},
		{
			{
				ModulePathParts: []string{"argparse"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"helpers", "database"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				//In this group and not in later group because this is the first group that references its parent
				ModulePathParts: []string{"tools", "hammer"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"tools", "saw"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
		},
		{
			{
				ModulePathParts: []string{"pydantic", "types", "numeric"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"pydantic", "types", "string"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
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
				ModulePathParts: []string{"helpers", "commandline"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"helpers", "imported"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"pydantic", "validators"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
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

	assert.Equal(t, expectedResult, result)
}

func TestOrganiseImports_CompositeGroups(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: []organisation.ImportStatementGroupRule{
			{
				ModulePaths:          [][]string{{"helpers"}, {"anokdependencyiguess"}, {"argparse"}},
				DependencySourceType: organisation.DependencySourceTypeExternal,
			},
			{
				ModulePaths: [][]string{{"helpers"}, {"dataclasses"}, {"somecoolmodule"}},
			},
			{
				DependencySourceType: organisation.DependencySourceTypeInternal,
			},
		},
	}

	testImports := generateTestImports()

	expectedResult := [][]parsing.ImportStatement{
		{
			{
				ModulePathParts: []string{"argparse"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"helpers", "imported"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
		},
		{
			{
				ModulePathParts: []string{"dataclasses"},
				Dependecies: []parsing.Dependecy{
					{Name: "dataclass", Alias: ""},
					{Name: "field", Alias: ""},
				},
				External: true,
			},
			{
				ModulePathParts: []string{"helpers", "commandline"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"helpers", "database"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"somecoolmodule"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
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
				ModulePathParts: []string{"tools", "hammer"},
				Dependecies:     []parsing.Dependecy{},
				External:        false,
			},
			{
				ModulePathParts: []string{"tools", "saw"},
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
				ModulePathParts: []string{"pydantic", "types", "numeric"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"pydantic", "types", "string"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
			{
				ModulePathParts: []string{"pydantic", "validators"},
				Dependecies:     []parsing.Dependecy{},
				External:        true,
			},
		},
	}

	result := testOrganiser.OrganiseImportStatements(testImports)

	assert.Equal(t, expectedResult, result)

}
