package view

import (
	"context"
	"time"

	"github.com/derailed/k9s/internal"
	"github.com/derailed/k9s/internal/ui"
	"github.com/gdamore/tcell"
	"github.com/rs/zerolog/log"
)

type Table struct {
	*ui.Table

	app     *App
	enterFn EnterFunc
}

func NewTable(gvr string) *Table {
	return &Table{
		Table: ui.NewTable(gvr),
	}
}

// Init initializes the component
func (t *Table) Init(ctx context.Context) (err error) {
	if t.app, err = extractApp(ctx); err != nil {
		return err
	}
	ctx = context.WithValue(ctx, internal.KeyStyles, t.app.Styles)
	t.Table.Init(ctx)
	t.bindKeys()
	t.GetModel().SetRefreshRate(time.Duration(t.app.Config.K9s.GetRefreshRate()) * time.Second)

	return nil
}

// Name returns the table name.
func (t *Table) Name() string { return t.BaseTitle }

// App returns the current app handle.
func (t *Table) App() *App {
	return t.app
}

// Start runs the component.
func (t *Table) Start() {
	t.Stop()
	t.SearchBuff().AddListener(t.app.Cmd())
	t.SearchBuff().AddListener(t)
}

// Stop terminates the component.
func (t *Table) Stop() {
	t.SearchBuff().RemoveListener(t.app.Cmd())
	t.SearchBuff().RemoveListener(t)
}

// SetEnterFn specifies the default enter behavior.
func (t *Table) SetEnterFn(f EnterFunc) {
	t.enterFn = f
}

// SetExtraActionsFn specifies custom keyboard behavior.
func (t *Table) SetExtraActionsFn(BoostActionsFunc) {}

// BufferChanged indicates the buffer was changed.
func (t *Table) BufferChanged(s string) {}

// BufferActive indicates the buff activity changed.
func (t *Table) BufferActive(state bool, k ui.BufferKind) {
	t.app.BufferActive(state, k)
}

func (t *Table) saveCmd(evt *tcell.EventKey) *tcell.EventKey {
	if path, err := saveTable(t.app.Config.K9s.CurrentCluster, t.BaseTitle, t.Path, t.GetFilteredData()); err != nil {
		t.app.Flash().Err(err)
	} else {
		t.app.Flash().Infof("File %s saved successfully!", path)
	}

	return nil
}

func (t *Table) bindKeys() {
	t.Actions().Add(ui.KeyActions{
		ui.KeySpace:         ui.NewSharedKeyAction("Mark", t.markCmd, false),
		tcell.KeyCtrlSpace:  ui.NewSharedKeyAction("Marks Clear", t.clearMarksCmd, false),
		tcell.KeyCtrlS:      ui.NewSharedKeyAction("Save", t.saveCmd, false),
		ui.KeySlash:         ui.NewSharedKeyAction("Filter Mode", t.activateCmd, false),
		tcell.KeyCtrlU:      ui.NewSharedKeyAction("Clear Filter", t.clearCmd, false),
		tcell.KeyBackspace2: ui.NewSharedKeyAction("Erase", t.eraseCmd, false),
		tcell.KeyBackspace:  ui.NewSharedKeyAction("Erase", t.eraseCmd, false),
		tcell.KeyDelete:     ui.NewSharedKeyAction("Erase", t.eraseCmd, false),
		ui.KeyShiftN:        ui.NewKeyAction("Sort Name", t.SortColCmd(0, true), false),
		ui.KeyShiftA:        ui.NewKeyAction("Sort Age", t.SortColCmd(-1, true), false),
	})
}

func (t *Table) markCmd(evt *tcell.EventKey) *tcell.EventKey {
	path := t.GetSelectedItem()
	if path == "" {
		return evt
	}
	t.ToggleMark()
	t.Refresh()

	return nil
}

func (t *Table) clearMarksCmd(evt *tcell.EventKey) *tcell.EventKey {
	path := t.GetSelectedItem()
	if path == "" {
		return evt
	}
	t.ClearMarks()

	return nil
}

func (t *Table) clearCmd(evt *tcell.EventKey) *tcell.EventKey {
	if !t.SearchBuff().IsActive() {
		return evt
	}
	t.SearchBuff().Clear()

	return nil
}

func (t *Table) eraseCmd(evt *tcell.EventKey) *tcell.EventKey {
	if t.SearchBuff().IsActive() {
		t.SearchBuff().Delete()
	}

	return nil
}

func (t *Table) activateCmd(evt *tcell.EventKey) *tcell.EventKey {
	log.Debug().Msgf("Table filter activated!")
	if t.app.InCmdMode() {
		log.Debug().Msgf("App Is in Command mode!")
		return evt
	}
	t.app.Flash().Info("Filter mode activated.")
	t.SearchBuff().SetActive(true)

	return nil
}
