package organisation

import (
	lp "importcleaner/internal/languageprocessing"
	"strings"
)

func CompareImportStatements(a lp.ImportStatement, b lp.ImportStatement) int {
	aDetails := a.GetGenericDetails()
	bDetails := b.GetGenericDetails()

	return strings.Compare(
		strings.Join(aDetails.ModulePathParts, ""),
		strings.Join(bDetails.ModulePathParts, ""),
	)
}
