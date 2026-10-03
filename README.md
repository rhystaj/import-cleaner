# Import Cleaner

A configurable tool for automatically cleaning organising imports.
**Still in early development - use at own risk!**

## Features

- Group group imports based on
    - Module Path
    - Whether they are internal/external

**Supported Languages:**
- Python
- Typescript

## Installation

### Build From Source

1) Clone repository
2) From root directory run `go build`

## Use

An *imports-config.yml* is required in the target directory.

To use, run
```bash
<path-to-executable>/importcleaner <target-directory>
```

### Configuration

Specified in *imports-config.yml* file.

Example configuration:
```yaml
language: python
ignoredPaths:
  - ignored_dir
  - submod/ignored_file.py
groupingRules:
  # Imports for external libraries part of the 'dataclasses' OR 'submod' modules
  - modulePaths: ["dataclasses", "submod"]
    dependecySourceType: external
  # Import for internal dependencies
  - dependencySourceType: internal
  ## Any remaining imports (i.e. imports for external dependencies not part of the 'dataclasses' OR 'submod' modules)
```

- `language` determines the target language. Current accepted values are 'python', 'typescript'
- `ignoredPaths` can be used to specify directories and files you don't want import cleaning applied to. An example might be the 'node_modules' folder in typescript or '.env' file in python 
- `groupingRules` specify how imports should be grouped. More detail on use below.

#### Specifying Grouping Rules

`groupingRules` accepts an array of group definition objects. Each entry defines a single group. The order the in which the groups are specifed are the order they appear in files after grouping has been applied

Each group definition defines one or more conditions an import must meet to be considered part of the group. Imports are considered to be part of the *first* group they meet all conditions for, based on the order in which the groups are defined (so it is recommended to define more specific rules first). If an import doesn't meet any conditions for a group it will applied to a 'miscellaneous' group automatically written last.

The rules that can be defined for a group are:

- `modulePaths` - the paths of the module from which dependencies are being imported. An import will be considered part of the group if it containes *any* of the listed modules
- `dependencySourceType` - can be:
  - 'internal' - the import defines dependencies which are part of the same codebase
  - 'external' - the import defines dependencies imported externally (e.g. via a package manager such as npm or pip)

## Future Plans

### Features

In rough precedence:
1) Further grouping rules, including language specific rules (such as grouping type imports in typescript)
2) Supporting further languages
3) Futher cleanup transformations beyond simple grouping (e.g. consolidating imports from the same module)
4) Automatic generation and use of index files

### Workflow
- Implement CI pipeline with unit tests
- Automated tests that run executable against test files in filesystem