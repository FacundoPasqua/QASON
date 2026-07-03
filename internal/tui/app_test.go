package tui

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// collectMsgs executes cmd and flattens the result. When cmd() yields a
// tea.BatchMsg (as returned by tea.Batch), it recurses into each of the
// batched commands so callers can find a specific message (e.g.
// actionDoneMsg) without caring whether it was wrapped in a batch.
func collectMsgs(cmd tea.Cmd) []tea.Msg {
	if cmd == nil {
		return nil
	}
	msg := cmd()
	if batch, ok := msg.(tea.BatchMsg); ok {
		var out []tea.Msg
		for _, c := range batch {
			out = append(out, collectMsgs(c)...)
		}
		return out
	}
	return []tea.Msg{msg}
}

// findActionDoneMsg searches msgs for an actionDoneMsg, returning it and
// whether one was found.
func findActionDoneMsg(msgs []tea.Msg) (actionDoneMsg, bool) {
	for _, m := range msgs {
		if done, ok := m.(actionDoneMsg); ok {
			return done, true
		}
	}
	return actionDoneMsg{}, false
}

// TestMenuNavigation verifies cursor movement on screenMenu clamps at the
// top (0) and bottom (2) of the 3-item menu, using both vim-style runes
// and arrow keys.
func TestMenuNavigation(t *testing.T) {
	tests := []struct {
		name    string
		start   int
		key     tea.KeyMsg
		wantEnd int
	}{
		{
			name:    "down from 0 moves to 1",
			start:   0,
			key:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")},
			wantEnd: 1,
		},
		{
			name:    "up from 0 stays at 0",
			start:   0,
			key:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")},
			wantEnd: 0,
		},
		{
			name:    "down at 2 stays at 2",
			start:   2,
			key:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")},
			wantEnd: 2,
		},
		{
			name:    "up from 2 moves to 1",
			start:   2,
			key:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")},
			wantEnd: 1,
		},
		{
			name:    "down arrow from 0 moves to 1",
			start:   0,
			key:     tea.KeyMsg{Type: tea.KeyDown},
			wantEnd: 1,
		},
		{
			name:    "up arrow from 2 moves to 1",
			start:   2,
			key:     tea.KeyMsg{Type: tea.KeyUp},
			wantEnd: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			m.cursor = tt.start

			updated, _ := m.Update(tt.key)
			m = updated.(Model)

			if m.cursor != tt.wantEnd {
				t.Errorf("cursor = %d, want %d", m.cursor, tt.wantEnd)
			}
		})
	}
}

// TestMenuSelectInstallRunsActionAndFinishes verifies that pressing enter
// on the Install menu item transitions to screenWorking, returns a
// non-nil command, and that draining the command yields an actionDoneMsg
// built from the injected install func. Feeding that message back in
// should land on screenDone with the skill count visible in the view.
func TestMenuSelectInstallRunsActionAndFinishes(t *testing.T) {
	m := NewWithActions(
		func() (int, int, error) { return 3, 35, nil },
		func() error { return nil },
	)
	m.cursor = 0

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)

	if m.screen != screenWorking {
		t.Fatalf("screen = %v, want screenWorking", m.screen)
	}
	if cmd == nil {
		t.Fatal("Update returned nil cmd, want non-nil cmd to run the install action")
	}

	msgs := collectMsgs(cmd)
	done, ok := findActionDoneMsg(msgs)
	if !ok {
		t.Fatalf("no actionDoneMsg found among drained messages: %#v", msgs)
	}
	if done.agents != 3 || done.skills != 35 {
		t.Errorf("actionDoneMsg = %+v, want agents=3 skills=35", done)
	}

	updated, _ = m.Update(done)
	m = updated.(Model)

	if m.screen != screenDone {
		t.Fatalf("screen = %v, want screenDone", m.screen)
	}
	if !strings.Contains(m.View(), "35") {
		t.Errorf("View() = %q, want it to contain %q", m.View(), "35")
	}
}

// TestMenuSelectUninstall verifies that pressing enter on the Uninstall
// menu item invokes the injected uninstall func and, once the resulting
// actionDoneMsg is fed back into Update, lands on screenDone with a view
// that mentions uninstall.
func TestMenuSelectUninstall(t *testing.T) {
	called := false
	m := NewWithActions(
		func() (int, int, error) { return 0, 0, nil },
		func() error {
			called = true
			return nil
		},
	)
	m.cursor = 1

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)

	if m.screen != screenWorking {
		t.Fatalf("screen = %v, want screenWorking", m.screen)
	}
	if cmd == nil {
		t.Fatal("Update returned nil cmd, want non-nil cmd to run the uninstall action")
	}

	msgs := collectMsgs(cmd)
	done, ok := findActionDoneMsg(msgs)
	if !ok {
		t.Fatalf("no actionDoneMsg found among drained messages: %#v", msgs)
	}
	if !done.uninstall {
		t.Errorf("actionDoneMsg.uninstall = false, want true")
	}
	if !called {
		t.Error("fake uninstall func was not called")
	}

	updated, _ = m.Update(done)
	m = updated.(Model)

	if m.screen != screenDone {
		t.Fatalf("screen = %v, want screenDone", m.screen)
	}
	if !strings.Contains(strings.ToLower(m.View()), "uninstall") {
		t.Errorf("View() = %q, want it to contain %q (case-insensitive)", m.View(), "uninstall")
	}
}

// TestQuitPaths verifies every path that should quit the program from
// screenMenu: selecting Quit with enter, "q", and ctrl+c.
func TestQuitPaths(t *testing.T) {
	tests := []struct {
		name       string
		setupModel func() Model
		key        tea.KeyMsg
	}{
		{
			name: "enter on Quit item",
			setupModel: func() Model {
				m := New()
				m.cursor = 2
				return m
			},
			key: tea.KeyMsg{Type: tea.KeyEnter},
		},
		{
			name: "q key on menu",
			setupModel: func() Model {
				return New()
			},
			key: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")},
		},
		{
			name: "ctrl+c on menu",
			setupModel: func() Model {
				return New()
			},
			key: tea.KeyMsg{Type: tea.KeyCtrlC},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setupModel()

			_, cmd := m.Update(tt.key)
			if cmd == nil {
				t.Fatal("Update returned nil cmd, want a cmd that yields tea.QuitMsg")
			}

			msg := cmd()
			if _, ok := msg.(tea.QuitMsg); !ok {
				t.Errorf("cmd() = %#v (%T), want tea.QuitMsg", msg, msg)
			}
		})
	}
}

// TestActionErrorShownOnDoneScreen verifies that feeding an actionDoneMsg
// carrying an error transitions to screenDone and surfaces the error text
// in the rendered view.
func TestActionErrorShownOnDoneScreen(t *testing.T) {
	m := New()

	updated, _ := m.Update(actionDoneMsg{err: errors.New("boom disk full")})
	m = updated.(Model)

	if m.screen != screenDone {
		t.Fatalf("screen = %v, want screenDone", m.screen)
	}
	if !strings.Contains(m.View(), "boom disk full") {
		t.Errorf("View() = %q, want it to contain %q", m.View(), "boom disk full")
	}
}

// TestDoneScreenEnterQuits verifies that once on screenDone, pressing
// enter returns a command that yields tea.QuitMsg.
func TestDoneScreenEnterQuits(t *testing.T) {
	m := New()
	updated, _ := m.Update(actionDoneMsg{agents: 3, skills: 35})
	m = updated.(Model)

	if m.screen != screenDone {
		t.Fatalf("precondition failed: screen = %v, want screenDone", m.screen)
	}

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Update returned nil cmd, want a cmd that yields tea.QuitMsg")
	}

	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Errorf("cmd() = %#v (%T), want tea.QuitMsg", msg, msg)
	}
}

// TestViewShowsBrandAndMenu verifies that the fresh model's menu view
// contains the QASON brand and all three menu items.
func TestViewShowsBrandAndMenu(t *testing.T) {
	m := New()
	view := m.View()

	for _, want := range []string{"QASON", "Install", "Uninstall", "Quit"} {
		if !strings.Contains(view, want) {
			t.Errorf("View() = %q, want it to contain %q", view, want)
		}
	}
}
