package organisation_test

import (
	"importcleaner/internal/config"
	lp "importcleaner/internal/languageprocessing"
	"importcleaner/internal/organisation"
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
				ModulePathParts: []string{"tools", "saw"},
				Dependecies:     []lp.GenericDependecy{},
				External:        false,
			},
		},
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"pydantic", "types", "numeric"},
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
				ModulePathParts: []string{"pydantic", "types", "string"},
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
		testImportStatement{
			details: lp.ImportStatementGenericDetails{
				ModulePathParts: []string{"tools", "hammer"},
				Dependecies:     []lp.GenericDependecy{},
				External:        false,
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

func TestOrganiseImports_NoGroupsConfigured(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: []organisation.ImportStatementGroupRule{},
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
					ModulePathParts: []string{"pydantic", "types", "numeric"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "string"},
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
					ModulePathParts: []string{"somecoolmodule"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "hammer"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "saw"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
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

	result := testOrganiser.OrganiseImportStatements(testImports)

	assert.Equal(t, expectedResult, result)
}

func TestOrganiseImports_ExternalGroup(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: []organisation.ImportStatementGroupRule{
			organisation.NewGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
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
					ModulePathParts: []string{"pydantic", "types", "numeric"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "string"},
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
					ModulePathParts: []string{"tools", "hammer"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "saw"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
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

	result := testOrganiser.OrganiseImportStatements(testImports)

	assert.Equal(t, expectedResult, result)
}

func TestOrganiseImports_InternalGroup(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: []organisation.ImportStatementGroupRule{
			organisation.NewGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
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
					ModulePathParts: []string{"tools", "hammer"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "saw"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
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
					ModulePathParts: []string{"pydantic", "types", "numeric"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "string"},
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
		},
	}

	result := testOrganiser.OrganiseImportStatements(testImports)

	assert.Equal(t, expectedResult, result)
}

func TestOrganiseImports_GroupByModulePaths(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: []organisation.ImportStatementGroupRule{
			organisation.NewGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				ModulePaths: [][]string{{"dataclasses"}, {"somecoolmodule"}, {"abc"}},
			}),
			organisation.NewGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				ModulePaths: [][]string{{"tools"}, {"helpers", "database"}, {"argparse"}},
			}),
			organisation.NewGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				ModulePaths: [][]string{{"tools", "hammer"}, {"pydantic", "types"}},
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
				//In this group and not in later group because this is the first group that references its parent
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "hammer"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "saw"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
		},
		{
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "numeric"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "string"},
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

	result := testOrganiser.OrganiseImportStatements(testImports)

	assert.Equal(t, expectedResult, result)
}

func TestOrganiseImports_CompositeGroups(t *testing.T) {

	testOrganiser := organisation.ImportStatementOrganiserImpl{
		GroupRules: []organisation.ImportStatementGroupRule{
			organisation.NewGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				DependencySourceType: config.DependencySourceTypeExternal,
				ModulePaths:          [][]string{{"helpers"}, {"anokdependencyiguess"}, {"argparse"}},
			}),
			organisation.NewGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
				ModulePaths: [][]string{{"helpers"}, {"dataclasses"}, {"somecoolmodule"}},
			}),
			organisation.NewGroupingRuleFromDefinition(config.ImportStatementGroupRuleDefinition{
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
					ModulePathParts: []string{"tools", "hammer"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "saw"},
					Dependecies:     []lp.GenericDependecy{},
					External:        false,
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
					ModulePathParts: []string{"pydantic", "types", "numeric"},
					Dependecies:     []lp.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: lp.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "string"},
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
		},
	}

	result := testOrganiser.OrganiseImportStatements(testImports)

	assert.Equal(t, expectedResult, result)

}
