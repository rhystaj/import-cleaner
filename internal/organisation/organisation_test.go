package organisation_test

import (
	"importcleaner/internal/config"
	"importcleaner/internal/organisation"
	"importcleaner/internal/parsing"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testImportStatement struct {
	details parsing.ImportStatementGenericDetails
}

func (s testImportStatement) AsText() (string, error) {
	panic("not implemented")
}

func (s testImportStatement) GetGenericDetails() parsing.ImportStatementGenericDetails {
	return s.details
}

func generateTestImports() []parsing.ImportStatement {
	return []parsing.ImportStatement{
		testImportStatement{
			details: parsing.ImportStatementGenericDetails{
				ModulePathParts: []string{"helpers", "database"},
				Dependecies:     []parsing.GenericDependecy{},
				External:        false,
			},
		},
		testImportStatement{
			details: parsing.ImportStatementGenericDetails{
				ModulePathParts: []string{"somecoolmodule"},
				Dependecies:     []parsing.GenericDependecy{},
				External:        false,
			},
		},
		testImportStatement{
			details: parsing.ImportStatementGenericDetails{
				ModulePathParts: []string{"dataclasses"},
				Dependecies: []parsing.GenericDependecy{
					{Name: "dataclass", Alias: ""},
					{Name: "field", Alias: ""},
				},
				External: true,
			},
		},
		testImportStatement{
			details: parsing.ImportStatementGenericDetails{
				ModulePathParts: []string{"tools", "saw"},
				Dependecies:     []parsing.GenericDependecy{},
				External:        false,
			},
		},
		testImportStatement{
			details: parsing.ImportStatementGenericDetails{
				ModulePathParts: []string{"pydantic", "types", "numeric"},
				Dependecies:     []parsing.GenericDependecy{},
				External:        true,
			},
		},
		testImportStatement{
			details: parsing.ImportStatementGenericDetails{
				ModulePathParts: []string{"zzz"},
				Dependecies: []parsing.GenericDependecy{
					{Name: "zzzz", Alias: ""},
				},
				External: false,
			},
		},
		testImportStatement{
			details: parsing.ImportStatementGenericDetails{
				ModulePathParts: []string{"pydantic", "validators"},
				Dependecies:     []parsing.GenericDependecy{},
				External:        true,
			},
		},
		testImportStatement{
			details: parsing.ImportStatementGenericDetails{
				ModulePathParts: []string{"helpers", "commandline"},
				Dependecies:     []parsing.GenericDependecy{},
				External:        false,
			},
		},
		testImportStatement{
			details: parsing.ImportStatementGenericDetails{
				ModulePathParts: []string{"anokdependencyiguess"},
				Dependecies: []parsing.GenericDependecy{
					{Name: "thing", Alias: ""},
					{Name: "doodad", Alias: ""},
				},
				External: false,
			},
		},
		testImportStatement{
			details: parsing.ImportStatementGenericDetails{
				ModulePathParts: []string{"pydantic", "types", "string"},
				Dependecies:     []parsing.GenericDependecy{},
				External:        true,
			},
		},
		testImportStatement{
			details: parsing.ImportStatementGenericDetails{
				ModulePathParts: []string{"helpers", "imported"},
				Dependecies:     []parsing.GenericDependecy{},
				External:        true,
			},
		},
		testImportStatement{
			details: parsing.ImportStatementGenericDetails{
				ModulePathParts: []string{"tools", "hammer"},
				Dependecies:     []parsing.GenericDependecy{},
				External:        false,
			},
		},
		testImportStatement{
			details: parsing.ImportStatementGenericDetails{
				ModulePathParts: []string{"abc"},
				Dependecies:     []parsing.GenericDependecy{},
				External:        true,
			},
		},
		testImportStatement{
			details: parsing.ImportStatementGenericDetails{
				ModulePathParts: []string{"argparse"},
				Dependecies:     []parsing.GenericDependecy{},
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

	expectedResult := [][]parsing.ImportStatement{
		{
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"abc"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"anokdependencyiguess"},
					Dependecies: []parsing.GenericDependecy{
						{Name: "thing", Alias: ""},
						{Name: "doodad", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"argparse"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"dataclasses"},
					Dependecies: []parsing.GenericDependecy{
						{Name: "dataclass", Alias: ""},
						{Name: "field", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "commandline"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "database"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "imported"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "numeric"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "string"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "validators"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"somecoolmodule"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "hammer"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "saw"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"zzz"},
					Dependecies: []parsing.GenericDependecy{
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

	expectedResult := [][]parsing.ImportStatement{
		{
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"abc"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"argparse"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"dataclasses"},
					Dependecies: []parsing.GenericDependecy{
						{Name: "dataclass", Alias: ""},
						{Name: "field", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "imported"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "numeric"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "string"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "validators"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
		},
		{
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"anokdependencyiguess"},
					Dependecies: []parsing.GenericDependecy{
						{Name: "thing", Alias: ""},
						{Name: "doodad", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "commandline"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "database"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"somecoolmodule"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "hammer"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "saw"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"zzz"},
					Dependecies: []parsing.GenericDependecy{
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

	expectedResult := [][]parsing.ImportStatement{
		{
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"anokdependencyiguess"},
					Dependecies: []parsing.GenericDependecy{
						{Name: "thing", Alias: ""},
						{Name: "doodad", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "commandline"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "database"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"somecoolmodule"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "hammer"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "saw"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"zzz"},
					Dependecies: []parsing.GenericDependecy{
						{Name: "zzzz", Alias: ""},
					},
					External: false,
				},
			},
		},
		{
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"abc"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"argparse"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"dataclasses"},
					Dependecies: []parsing.GenericDependecy{
						{Name: "dataclass", Alias: ""},
						{Name: "field", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "imported"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "numeric"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "string"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "validators"},
					Dependecies:     []parsing.GenericDependecy{},
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

	expectedResult := [][]parsing.ImportStatement{
		{
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"abc"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"dataclasses"},
					Dependecies: []parsing.GenericDependecy{
						{Name: "dataclass", Alias: ""},
						{Name: "field", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"somecoolmodule"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
		},
		{
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"argparse"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "database"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				//In this group and not in later group because this is the first group that references its parent
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "hammer"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "saw"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
		},
		{
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "numeric"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "string"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
		},
		{
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"anokdependencyiguess"},
					Dependecies: []parsing.GenericDependecy{
						{Name: "thing", Alias: ""},
						{Name: "doodad", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "commandline"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "imported"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "validators"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"zzz"},
					Dependecies: []parsing.GenericDependecy{
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

	expectedResult := [][]parsing.ImportStatement{
		{
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"argparse"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "imported"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
		},
		{
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"dataclasses"},
					Dependecies: []parsing.GenericDependecy{
						{Name: "dataclass", Alias: ""},
						{Name: "field", Alias: ""},
					},
					External: true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "commandline"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"helpers", "database"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"somecoolmodule"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
		},
		{
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"anokdependencyiguess"},
					Dependecies: []parsing.GenericDependecy{
						{Name: "thing", Alias: ""},
						{Name: "doodad", Alias: ""},
					},
					External: false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "hammer"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"tools", "saw"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        false,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"zzz"},
					Dependecies: []parsing.GenericDependecy{
						{Name: "zzzz", Alias: ""},
					},
					External: false,
				},
			},
		},
		{
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"abc"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "numeric"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "types", "string"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
			testImportStatement{
				details: parsing.ImportStatementGenericDetails{
					ModulePathParts: []string{"pydantic", "validators"},
					Dependecies:     []parsing.GenericDependecy{},
					External:        true,
				},
			},
		},
	}

	result := testOrganiser.OrganiseImportStatements(testImports)

	assert.Equal(t, expectedResult, result)

}
