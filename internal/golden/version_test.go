package golden

import (
	"context"
	"regexp"
	"strings"
	"testing"
)

// Nothing compared any version surface before this file, which is how
// they came to disagree. Upstream answers one question with one
// constant from four places -- MUF VERSION (prim_version,
// p_stack.c:817), MPI {version} (mfn_version, mfuns.c:129), the
// compiler's __version define (compile.c:707) and @version's first
// line (do_version, help.c:587), each returning VERSION from
// config.h:24. Emerald answered it three different ways: the product
// name and the build version from MUF, upstream's constant from MPI,
// and the build version from the define.
//
// The MUF, MPI and define answers compare **verbatim**, with no mask
// at all, which is the strong half of this case.
//
// No parentheses inside the MUF comments: ")" ends a "( ... )"
// comment wherever it appears.
const versionSource = `$ifdef __version>Muck2.2fb
$define vnewer "yes" $enddef
$else
$define vnewer "no" $enddef
$endif
$ifdef __version<Muck2.3
$define volder "yes" $enddef
$else
$define volder "no" $enddef
$endif
: t[ s -- ] me @ s @ notify ;
: main
  ( What prim_version pushes. )
  version "MUF " swap strcat t

  ( The same constant reached through MPI. These two had no
    reason to agree and did not. parsempi takes four: the object,
    the MPI text, the argument string and the private flag. )
  me @ "{version}" "" 1 parsempi "MPI " swap strcat t

  ( And the compiler define, exercised the one way the manual
    uses it -- a textual comparison in a directive, not a value.
    Expanding it as a value is a compile error on both servers,
    because insert_def installs it unquoted. )
  vnewer "NEWER " swap strcat t
  volder "OLDER " swap strcat t
;`

var versionScript = Script{
	// Two lines upstream, three here: the third names this
	// server, which upstream has no equivalent of and which
	// maskOptions drops.
	"@version",

	// The three in-game surfaces, compared as they come.
	"test",
}

var (
	// optionsLine is do_version's second line. Its value cannot
	// match and the difference is the harness's and the port's,
	// not a fault: the oracle is built --without-ssl
	// (deploy/golden/Containerfile.fbmuck), so it has no SSL
	// token, and its config.h defines SPAWN_HOST_RESOLVER on
	// Linux, so it has RESOLVER -- where Emerald is TLS-only and
	// resolves names on its own goroutines. The line and its
	// position are kept; only the value is masked.
	optionsLine = regexp.MustCompile(`^Options: .*$`)
	// emeraldLine is the third line, which is this server's own
	// and deliberately absent upstream.
	emeraldLine = regexp.MustCompile(`^Emerald FB .+$`)
)

// maskOptions blanks the compile-options value and drops the line
// upstream does not print. It leaves "Version:" alone, because that
// line agreeing verbatim is the point of the whole case.
func maskOptions(s string) string {
	var kept []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSuffix(line, "\r")
		if emeraldLine.MatchString(line) {
			continue
		}
		if optionsLine.MatchString(line) {
			line = "Options: <build>"
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n")
}

// TestVersionMatchesFuzzball compares the four surfaces.
func TestVersionMatchesFuzzball(t *testing.T) {
	requireOracle(t)
	ctx := context.Background()

	fx, err := WriteFixture(t.TempDir(), versionSource)
	if err != nil {
		t.Fatal(err)
	}
	script := versionScript
	oracle, err := RunOracleSteps(ctx, fx, script, nil)
	if err != nil {
		t.Fatalf("driving the oracle: %v", err)
	}
	emerald, err := RunEmeraldSteps(ctx, fx, script, nil)
	if err != nil {
		t.Fatalf("driving this server: %v", err)
	}
	for i, cmd := range script {
		var want, got string
		if i < len(oracle) {
			want = oracle[i]
		}
		if i < len(emerald) {
			got = emerald[i]
		}
		if diffs := Compare(maskOptions(want),
			maskOptions(got)); len(diffs) > 0 {
			t.Errorf("%q\n%s", cmd, Render(diffs))
		}
	}

	// A passing golden case is not evidence the code ran: the
	// three probes would agree on an abort too. So check the
	// transcript actually carries all three answers.
	if len(emerald) < 2 {
		t.Fatalf("no transcript for %q", script[1])
	}
	for _, want := range []string{
		"MUF ", "MPI ", "NEWER ", "OLDER ",
	} {
		if !strings.Contains(emerald[1], want) {
			t.Errorf("the probe never reached %q: %q",
				want, emerald[1])
		}
	}
}
