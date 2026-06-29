package main

import (
	"fmt"
	"os"

	"github.com/xmasengine/xmas/xgal"
	"github.com/xmasengine/xmas/xui"
)

var (
	iconSword  = newIcon(16, 16)
	iconStar   = newIcon(16, 16)
	iconBag    = newIcon(16, 16)
	iconBoot   = newIcon(16, 16)
	iconFire   = newIcon(16, 16)
	iconIce    = newIcon(16, 16)
	iconBolt   = newIcon(16, 16)
	iconPotion = newIcon(16, 16)
	iconEther  = newIcon(16, 16)
	iconHeart  = newIcon(16, 16)
	iconDrop   = newIcon(16, 16)
	iconCoin   = newIcon(16, 16)
	iconLvl    = newIcon(16, 16)
)

func newIcon(w, h int) *xgal.Surface {
	return xgal.Prepare(w, h)
}

func init() {
	// sword
	xgal.Line(iconSword, 8, 2, 8, 12, 2, xgal.Wash(200, 200, 200, 255))
	xgal.Line(iconSword, 5, 10, 11, 10, 2, xgal.Wash(160, 120, 60, 255))
	xgal.Line(iconSword, 4, 12, 12, 12, 1, xgal.Wash(160, 120, 60, 255))
	xgal.Box(iconSword, xgal.Rect(7, 1, 10, 3), xgal.Wash(200, 200, 200, 255))

	// magic star
	xgal.Ink(iconStar, xgal.BuiltinFace, xgal.Wash(200, 100, 255, 255), 3, 1, "+")

	// bag
	xgal.Box(iconBag, xgal.Rect(3, 4, 13, 14), xgal.Wash(160, 100, 50, 255))
	xgal.Outline(iconBag, xgal.Rect(3, 4, 13, 14), 1, xgal.Wash(100, 60, 20, 255))
	xgal.Box(iconBag, xgal.Rect(5, 2, 11, 5), xgal.Wash(160, 100, 50, 255))
	xgal.Outline(iconBag, xgal.Rect(5, 2, 11, 5), 1, xgal.Wash(100, 60, 20, 255))

	// boot
	xgal.Box(iconBoot, xgal.Rect(4, 3, 10, 11), xgal.Wash(140, 80, 40, 255))
	xgal.Box(iconBoot, xgal.Rect(3, 9, 14, 13), xgal.Wash(140, 80, 40, 255))
	xgal.Outline(iconBoot, xgal.Rect(3, 9, 14, 13), 1, xgal.Wash(80, 50, 20, 255))

	// flame
	xgal.Disk(iconFire, xgal.Pt(8, 9), 5, xgal.Wash(255, 100, 0, 255))
	xgal.Disk(iconFire, xgal.Pt(7, 7), 3, xgal.Wash(255, 200, 0, 255))
	xgal.Disk(iconFire, xgal.Pt(8, 5), 2, xgal.Wash(255, 255, 100, 255))

	// ice crystal
	xgal.Line(iconIce, 8, 2, 8, 14, 2, xgal.Wash(150, 200, 255, 255))
	xgal.Line(iconIce, 4, 6, 12, 6, 1, xgal.Wash(150, 200, 255, 255))
	xgal.Line(iconIce, 5, 10, 11, 10, 1, xgal.Wash(150, 200, 255, 255))
	xgal.Disk(iconIce, xgal.Pt(8, 8), 4, xgal.Wash(180, 220, 255, 80))

	// bolt
	xgal.Box(iconBolt, xgal.Rect(9, 1, 13, 7), xgal.Wash(255, 220, 50, 255))
	xgal.Box(iconBolt, xgal.Rect(3, 6, 12, 10), xgal.Wash(255, 220, 50, 255))
	xgal.Box(iconBolt, xgal.Rect(4, 9, 8, 15), xgal.Wash(255, 220, 50, 255))

	// potion
	xgal.Box(iconPotion, xgal.Rect(6, 2, 10, 4), xgal.Wash(200, 50, 50, 255))
	xgal.Box(iconPotion, xgal.Rect(5, 4, 11, 13), xgal.Wash(200, 50, 50, 255))
	xgal.Outline(iconPotion, xgal.Rect(5, 4, 11, 13), 1, xgal.Wash(150, 20, 20, 255))
	xgal.Box(iconPotion, xgal.Rect(7, 6, 9, 10), xgal.Wash(100, 200, 100, 255))

	// ether
	xgal.Box(iconEther, xgal.Rect(6, 2, 10, 4), xgal.Wash(50, 100, 220, 255))
	xgal.Box(iconEther, xgal.Rect(5, 4, 11, 13), xgal.Wash(50, 100, 220, 255))
	xgal.Outline(iconEther, xgal.Rect(5, 4, 11, 13), 1, xgal.Wash(20, 50, 160, 255))
	xgal.Box(iconEther, xgal.Rect(7, 6, 9, 10), xgal.Wash(150, 200, 255, 255))

	// heart (HP)
	xgal.Disk(iconHeart, xgal.Pt(5, 6), 3, xgal.Wash(255, 50, 50, 255))
	xgal.Disk(iconHeart, xgal.Pt(11, 6), 3, xgal.Wash(255, 50, 50, 255))
	xgal.Box(iconHeart, xgal.Rect(5, 5, 11, 9), xgal.Wash(255, 50, 50, 255))
	xgal.Box(iconHeart, xgal.Rect(6, 9, 10, 12), xgal.Wash(255, 50, 50, 255))

	// drop (MP)
	xgal.Disk(iconDrop, xgal.Pt(8, 9), 5, xgal.Wash(80, 120, 255, 255))
	xgal.Box(iconDrop, xgal.Rect(6, 3, 10, 8), xgal.Wash(80, 120, 255, 255))

	// coin (GP)
	xgal.Disk(iconCoin, xgal.Pt(8, 8), 6, xgal.Wash(255, 200, 50, 255))
	xgal.Circle(iconCoin, xgal.Pt(8, 8), 6, 1, xgal.Wash(180, 130, 20, 255))
	xgal.Ink(iconCoin, xgal.BuiltinFace, xgal.Wash(180, 130, 20, 255), 6, 5, "G")

	// arrow up (Level)
	xgal.Box(iconLvl, xgal.Rect(7, 3, 9, 12), xgal.Wash(100, 220, 100, 255))
	xgal.Box(iconLvl, xgal.Rect(5, 5, 11, 7), xgal.Wash(100, 220, 100, 255))
}

type guiDemo struct {
	root   *xui.Layer
	popup  *xui.PopupLayer
	talk   *xui.TalkLayer
	ring   *xui.RingLayer
	screen *xui.ScreenLayer

	hud       *xui.HUDLayer
	hpBar     *xui.HUDBar
	mpBar     *xui.HUDBar
	levelStat *xui.HUDStat
	gpStat    *xui.HUDStat

	// synthetic game state
	hp, maxHP int
	mp, maxMP int
	level     int
	gp        int

	confirmPressed bool
	cancelPressed  bool
	dirDx, dirDy   int
}

const screenW = 320
const screenH = 240

func (g *guiDemo) Update() error {
	g.root.Poll()

	// track one-shot key presses for the frame
	g.confirmPressed = xgal.Tap(xgal.KeyEnter) || xgal.Tap(xgal.KeySpace)
	g.cancelPressed = xgal.Tap(xgal.KeyEscape)
	g.dirDx, g.dirDy = 0, 0
	if xgal.Tap(xgal.KeyArrowLeft) || xgal.Key(xgal.KeyA) {
		g.dirDx = -1
	}
	if xgal.Tap(xgal.KeyArrowRight) || xgal.Key(xgal.KeyD) {
		g.dirDx = 1
	}
	if xgal.Tap(xgal.KeyArrowUp) || xgal.Key(xgal.KeyW) {
		g.dirDy = -1
	}
	if xgal.Tap(xgal.KeyArrowDown) || xgal.Key(xgal.KeyS) {
		g.dirDy = 1
	}

	// toggle popup
	if xgal.Tap(xgal.KeyZ) && g.popup == nil {
		g.popup = xui.Popup(screenW, screenH, "The Whispering Woods")
		g.root.Add(g.popup)
	}

	// toggle talk window
	if xgal.Tap(xgal.KeyX) && g.talk == nil {
		// create a simple portrait (32×32 coloured square)
		portrait := xgal.Prepare(32, 32)
		xgal.Box(portrait, xgal.Rect(0, 0, 32, 32), xgal.Wash(120, 80, 40, 255))
		xgal.Disk(portrait, xgal.Pt(16, 8), 6, xgal.Wash(240, 200, 160, 255))
		xgal.Line(portrait, 8, 20, 24, 20, 2, xgal.Wash(0, 0, 0, 255))

		g.talk = xui.Talk(
			xgal.Rect(60, screenH*6/8, screenW-50, screenH*7/8),
			portrait,
			[]string{
				"Welcome, traveller!\nThe forest is dangerous ahead.",
				"Be sure to stock up on supplies before proceeding.",
				"May the good Lord be with you!",
			},
			func() bool { return g.confirmPressed },
		)
		g.root.Add(g.talk)
	}

	// toggle ring menu
	if xgal.Tap(xgal.KeyC) && g.ring == nil {
		subRing := xui.Ring(0, 0, 45, []xui.RingItem{
			{Label: "Fire", Icon: xui.Icon{Image: iconFire}, Action: func() { fmt.Println("ring: Fire") }},
			{Label: "Ice", Icon: xui.Icon{Image: iconIce}, Action: func() { fmt.Println("ring: Ice") }},
			{Label: "Thunder", Icon: xui.Icon{Image: iconBolt}, Action: func() { fmt.Println("ring: Thunder") }},
		})
		g.ring = xui.Ring(screenW/2, screenH/2, 80, []xui.RingItem{
			{Label: "Attack", Icon: xui.Icon{Image: iconSword}, Action: func() { fmt.Println("ring: Attack") }},
			{Label: "Magic", Icon: xui.Icon{Image: iconStar}, SubRing: subRing},
			{Label: "Items", Icon: xui.Icon{Image: iconBag}, Action: func() { fmt.Println("ring: Items") }},
			{Label: "Run", Icon: xui.Icon{Image: iconBoot}, Action: func() { fmt.Println("ring: Run") }},
		})
		g.root.Add(g.ring)
	}

	// toggle full-screen menu
	if xgal.Tap(xgal.KeyV) && g.screen == nil {
		itemsTab := xui.ScreenTab{
			Label: "Items",
			Items: []xui.TabItem{
				{Icon: xui.Icon{Image: iconPotion}, Label: "Potion", Value: "x3", Activate: func() { fmt.Println("used Potion") }},
				{Icon: xui.Icon{Image: iconEther}, Label: "Ether", Value: "x1", Activate: func() { fmt.Println("used Ether") }},
				{Icon: xui.Icon{Image: iconPotion}, Label: "Elixir", Value: "x0", Activate: func() { fmt.Println("no Elixir left!") }},
			},
		}
		settingsTab := xui.ScreenTab{
			Label: "Settings",
			Items: []xui.TabItem{
				{Label: "BGM Volume", Options: []string{"Off", "Low", "Mid", "High"}, OptIdx: 2},
				{Label: "SFX Volume", Options: []string{"Off", "Low", "Mid", "High"}, OptIdx: 2},
				{Label: "Fullscreen", Bool: new(bool)},
				{Label: "Hard Mode", Bool: new(bool)},
			},
		}

		g.screen = xui.Screen(screenW, screenH,
			[]xui.ScreenTab{
				{Label: "Status", Items: []xui.TabItem{
					{Icon: xui.Icon{Image: iconLvl}, Label: "Level", Value: fmt.Sprintf("%d", g.level)},
					{Icon: xui.Icon{Image: iconHeart}, Label: "HP", Value: fmt.Sprintf("%d/%d", g.hp, g.maxHP)},
					{Icon: xui.Icon{Image: iconDrop}, Label: "MP", Value: fmt.Sprintf("%d/%d", g.mp, g.maxMP)},
					{Icon: xui.Icon{Image: iconCoin}, Label: "GP", Value: fmt.Sprintf("%d", g.gp)},
				}},
				itemsTab,
				settingsTab,
			},
		)
		g.screen.Close = func() bool { return g.cancelPressed }
		g.root.Add(g.screen)
	}

	// remove finished widgets
	var kept []xui.Widget
	for _, kid := range g.root.Kids {
		switch kid.(type) {
		case *xui.PopupLayer:
			g.popup = nil
		case *xui.TalkLayer:
			g.talk = nil
		case *xui.RingLayer:
			g.ring = nil
		case *xui.ScreenLayer:
			g.screen = nil
		}
		kept = append(kept, kid)
	}
	g.root.Kids = kept

	return nil
}

func (g *guiDemo) Draw(screen *xgal.Surface) {
	xgal.Clear(screen, xgal.Wash(20, 20, 40, 255))
	// grid background
	for x := 0; x < screenW; x += 32 {
		xgal.Line(screen, x, 0, x, screenH, 1, xgal.Wash(30, 30, 50, 255))
	}
	for y := 0; y < screenH; y += 32 {
		xgal.Line(screen, 0, y, screenW, y, 1, xgal.Wash(30, 30, 50, 255))
	}

	// help text
	xgal.Ink(screen, xgal.BuiltinFace, xgal.Wash(200, 200, 200, 255), 10, screenH-40,
		"Z:Popup  X:Talk  C:Ring  V:Screen  Arrows/ WASD: Navigate  Enter: Confirm  Esc: Cancel")

	g.root.Render(screen)
}

func (g *guiDemo) Layout(w, h int) (int, int) {
	g.root.Place(xgal.Rect(0, 0, screenW, screenH))
	return screenW, screenH
}

func main() {
	g := &guiDemo{
		hp:    73,
		maxHP: 100,
		mp:    42,
		maxMP: 50,
		level: 5,
		gp:    340,
	}

	g.root = &xui.Layer{}
	g.root.Style = xui.DefaultStyle()

	g.hud = xui.HUD(screenW)
	g.hpBar = g.hud.AddBar("HP", g.hp, g.maxHP, xgal.Wash(255, 0, 0, 255))
	g.mpBar = g.hud.AddBar("MP", g.mp, g.maxMP, xgal.Wash(0, 0, 255, 255))

	g.levelStat = g.hud.AddStat("Mia Level: %d", g.level)
	g.gpStat = g.hud.AddStat("GP: %d", g.gp)
	g.root.Add(g.hud)

	xgal.Screen(-1, -1, "xui Game Widgets")
	xgal.Decorate(true)
	if err := xgal.Play(g); err != nil {
		fmt.Fprintf(os.Stderr, "RUNTIME: %v\n", err)
		os.Exit(2)
	}
}
