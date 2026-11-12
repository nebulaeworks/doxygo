# doxygo

A Doxygen input filter for Go.

Doxygen has no Go parser. doxygo is an `INPUT_FILTER` that rewrites each Go
file into C++-shaped declarations Doxygen understands, so Go code gets the
same documentation as a C++ project: namespaces, class pages, member lists,
inheritance and collaboration diagrams, include and directory graphs, call
and caller graphs, cross-references and annotated source.

## Installing

```sh
go install github.com/nebulaeworks/doxygo@latest
```

or, from a checkout, `make build` (binary in `bin/`). It needs no
dependencies besides Go's standard library, and Doxygen to use it.

## Using it

In your Doxyfile:

```
EXTENSION_MAPPING      = go=C++
FILE_PATTERNS          = *.go *.md
FILTER_PATTERNS        = *.go=/path/to/doxygo
EXCLUDE_PATTERNS       = *_test.go
JAVADOC_AUTOBRIEF      = YES

# doxygo writes the cross-references and graphs itself:
REFERENCED_BY_RELATION = NO
REFERENCES_RELATION    = NO
CALL_GRAPH             = NO
CALLER_GRAPH           = NO
```

Run `doxygen` as usual: doxygo finds the module (`go.mod`) containing each
file on its own. The rest of a Doxyfile is as for C++ (`HAVE_DOT`,
`SOURCE_BROWSER`, `EXTRACT_PRIVATE` and so on).

## Comments

Either comment style works, and can be mixed:

```go
/**
 * Lookup returns the entity with this RefID, and whether there is one.
 *
 * @param refID The entity's Doxygen id.
 * @return The entity, and whether there is one.
 */
func (x *Index) Lookup(refID string) (*Entity, bool) {

// Lookup returns the entity with this RefID, and whether there is one.
//
// @param refID The entity's Doxygen id.
// @return The entity, and whether there is one.
func (x *Index) Lookup(refID string) (*Entity, bool) {
```

with `/** One line. */` (or `/**< ... */` after a member) for fields.
Doxygen commands and Markdown work as in C++. A file is documented with a
`/** @file ... */` comment before its `package` clause, which Go ignores.

Go doc links are resolved: `[Name]`, `[Type.Method]` and `[pkg.Name]` become
Doxygen references to the module's own symbols, and pkg.go.dev links to
other modules'. A method documented only as "Name implements Iface." gets
Iface's documentation.

## What becomes what

| Go                                | Doxygen                                           |
|-----------------------------------|---------------------------------------------------|
| package                           | namespace: module's last element, then directories |
| struct, embedded field            | struct, base class                                |
| interface                         | class of pure virtual methods                     |
| type implementing an interface    | class with the interface as a base class          |
| other named type with methods     | class, noting its Go type                         |
| other named type                  | `using` alias                                     |
| `*T`, `[]T`, `map[K]V`, `func`    | `T*`, `slice<T>`, `map<K, V>`, `func<R(P)>`       |
| several results                   | `tuple<...>`                                      |
| import                            | include (of the package's file, within the module) |
| unexported name                   | private                                           |

Every declaration stays on its source line, so "Definition at line N" and
the source browser point at the real code.

## References and graphs

Doxygen can't follow calls in Go, so doxygo resolves them with the Go type
checker (`go/types`, over the whole module) and adds, to each entity's
documentation, "References" and "Referenced by" lists and, for functions,
collapsible call and caller graphs. These replace Doxygen's own, which is
why those settings are turned off.

Type-checking the module for every file would be slow, so the result is
cached in the user cache directory, keyed by the module's source. Set
`DOXYGO_CACHE=off` to disable the cache, or to a directory to use it
instead.

## Limits

- Calls through interface values or function values aren't followed:
  doxygo links to the interface method, not the implementations.
- A package that doesn't compile is checked as far as it can be; references
  involving the broken parts are missing.
- Generic types are shown as C++ templates; instantiations aren't tracked.
- Build tags: files that don't build on the machine running Doxygen are
  skipped.

## Development

```sh
make test   # go vet + go test
make docs   # doxygo's own documentation, in build/docs/html
```

doxygo documents itself, and its docs build fails on anything undocumented,
parameters and results included.

## License

[MPL-2.0](LICENSE).
