package python_test

import (
	"importcleaner/internal/languageprocessing/python"
	"testing"
)

func Test_PythonImportStatement_AsText(t *testing.T) {
	AssertStatementCorrectlyConvertedToText(t, python.PythonImportStatement{
		ModulePathParts: []string{"somemod"},
		ModuleAlias:     "",
		Dependecies:     []python.PythonDependecy{},
		External:        true,
	}, "import somemod")

	AssertStatementCorrectlyConvertedToText(t, python.PythonImportStatement{
		ModulePathParts: []string{"bigmod"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "something", Alias: ""},
			{Name: "thingwithalias", Alias: "twa"},
			{Name: "anotherthing", Alias: ""},
		},
		External: true,
	}, "from bigmod import something, thingwithalias as twa, anotherthing")

	AssertStatementCorrectlyConvertedToText(t, python.PythonImportStatement{
		ModulePathParts: []string{"mod", "with", "multiple", "parts"},
		ModuleAlias:     "",
		Dependecies: []python.PythonDependecy{
			{Name: "something", Alias: "s"},
			{Name: "thingwithalias", Alias: ""},
			{Name: "anotherthing", Alias: "at"},
		},
		External: true,
	}, "from mod.with.multiple.parts import something as s, thingwithalias, anotherthing as at")
}
