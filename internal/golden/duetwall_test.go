package golden

import (
	"context"
	"testing"
)

// `@wall` goes to every **descriptor** through `notify_listeners`
// (`speech.c:125`), not straight to each connection — and this sent
// it straight.
//
// Two things follow, and neither could be observed when this case was
// written. The ignore filter applies — but `@wall` is wizard-only
// and `ignore_is_ignoring_sub` refuses outright when either owner is
// a wizard, so for any shouter who may use the command the filter is
// **dead**, and that one still cannot be seen. The other is the
// per-descriptor arithmetic: a player with two connections is
// notified once per descriptor, each notify reaching all of their
// descriptors, so they hear it twice over. That needed a third seat,
// which exists now — see `trio_test.go`. The change here was the
// faithful shape rather than a fix: it puts `@wall` through the one
// function that holds the rules.
//
// What two seats *do* show here is that a shout crosses rooms where
// `say` does not, and that a quelled wizard is refused the command
// entirely.
//
// The connect and disconnect announcements go with it, and what this
// comment said about their gating was **backwards** — the third
// seat corrected it. The "has connected." line is gated on DARK alone
// and fires on *every* connection, as do the propqueues; the
// `connect` **action** and the puppets waking are what wait for the
// first, which upstream's own comment calls odd. What two seats add
// here is that the *other* player sees the arrival at all.
const duetWallSource = `: main 1 pop ;`

var duetWallSetup = Script{
	"@tune penny_rate=0",
	"@pcreate Bob=secret",
}

var duetWall = Duet{
	Say("@teleport *Bob=#0"),

	// Heard by both, with the shouter's own name on it.
	Say("@wall everybody out"),

	// `say` for contrast: a room line, which both hear because
	// they are in the same room.
	Say("say in here"),

	// And from the other side, to show the wizard-only command
	// refusing a mortal.
	Hear("@wall let me out"),

	// Bob leaves the room: a wall still reaches him, which is the
	// whole difference between @wall and say.
	Say("@dig Far"),
	Say("@teleport *Bob=#5"),
	Say("@wall still here"),
	Say("say only in here"),

	// A quelled wizard may not @wall at all.
	Say("@set me=quell"),
	Say("@wall can you hear me"),
	Say("@set me=!quell"),
	Say("@wall and now"),
}

// TestDuetWallMatchesFuzzball compares both seats, step by step.
func TestDuetWallMatchesFuzzball(t *testing.T) {
	requireOracle(t)
	ctx := context.Background()

	fx, err := WriteFixture(t.TempDir(), duetWallSource)
	if err != nil {
		t.Fatal(err)
	}
	oracle, err := RunOracleDuet(ctx, fx, duetWallSetup,
		bobOnly, duetWall)
	if err != nil {
		t.Fatalf("driving the oracle: %v", err)
	}
	emerald, err := RunEmeraldDuet(ctx, fx, duetWallSetup,
		bobOnly, duetWall)
	if err != nil {
		t.Fatalf("driving this server: %v", err)
	}
	CompareDuets(t, duetWall, oracle, emerald)
}
