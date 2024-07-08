package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/input"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

type C = layout.Context
type D = layout.Dimensions

type Pos2 struct {
	x float32
	y float32
}

type AppState struct {
	tree   *ProcessedTree
	canvas *TreeCanvas
}

func main() {
	tree, _ := LoadTreeExport()
	pTree := ProcessTree(tree)
	canvas := TreeCanvas{scale: 1.0}
	state := AppState{
		tree:   &pTree,
		canvas: &canvas,
	}

	go func() {
		window := new(app.Window)
		err := run(window, state)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window, state AppState) error {
	theme := material.NewTheme()
	var ops op.Ops
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			RenderApp(gtx, theme, state, e.Source)

			e.Frame(gtx.Ops)
		}
	}
}

func RenderApp(gtx C, theme *material.Theme, state AppState, i input.Source) {
	layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return TopBar(gtx, theme)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEnd}.Layout(gtx,
				layout.Flexed(0.18, func(gtx C) D {
					return SideBar(gtx, theme)
				}),
				layout.Flexed(0.82, func(gtx C) D {
					return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd}.Layout(gtx,
						layout.Flexed(0.95, func(gtx C) D {
							return state.canvas.Layout(gtx, i, state.tree)
						}),
						layout.Rigid(func(gtx C) D {
							return BottomBar(gtx, theme)
						}),
					)
				}),
			)
		}),
	)
}

func TopBar(gtx C, theme *material.Theme) D {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEnd}.Layout(gtx, layout.Rigid(func(gtx C) D {
		text := material.Body1(theme, "Top Bar")
		return text.Layout(gtx)
	}))
}

func SideBar(gtx C, theme *material.Theme) D {
	return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd}.Layout(gtx, layout.Rigid(func(gtx C) D {
		text := material.Body1(theme, "Side Bar")
		return text.Layout(gtx)
	}))
}

func BottomBar(gtx C, theme *material.Theme) D {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceEnd}.Layout(gtx, layout.Rigid(func(gtx C) D {
		text := material.Body1(theme, "Bottom Bar")
		return text.Layout(gtx)
	}))
}
