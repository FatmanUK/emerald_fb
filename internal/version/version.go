// Package version holds the two build facts a MUCK reports about
// itself, in one place because four surfaces read the first of them
// and had drifted apart.
//
// They are upstream's `VERSION` (config.h:24) and `compile_options`
// (game.c:66). Emerald's *own* release is `game.Version`, set at link
// time, and is a different thing: it says which build this is, where
// these two say what a program may expect of it.
package version

// Muck is upstream's VERSION (config.h:24), reported verbatim.
//
// Four surfaces return it and upstream returns the same constant from
// every one: MUF VERSION (prim_version, p_stack.c:817), MPI {version}
// (mfn_version, mfuns.c:129), the compiler's __version define
// (compile.c:707), and @version's first line (do_version,
// help.c:587). A program asks this to decide what it may use, so the
// answer has to be the compatibility version and not a product name
// -- "Fuzzball Emerald dev" told a program nothing it could act on,
// and told it in a different shape from what MPI told it.
//
// It is pinned rather than derived: bumping the fuzzball submodule to
// a release whose config.h says something else should be a deliberate
// edit here, because every world's MUF has been written against
// whatever this says.
const Muck = "Muck2.2fb7.21"

// CompileOptions is upstream's compile_options (game.c:66): the
// optional subsystems this build has, space-separated, with a
// trailing space on every token.
//
// Three of upstream's nine apply here. GODPRIV because `#1` is
// MLEV_GOD at 255 and TUNE_MLEV reserves parameters to it; MCPGUI
// because internal/mcp draws dialogs; SSL because Emerald is
// TLS-only, so the one upstream makes optional is the one that can
// never be absent. The other six are not Emerald's: DISKBASE is
// replaced by Postgres, MEMPROF is malloc profiling, MEMCLEANUP and
// DEBUG are C build modes, IPFWD is upstream's forwarding host, and
// RESOLVER is its *separate resolver process* where Emerald resolves
// names on its own goroutines -- which is upstream's `#ifndef
// SPAWN_HOST_RESOLVER` shape, not a missing feature.
//
// The order is upstream's source order, which is the order of its
// `#ifdef` blocks rather than alphabetical.
//
// Read it with tune.Param.Active, which is MOD_ENABLED (game.h:30) --
// a *substring* test, so the token "MCP" resolves against "MCPGUI ".
const CompileOptions = "GODPRIV MCPGUI SSL "
