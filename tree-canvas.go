package main

import (
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
	defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
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
	switch 5 {
	case 0:
		op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(c.scale, c.scale))).Add(gtx.Ops)
		for _, node := range tree.Nodes {
			if node.IsProxy {
				continue
			}
			drawCircle(gtx, node.Position.Mul(c.scale), 20, color.NRGBA{R: 0xff, A: 0xff})
		}
	case 1:
		op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(c.scale, c.scale))).Add(gtx.Ops)
		for _, node := range tree.Nodes {
			if node.IsProxy {
				continue
			}
			drawEllipse(gtx, node.Position.Mul(c.scale), 20, color.NRGBA{R: 0xff, A: 0xff})
		}
	case 2:
		op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(c.scale, c.scale))).Add(gtx.Ops)
		for _, node := range tree.Nodes {
			if node.IsProxy {
				continue
			}
			squashCircle(gtx, node.Position.Mul(c.scale), 20, color.NRGBA{R: 0xff, A: 0xff})
		}
	case 3:
		for _, node := range tree.Nodes {
			if node.IsProxy {
				continue
			}
			squashCircle(gtx, node.Position.Mul(c.scale), 20*c.scale, color.NRGBA{R: 0xff, A: 0xff})
		}
	case 4:
		for _, node := range tree.Nodes {
			if node.IsProxy {
				continue
			}
			squashCircle2(gtx, node.Position.Mul(c.scale), 20*c.scale, color.NRGBA{R: 0xff, A: 0xff})
		}
	case 5:
		var path clip.Path
		path.Begin(gtx.Ops)
		for _, node := range tree.Nodes {
			if node.IsProxy {
				continue
			}
			addSquashCircle(&path, node.Position.Mul(c.scale), 20*c.scale)
		}
		paint.FillShape(gtx.Ops, color.NRGBA{R: 0xff, A: 0xff}, clip.Outline{Path: path.End()}.Op())
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

func drawEllipse(gtx layout.Context, center f32.Point, radius int, clr color.NRGBA) {
	bounds := image.Rect(center.Round().X-radius, center.Round().Y-radius, center.Round().X+radius, center.Round().Y+radius)
	paint.FillShape(gtx.Ops, clr, clip.Ellipse(bounds).Op(gtx.Ops))
}

func squashCircle(gtx layout.Context, p f32.Point, r float32, color color.NRGBA) {
	defer op.Affine(f32.Affine2D{}.Offset(p)).Push(gtx.Ops).Pop()

	var path clip.Path
	path.Begin(gtx.Ops)
	path.Move(f32.Pt(0, -r))
	path.Cube(f32.Pt(r, 0), f32.Pt(r, 2*r*0.75), f32.Pt(0, 2*r*0.75))
	path.Cube(f32.Pt(-r, 0), f32.Pt(-r, -2*r*0.75), f32.Pt(0, -2*r*0.75))
	defer clip.Outline{Path: path.End()}.Op().Push(gtx.Ops).Pop()

	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func squashCircle2(gtx layout.Context, p f32.Point, r float32, color color.NRGBA) {
	var path clip.Path
	path.Begin(gtx.Ops)
	path.MoveTo(f32.Pt(p.X, p.Y-r))
	path.Cube(f32.Pt(r, 0), f32.Pt(r, 2*r*0.75), f32.Pt(0, 2*r*0.75))
	path.Cube(f32.Pt(-r, 0), f32.Pt(-r, -2*r*0.75), f32.Pt(0, -2*r*0.75))
	defer clip.Outline{Path: path.End()}.Op().Push(gtx.Ops).Pop()

	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func addSquashCircle(path *clip.Path, p f32.Point, r float32) {
	path.MoveTo(f32.Pt(p.X+0, p.Y-r))
	path.Cube(f32.Pt(r, 0), f32.Pt(r, 2*r*0.75), f32.Pt(0, 2*r*0.75))
	path.Cube(f32.Pt(-r, 0), f32.Pt(-r, -2*r*0.75), f32.Pt(0, -2*r*0.75))
}
