package main

type Camera struct {
	pan        Pos2
	zoom       float32
	zoomCenter Pos2
}

func (c Camera) Pan(delta Pos2) {
	c.pan.x += delta.x
	c.pan.y += delta.y
}

func (c Camera) Zoom(zoomFactor float32, mousePos Pos2) {
	beforeZoomWorld := c.ScreenToWorld(mousePos)
	c.zoom *= zoomFactor
	afterZoomWorld := c.ScreenToWorld(mousePos)

	c.pan.x += (afterZoomWorld.x - beforeZoomWorld.x) * c.zoom
	c.pan.y += (afterZoomWorld.y - beforeZoomWorld.y) * c.zoom
}

func (c Camera) WorldToScreen(world Pos2) Pos2 {
	screenX := c.zoomCenter.x + (world.x-c.zoomCenter.x)*c.zoom + c.pan.x
	screenY := c.zoomCenter.y + (world.y-c.zoomCenter.y)*c.zoom + c.pan.y
	return Pos2{x: screenX, y: screenY}
}

func (c Camera) ScreenToWorld(screen Pos2) Pos2 {
	worldX := (screen.x-c.zoomCenter.x-c.pan.x)/c.zoom + c.zoomCenter.x
	worldY := (screen.y-c.zoomCenter.y-c.pan.y)/c.zoom + c.zoomCenter.y
	return Pos2{x: worldX, y: worldY}
}
