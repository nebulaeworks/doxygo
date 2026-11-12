/**
 * Command doxygo is a Doxygen INPUT_FILTER that lets Doxygen document Go.
 *
 * Doxygen has no Go parser, so doxygo rewrites one Go file into C++-shaped
 * declarations Doxygen's C++ parser understands, keeping every doc comment:
 *
 *   - a package becomes a namespace named after the module's last path
 *     element and the package's directories (e.g. mod::internal::render),
 *     documented by its package comment;
 *   - a struct becomes a struct (embedded fields become base classes), an
 *     interface a class of pure virtual methods, and any other named type
 *     with methods a class; other named types become "using" aliases;
 *   - functions, methods, constants and variables keep their names, with Go
 *     types spelled the C++ way: *T as T*, []T as slice<T>, pkg.T as
 *     pkg::T, map[K]V as map<K, V>, several results as tuple<...>;
 *   - interfaces a type implements become base classes, and imports become
 *     includes, so inheritance, include and directory graphs work.
 *
 * Doc comments are `//` or block comments, with Doxygen commands and
 * Markdown as in C++, and Go doc links ([Name], [pkg.Name.Method]) turned
 * into Doxygen references. Doxygen can't follow calls in Go, so doxygo
 * resolves them with the Go type checker and writes the "References" and
 * "Referenced by" lists and the call and caller graphs itself.
 *
 * Each declaration is written on the line it occupies in the Go source, so
 * Doxygen's line numbers and source links point at the real code.
 *
 * Usage, in a Doxyfile:
 *
 * 	EXTENSION_MAPPING      = go=C++
 * 	FILTER_PATTERNS        = *.go=path/to/doxygo
 * 	REFERENCED_BY_RELATION = NO
 * 	REFERENCES_RELATION    = NO
 * 	CALL_GRAPH             = NO
 * 	CALLER_GRAPH           = NO
 *
 * The last four turn Doxygen's own cross-references off, which would be
 * wrong for Go. The index doxygo builds is cached in the user's cache
 * directory; the DOXYGO_CACHE environment variable names another directory,
 * or "off".
 */
package main

import (
	"fmt"
	"os"
)

/**
 * main filters the file named by the only argument to standard output.
 */
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: doxygo <file.go>")
		os.Exit(2)
	}
	out, err := Filter(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "doxygo:", err)
		os.Exit(1)
	}
	os.Stdout.WriteString(out)
}
