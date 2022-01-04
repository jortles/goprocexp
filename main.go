package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/process"
	"goprocexep/helpers"
	"sort"
	"strconv"
)

// Window handles windows
type Window func() (title string, content tview.Primitive)

type InfoUI struct {
	app   *tview.Application
	panel *tview.Flex
}

type ProcessView struct {
	Layout     *tview.Pages
	Table      *tview.Table
	Logger     *helpers.Logger
	detailsBox *tview.TextView
	details    *tview.Table
	procchan   chan helpers.Notification
}

type Process struct {
	Name          string
	Pid           int32
	CPUPercent    float64
	Username      string
	Executable    string
	CommandLine   string
	MemoryPercent float32
}

type ProcessDetails struct {
	CmdLine       string
	MemoryPercent float32
}

func (p *ProcessDetails) commandLine() string {
	procs, _ := process.Processes()
	for _, proc := range procs {
		p.CmdLine, _ = proc.Cmdline()
	}

	return p.CmdLine
}

func (p *ProcessDetails) memPercent() float32 {
	procs, _ := process.Processes()
	for _, proc := range procs {
		p.MemoryPercent, _ = proc.MemoryPercent()
	}

	return p.MemoryPercent
}

func (view *ProcessView) GetView() (title string, content tview.Primitive) {
	return "Processes", view.Layout
}

func (view *ProcessView) Init(app *tview.Application) {
	procs, _ := process.Processes()
	//var id string
	var s Process

	view.Layout = tview.NewPages()
	mainLayout := tview.NewFlex()
	mainLayout.SetDirection(tview.FlexRow).SetBorder(false)

	view.Table = tview.NewTable()
	view.Table.SetFixed(1, 1)
	view.Table.SetBorders(false).SetSeparator(tview.Borders.Vertical)
	view.Table.SetBorderPadding(0, 0, 0, 0)
	view.Table.SetSelectable(true, false)

	view.Table.SetCell(0, 1, tview.NewTableCell("PID").SetTextColor(tcell.ColorMediumPurple).SetSelectable(false).SetAlign(tview.AlignLeft))
	view.Table.SetCell(0, 2, tview.NewTableCell("Process Name").SetTextColor(tcell.ColorMediumPurple).SetSelectable(false).SetAlign(tview.AlignCenter))
	view.Table.SetCell(0, 3, tview.NewTableCell("CPU %").SetTextColor(tcell.ColorMediumPurple).SetSelectable(false).SetAlign(tview.AlignLeft))
	view.Table.SetCell(0, 4, tview.NewTableCell("User").SetTextColor(tcell.ColorMediumPurple).SetSelectable(false))
	view.Table.SetCell(0, 5, tview.NewTableCell("Path").SetTextColor(tcell.ColorMediumPurple).SetSelectable(false))

	detailsView := tview.NewFlex()
	detailsView.SetDirection(tview.FlexRow)
	detailsView.SetBorder(true).SetBorderColor(tcell.ColorIndigo).SetTitle(" PROCESS DETAILS ").SetTitleColor(tcell.ColorRed)

	view.details = tview.NewTable()
	view.details.SetBorders(false)
	view.details.SetCellSimple(1, 0, "Command:")
	view.details.GetCell(1, 0).SetAlign(tview.AlignRight)
	view.details.SetCellSimple(2, 0, "Memory %:")
	view.details.GetCell(2, 0).SetAlign(tview.AlignRight)

	view.detailsBox = tview.NewTextView().SetWrap(false).SetDynamicColors(true)

	func() {
		for _, proc := range procs {
			n := view.Table.GetRowCount()

			s.Name, _ = proc.Name()
			s.Pid = proc.Pid
			s.CPUPercent, _ = proc.CPUPercent()
			s.Username, _ = proc.Username()
			s.Executable, _ = proc.Exe()
			s.CommandLine, _ = proc.Cmdline()
			s.MemoryPercent, _ = proc.MemoryPercent()

			view.Table.SetCell(n, 1, tview.NewTableCell(strconv.Itoa(int(s.Pid))))
			view.Table.SetCell(n, 2, tview.NewTableCell(s.Name))
			view.Table.SetCell(n, 3, tview.NewTableCell(strconv.Itoa(int(s.CPUPercent))))
			view.Table.SetCell(n, 4, tview.NewTableCell(s.Username))
			view.Table.SetCell(n, 5, tview.NewTableCell(s.Executable))

		}
	}()

	detailsView.AddItem(view.details, 0, 1, false)

	mainLayout.AddItem(view.Table, 0, 2, false)
	mainLayout.AddItem(detailsView, 0, 2, false)

	// Deal with item focus
	items := []tview.Primitive{view.Table, detailsView, view.detailsBox}

	/*
		view.Table.SetSelectionChangedFunc(func(row int, column int) {
			if row > view.Table.GetRowCount() || row < 0 {
				return
			}

			view.detailsBox.Clear()

			id = view.Table.GetCell(row, 1).Text
			if entry := view.Logger.GetEntry(id); entry != nil {
				view.writeDetails()
				view.detailsBox.ScrollToBeginning()
			}
		})
	*/

	// Tab to next layout
	mainLayout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			for index, primitive := range items {
				if primitive == app.GetFocus() {
					app.SetFocus(items[(index-1+len(items))%len(items)])
					return event
				}
			}

			app.SetFocus(items[0])
		}

		return event
	})

	view.Layout.AddPage("mainLayout", mainLayout, true, true)

}

func (view *ProcessView) reloadtable() {
	view.Table.Clear()

	view.Table.SetCell(0, 1, tview.NewTableCell("Process Name").SetTextColor(tcell.ColorMediumPurple).SetSelectable(false).SetAlign(tview.AlignCenter))
	view.Table.SetCell(0, 2, tview.NewTableCell("PID").SetTextColor(tcell.ColorMediumPurple).SetSelectable(false).SetAlign(tview.AlignCenter))
	view.Table.SetCell(0, 3, tview.NewTableCell("CPU Percent").SetTextColor(tcell.ColorMediumPurple).SetSelectable(false))
	view.Table.SetCell(0, 4, tview.NewTableCell("User").SetTextColor(tcell.ColorMediumPurple).SetSelectable(false))
	view.Table.SetCell(0, 5, tview.NewTableCell("Path").SetTextColor(tcell.ColorMediumOrchid).SetSelectable(false))

	var processentries []*helpers.Entry
	for _, value := range view.Logger.GetEntries() {
		processentries = append(processentries, value)
	}

	sort.Slice(processentries, func(i, j int) bool {
		return processentries[i].StartedDateTime.Before(processentries[j].StartedDateTime)
	})

	for _, v := range processentries {
		view.AddEntry(v, 2)
	}
}

func (view *ProcessView) AddEntry(e *helpers.Entry, t int) {
	n := view.Table.GetRowCount()
	procs, _ := process.Processes()

	switch t {
	case 0:
		for _, proc := range procs {
			name, _ := proc.Name()
			pid := proc.Pid
			cpu, _ := proc.CPUPercent()
			usr, _ := proc.Username()
			exe, _ := proc.Exe()

			view.Table.SetCell(n, 1, tview.NewTableCell(name).SetSelectable(true))
			view.Table.SetCell(n, 2, tview.NewTableCell(strconv.Itoa(int(pid))).SetSelectable(false))
			view.Table.SetCell(n, 3, tview.NewTableCell(strconv.Itoa(int(cpu))).SetSelectable(false))
			view.Table.SetCell(n, 4, tview.NewTableCell(usr).SetSelectable(false))
			view.Table.SetCell(n, 5, tview.NewTableCell(exe).SetSelectable(false))

		}
	}
}

func main() {
	app := tview.NewApplication()

	processView := new(ProcessView)
	processView.Init(app)

	pages := []Window{
		processView.GetView,
	}

	// Main Layout
	mainWindow := tview.NewPages()
	footer := tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetWrap(false)
	footer.SetHighlightedFunc(func(added, removed, remaining []string) {
		mainWindow.SwitchToPage(added[0])

		_, p := mainWindow.GetFrontPage()
		app.SetFocus(p)
	})

	// Create the pages for all slides
	prevPage := func() {
		slide, _ := strconv.Atoi(footer.GetHighlights()[0])
		slide = (slide - 1 + len(pages)) % len(pages)
		footer.Highlight(strconv.Itoa(slide)).
			ScrollToHighlight()
	}
	nextPage := func() {
		slide, _ := strconv.Atoi(footer.GetHighlights()[0])
		slide = (slide + 1) % len(pages)
		footer.Highlight(strconv.Itoa(slide)).
			ScrollToHighlight()
	}

	for index, slide := range pages {
		title, primitive := slide()
		mainWindow.AddPage(strconv.Itoa(index), primitive, true, index == 0)
		fmt.Fprintf(footer, `%d ["%d"][mediumpurple]%s[white][""] `, index+1, index, title)
	}
	footer.Highlight("0")

	// Create the main layout
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(mainWindow, 0, 1, false)
	layout.AddItem(footer, 1, 1, false)

	// Add a quit modal
	quitModal := tview.NewModal().
		SetText("Quit?").
		AddButtons([]string{"Quit", "Cancel"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Quit" {
			app.Stop()
		} else {
			mainWindow.HidePage("quitModal")
		}
	})

	mainWindow.AddPage("quitModal", quitModal, false, false)

	// Shortcuts to navigate the slides
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			mainWindow.ShowPage("quitModal")
			app.SetFocus(quitModal)
			return nil
		case tcell.KeyRight:
			nextPage()
		case tcell.KeyLeft:
			prevPage()
		}

		return event
	})

	if err := app.SetRoot(layout, true).EnableMouse(true).SetFocus(processView.Table).Run(); err != nil {
		panic(err)
	}
}
