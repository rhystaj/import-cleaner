# Grouping Rules

## All Languages

### `modulePaths`
The paths of the module from which dependencies are being imported. An import will be considered part of the group if it containes *any* of the listed modules

Example usage:
```yaml
modulePaths: ['dataclasses', 'abc']
```

---


### `dependencySourceType`
Possible values are
  - 'internal' - the import defines dependencies which are part of the same codebase
  - 'external' - the import defines dependencies imported externally (e.g. via a package manager such as npm or pip)

Example usage:
```yaml
dependencySourceType: internal
```

---

## Typescript

### `isImportType`
Whether the import is a type import, e.g.
```typescript
import type { SomeType, SomeOtherType } from './module'
```

Example Usage:
```yaml
isImportType: true
```