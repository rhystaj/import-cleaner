package organisation_test

import (
	"importcleaner/internal/config"
	lp "importcleaner/internal/languageprocessing"
	"importcleaner/internal/organisation"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testImportStatement struct {
	details lp.ImportStatementGenericDetails
}

func (s testImportStatement) AsText() (string, error) {
	panic("not implemented")
}

func (s testImportStatement) GetGenericDetails() lp.ImportStatementGenericDetails {
	return s.details
}

func (s testImportStatement) ConsolidationKey() string {
	return strings.Join(s.details.ModulePathParts, ".") + "|" + s.details.ModuleAlias
}

func (s testImportStatement) ConsolidateWith(other lp.ImportStatement) (lp.ImportStatement, error) {
	return testImportStatement{
		details: lp.ImportStatementGenericDetails{
			ModulePathParts: s.details.ModulePathParts,
			ModuleAlias:     s.details.ModuleAlias,
			Dependecies:     append(s.details.Dependecies, other.GetGenericDetails().Dependecies...),
			External:        s.details.External,
		},
	}, nil
}

func generateTestImports() []lp.ImportStatement {
	return []lp.ImportStatement{
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"helpers", "database"},
				Dependecies:     []lp.GenericDependecy{},
				External:        false,
			},
		},
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"somecoolmodule"},
				Dependecies:     []lp.GenericDependecy{},
				External:        false,
			},
		},
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"dataclasses"},
				Dependecies: []lp.GenericDependecy{
					{Name: "dataclass", Alias: ""},
					{Name: "field", Alias: ""},
				},
				External: true,
			},
		},
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"tools"},
				Dependecies: []lp.GenericDependecy{
					{Name: "Saw", Alias: ""},
				},
				External: false,
			},
		},
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"pydantic", "types"},
				Dependecies: []lp.GenericDependecy{
					{Name: "Boolean", Alias: ""},
					{Name: "Float", Alias: ""},
				},
				External: true,
			},
		},
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"tools", "power"},
				Dependecies: []lp.GenericDependecy{
					{Name: "Drill", Alias: ""},
				},
				External: false,
			},
		},
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"zzz"},
				Dependecies: []lp.GenericDependecy{
					{Name: "zzzz", Alias: ""},
				},
				External: false,
			},
		},
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"pydantic", "validators"},
				Dependecies:     []lp.GenericDependecy{},
				External:        true,
			},
		},
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"helpers", "commandline"},
				Dependecies:     []lp.GenericDependecy{},
				External:        false,
			},
		},
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"anokdependencyiguess"},
				Dependecies: []lp.GenericDependecy{
					{Name: "thing", Alias: ""},
					{Name: "doodad", Alias: ""},
				},
				External: false,
			},
		},
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"pydantic", "types"},
				Dependecies: []lp.GenericDependecy{
					{Name: "String", Alias: ""},
					{Name: "Integer", Alias: ""},
				},
				External: true,
			},
		},
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"helpers", "imported"},
				Dependecies:     []lp.GenericDependecy{},
				External:        true,
			},
		},
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"tools"},
				Dependecies: []lp.GenericDependecy{
					{Name: "Hammer", Alias: ""},
					{Name: "Screwdriver", Alias: "sd"},
				},
				External: false,
			},
		},
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"abc"},
				Dependecies:     []lp.GenericDependecy{},
				External:        true,
			},
		},
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"argparse"},
				Dependecies:     []lp.GenericDependecy{},
				External:        true,
			},
		},
	}
}

func TestOrganiseImports_NoGroupsConfigured_NoConsolidation(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl[lp.ImportStatement]{
		GroupRules: []organisation.ImportStatementGroupRule[lp.ImportStatement]{},
	}
	testImports := generateTestImports()

	expectedResult := [][]lp.ImportStatement{
		{
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"abc"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"anokdependencyiguess"},
					Dependecies: []lp.GenericDependecy{
						{Name: "thing", Alias: ""},
						{Name: "doodad", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"argparse"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"dataclasses"},
					Dependecies: []lp.GenericDependecy{
						{Name: "dataclass", Alias: ""},
						{Name: "field", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "commandline"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "database"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "imported"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Boolean", Alias: ""},
						{Name: "Float", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types"},
					Dependecies: []lp.GenericDependecy{
						{Name: "String", Alias: ""},
						{Name: "Integer", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "validators"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"somecoolmodule"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Saw", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Hammer", Alias: ""},
						{Name: "Screwdriver", Alias: "sd"},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "power"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Drill", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"zzz"},
					Dependecies: []lp.GenericDependecy{
						{Name: "zzzz", Alias: ""},
					},
					External: false,
				},
			},
		},
	}

	result, err := testOrganiser.OrganiseImportStatements(testImports)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestOrganiseImports_NoGroupsConfigured_Consolidation(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl[lp.ImportStatement]{
		GroupRules:                  []organisation.ImportStatementGroupRule[lp.ImportStatement]{},
		ConsolidateImportStatements: true,
	}
	testImports := generateTestImports()

	expectedResult := [][]lp.ImportStatement{
		{
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"abc"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"anokdependencyiguess"},
					Dependecies: []lp.GenericDependecy{
						{Name: "thing", Alias: ""},
						{Name: "doodad", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"argparse"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"dataclasses"},
					Dependecies: []lp.GenericDependecy{
						{Name: "dataclass", Alias: ""},
						{Name: "field", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "commandline"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "database"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "imported"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Boolean", Alias: ""},
						{Name: "Float", Alias: ""},
						{Name: "String", Alias: ""},
						{Name: "Integer", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "validators"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"somecoolmodule"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Saw", Alias: ""},
						{Name: "Hammer", Alias: ""},
						{Name: "Screwdriver", Alias: "sd"},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "power"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Drill", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"zzz"},
					Dependecies: []lp.GenericDependecy{
						{Name: "zzzz", Alias: ""},
					},
					External: false,
				},
			},
		},
	}

	result, err := testOrganiser.OrganiseImportStatements(testImports)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestOrganiseImports_ExternalGroup(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl[lp.ImportStatement]{
		GroupRules: []organisation.ImportStatementGroupRule[lp.ImportStatement]{
			organisation.NewGenericGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				DependencySourceType: config.DependencySourceTypeExternal,
			}),
		},
	}
	testImports := generateTestImports()

	expectedResult := [][]lp.ImportStatement{
		{
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"abc"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"argparse"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"dataclasses"},
					Dependecies: []lp.GenericDependecy{
						{Name: "dataclass", Alias: ""},
						{Name: "field", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "imported"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Boolean", Alias: ""},
						{Name: "Float", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types"},
					Dependecies: []lp.GenericDependecy{
						{Name: "String", Alias: ""},
						{Name: "Integer", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "validators"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
		},
		{
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"anokdependencyiguess"},
					Dependecies: []lp.GenericDependecy{
						{Name: "thing", Alias: ""},
						{Name: "doodad", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "commandline"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "database"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"somecoolmodule"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Saw", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Hammer", Alias: ""},
						{Name: "Screwdriver", Alias: "sd"},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "power"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Drill", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"zzz"},
					Dependecies: []lp.GenericDependecy{
						{Name: "zzzz", Alias: ""},
					},
					External: false,
				},
			},
		},
	}

	result, err := testOrganiser.OrganiseImportStatements(testImports)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestOrganiseImports_InternalGroup(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl[lp.ImportStatement]{
		GroupRules: []organisation.ImportStatementGroupRule[lp.ImportStatement]{
			organisation.NewGenericGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				DependencySourceType: config.DependencySourceTypeInternal,
			}),
		},
	}

	testImports := generateTestImports()

	expectedResult := [][]lp.ImportStatement{
		{
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"anokdependencyiguess"},
					Dependecies: []lp.GenericDependecy{
						{Name: "thing", Alias: ""},
						{Name: "doodad", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "commandline"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "database"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"somecoolmodule"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Saw", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Hammer", Alias: ""},
						{Name: "Screwdriver", Alias: "sd"},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "power"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Drill", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"zzz"},
					Dependecies: []lp.GenericDependecy{
						{Name: "zzzz", Alias: ""},
					},
					External: false,
				},
			},
		},
		{
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"abc"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"argparse"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"dataclasses"},
					Dependecies: []lp.GenericDependecy{
						{Name: "dataclass", Alias: ""},
						{Name: "field", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "imported"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Boolean", Alias: ""},
						{Name: "Float", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types"},
					Dependecies: []lp.GenericDependecy{
						{Name: "String", Alias: ""},
						{Name: "Integer", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "validators"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
		},
	}

	result, err := testOrganiser.OrganiseImportStatements(testImports)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestOrganiseImports_GroupByModulePaths(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl[lp.ImportStatement]{
		GroupRules: []organisation.ImportStatementGroupRule[lp.ImportStatement]{
			organisation.NewGenericGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				ModulePaths: [][]string{{"dataclasses"}, {"somecoolmodule"}, {"abc"}},
			}),
			organisation.NewGenericGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				ModulePaths: [][]string{{"tools"}, {"helpers", "database"}, {"argparse"}},
			}),
			organisation.NewGenericGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				ModulePaths: [][]string{{"tools", "power"}, {"pydantic", "types"}},
			}),
		},
	}

	testImports := generateTestImports()

	expectedResult := [][]lp.ImportStatement{
		{
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"abc"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"dataclasses"},
					Dependecies: []lp.GenericDependecy{
						{Name: "dataclass", Alias: ""},
						{Name: "field", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"somecoolmodule"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
		},
		{
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"argparse"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "database"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Saw", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Hammer", Alias: ""},
						{Name: "Screwdriver", Alias: "sd"},
					},
					External: false,
				},
			},
			testImportStatement{
				//In this group and not in later group because this is the first group that references its parent
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "power"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Drill", Alias: ""},
					},
					External: false,
				},
			},
		},
		{
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Boolean", Alias: ""},
						{Name: "Float", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types"},
					Dependecies: []lp.GenericDependecy{
						{Name: "String", Alias: ""},
						{Name: "Integer", Alias: ""},
					},
					External: true,
				},
			},
		},
		{
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"anokdependencyiguess"},
					Dependecies: []lp.GenericDependecy{
						{Name: "thing", Alias: ""},
						{Name: "doodad", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "commandline"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "imported"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "validators"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"zzz"},
					Dependecies: []lp.GenericDependecy{
						{Name: "zzzz", Alias: ""},
					},
					External: false,
				},
			},
		},
	}

	result, err := testOrganiser.OrganiseImportStatements(testImports)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestOrganiseImports_CompositeGroups(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl[lp.ImportStatement]{
		GroupRules: []organisation.ImportStatementGroupRule[lp.ImportStatement]{
			organisation.NewGenericGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				DependencySourceType: config.DependencySourceTypeExternal,
				ModulePaths:          [][]string{{"helpers"}, {"anokdependencyiguess"}, {"argparse"}},
			}),
			organisation.NewGenericGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				ModulePaths: [][]string{{"helpers"}, {"dataclasses"}, {"somecoolmodule"}},
			}),
			organisation.NewGenericGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				DependencySourceType: config.DependencySourceTypeInternal,
			}),
		},
	}

	testImports := generateTestImports()

	expectedResult := [][]lp.ImportStatement{
		{
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"argparse"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "imported"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
		},
		{
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"dataclasses"},
					Dependecies: []lp.GenericDependecy{
						{Name: "dataclass", Alias: ""},
						{Name: "field", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "commandline"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "database"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"somecoolmodule"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
		},
		{
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"anokdependencyiguess"},
					Dependecies: []lp.GenericDependecy{
						{Name: "thing", Alias: ""},
						{Name: "doodad", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Saw", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Hammer", Alias: ""},
						{Name: "Screwdriver", Alias: "sd"},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "power"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Drill", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"zzz"},
					Dependecies: []lp.GenericDependecy{
						{Name: "zzzz", Alias: ""},
					},
					External: false,
				},
			},
		},
		{
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"abc"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types"},
					Dependecies: []lp.GenericDependecy{
						{Name: "Boolean", Alias: ""},
						{Name: "Float", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types"},
					Dependecies: []lp.GenericDependecy{
						{Name: "String", Alias: ""},
						{Name: "Integer", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "validators"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
		},
	}

	result, err := testOrganiser.OrganiseImportStatements(testImports)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}
