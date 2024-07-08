package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"

	"gioui.org/f32"
)

type Node struct {
	Group      int  `json:"Group"`
	Orbit      int  `json:"Orbit"`
	OrbitIndex int  `json:"OrbitIndex"`
	IsProxy    bool `json:"IsProxy"`
}

type Group struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type TreeExport struct {
	Groups map[string]Group `json:"groups"`
	Nodes  map[string]Node  `json:"nodes"`
}

type ProcessedNode struct {
	Position f32.Point
	IsProxy  bool
}

type ProcessedTree struct {
	Nodes map[int]ProcessedNode
}

func LoadTreeExport() (*TreeExport, error) {
	var data *TreeExport
	content, err := os.ReadFile("data.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (t *TreeExport) NodePosition(node *Node) f32.Point {
	var ORBIT_RADII = [...]int{0, 82, 162, 335, 493, 662, 846}
	var ORBIT_ANGLES_16 = [...]int{0, 30, 45, 60, 90, 120, 135, 150, 180, 210, 225, 240, 270, 300, 315, 330}
	var ORBIT_ANGLES_40 = [...]int{0, 10, 20, 30, 40, 45, 50, 60, 70, 80, 90, 100, 110, 120, 130, 135, 140, 150, 160, 170, 180, 190, 200, 210, 220, 225, 230, 240, 250, 260, 270, 280, 290, 300, 310, 315, 320, 330, 340, 350}
	var ORBIT_NODES = [...]int{1, 6, 16, 16, 40, 72, 72}
	group := t.Groups[strconv.Itoa(node.Group)]
	radius := float64(ORBIT_RADII[node.Orbit])
	skillsOnOrbit := ORBIT_NODES[node.Orbit]
	orbitIndex := node.OrbitIndex
	twoPi := math.Pi

	var angle float64
	if skillsOnOrbit == 16 {
		angle = float64(ORBIT_ANGLES_16[orbitIndex])
	} else if skillsOnOrbit == 40 {
		angle = float64(ORBIT_ANGLES_40[orbitIndex])
	} else {
		angle = twoPi / float64(skillsOnOrbit) * float64(orbitIndex)
	}
	angle = ToRadians(angle)

	x := group.X + float32(radius*math.Sin(angle))
	y := group.Y - float32(radius*math.Cos(angle))

	return f32.Point{X: x, Y: y}
}

func ToRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func ProcessTree(t *TreeExport) ProcessedTree {
	var tree = ProcessedTree{Nodes: make(map[int]ProcessedNode)}
	fmt.Println(len(t.Nodes))
	for nodeId, node := range t.Nodes {
		var pNode = ProcessedNode{Position: t.NodePosition(&node), IsProxy: node.IsProxy}
		pId, _ := strconv.Atoi(nodeId)
		tree.Nodes[pId] = pNode
	}
	return tree
}
