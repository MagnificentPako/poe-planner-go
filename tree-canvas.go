package main

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/input"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type TreeCanvas struct {
	offsetX, offsetY         float32
	scale                    float32
	dragging                 bool
	dragStartX, dragStartY   float32
	zoomCenterX, zoomCenterY float32
}

var canvasTag = new(bool)

func (c *TreeCanvas) Layout(gtx layout.Context, i input.Source, tree *ProcessedTree) layout.Dimensions {
	defer clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()
	event.Op(gtx.Ops, canvasTag)

	for {
		ev, ok := i.Event(pointer.Filter{Target: canvasTag, Kinds: pointer.Press | pointer.Release | pointer.Drag | pointer.Scroll, ScrollY: pointer.ScrollRange{Min: -1, Max: 1}})
		if !ok {
			break
		}

		if x, ok := ev.(pointer.Event); ok {
			switch x.Kind {
			case pointer.Press:
				c.dragging = true
				c.dragStartX = x.Position.X
				c.dragStartY = x.Position.Y
			case pointer.Drag:
				if c.dragging {
					c.offsetX += x.Position.X - c.dragStartX
					c.offsetY += x.Position.Y - c.dragStartY
					c.dragStartX = x.Position.X
					c.dragStartY = x.Position.Y
				}
			case pointer.Release:
				c.dragging = false
			case pointer.Scroll:
				if x.Scroll.Y < 0 {
					c.scale *= 1.1
				} else if x.Scroll.Y > 0 {
					c.scale /= 1.1
				}
			}
		}
	}

	op.Offset(image.Pt(int(c.offsetX), int(c.offsetY))).Add(gtx.Ops)
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(c.scale, c.scale))).Add(gtx.Ops)

	fmt.Println(len(tree.Nodes))

	for _, node := range tree.Nodes {
		if node.IsProxy {
			continue
		}
		drawCircle(gtx, node.Position, 20, color.NRGBA{R: 0xff, A: 0xff})
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func drawCircle(gtx layout.Context, center f32.Point, radius int, clr color.NRGBA) {
	bounds := image.Rect(center.Round().X-radius, center.Round().Y-radius, center.Round().X+radius, center.Round().Y+radius)
	defer clip.UniformRRect(bounds, 20).Op(gtx.Ops).Push(gtx.Ops).Pop()
	defer clip.Rect{Min: image.Pt(center.Round().X-radius, center.Round().Y-radius), Max: image.Pt(center.Round().X+radius, center.Round().Y+radius)}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: clr}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}
