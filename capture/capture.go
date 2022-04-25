package capture

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"vtools/tools"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Sniper struct {
	widgets.QWidget
	_             func(string) `signal:"setOcrData"`
	desktopPixmap *gui.QPixmap
	capturePixmap *gui.QPixmap
	pixmapBegain  *core.QPoint
	pixmapEnd     *core.QPoint
	isMousePress  bool
}

func (si *Sniper) init() {
	si.SetWindowTitle("Sniper")
	si.SetWindowFlags(core.Qt__Dialog | core.Qt__FramelessWindowHint | core.Qt__WindowStaysOnTopHint)
	si.SetWindowState(core.Qt__WindowActive | core.Qt__WindowFullScreen)
	palette := gui.NewQPalette()
	si.SetPalette(palette)
}

func (si *Sniper) clear() {
	if !si.desktopPixmap.IsNull() {
		si.desktopPixmap.DestroyQPixmap()
	}
	if !si.capturePixmap.IsNull() {
		si.capturePixmap.DestroyQPixmap()
	}

	si.capturePixmap = nil
	si.pixmapBegain = nil
	si.pixmapEnd = nil
	si.capturePixmap = nil
}

func (si *Sniper) Recognize() {
	// si.clear()
	screen := gui.QGuiApplication_PrimaryScreen()
	//desktop window id is 0
	si.desktopPixmap = screen.GrabWindow(0, 0, 0, screen.Size().Width(), screen.Size().Height())

	si.Palette().SetBrush(si.BackgroundRole(), gui.NewQBrush7(si.desktopPixmap))
	si.pixmapBegain = core.NewQPoint()
	si.pixmapEnd = core.NewQPoint()
	si.SetMouseTracking(true)
	gui.QGuiApplication_SetOverrideCursor(gui.NewQCursor2(core.Qt__CrossCursor))
	si.eventSet()
	si.Show()
}

func (si *Sniper) eventSet() {
	si.mouseMoveEventSet()
	si.mousePressEventSet()
	si.mouseReleaseEventSet()
	si.paintEventSet()
	si.keyPressEventSet()
}

func (si *Sniper) mouseMoveEventSet() {
	si.ConnectMouseMoveEvent(func(event *gui.QMouseEvent) {
		// fmt.Println("mouseMoveEventSet")
		if si.isMousePress {
			si.pixmapEnd = event.Pos()
			si.Update()
		}
	})
}

func (si *Sniper) mousePressEventSet() {
	si.ConnectMousePressEvent(func(event *gui.QMouseEvent) {
		// fmt.Println("mousePressEventSet")
		if event.Button() == core.Qt__LeftButton {
			si.isMousePress = true
			si.pixmapBegain = event.Pos()
		}
	})
}

func (si *Sniper) mouseReleaseEventSet() {
	si.ConnectMouseReleaseEvent(func(event *gui.QMouseEvent) {
		si.isMousePress = false
		si.pixmapEnd = event.Pos()
		gui.QGuiApplication_SetOverrideCursor(gui.NewQCursor2(core.Qt__ArrowCursor))
	})
}

func (si *Sniper) paintEventSet() {
	si.ConnectPaintEvent(func(event *gui.QPaintEvent) {
		painter := gui.NewQPainter2(si)
		shadowColor := gui.NewQColor3(0, 0, 0, 100)
		painter.FillRect6(si.desktopPixmap.Rect(), shadowColor)
		if si.isMousePress {
			pen := gui.NewQPen()
			pen.SetColor(gui.NewQColor6("white"))
			// pen.SetStyle(core.Qt__SolidLine)
			pen.SetStyle(core.Qt__DashDotDotLine)
			pen.SetWidth(1)
			pen.SetCapStyle(core.Qt__FlatCap)
			painter.SetPen(pen)
			selectRect := si.getRect(si.pixmapBegain, si.pixmapEnd)
			si.capturePixmap = si.desktopPixmap.Copy(selectRect)
			// si.capturePixmap.Save("capture.png", "PNG", 100)
			painter.DrawPixmap8(selectRect.TopLeft(), si.capturePixmap)
			painter.DrawRect3(selectRect)
		}
		painter.End()
	})
}

func (si *Sniper) keyPressEventSet() {
	si.ConnectKeyPressEvent(func(event *gui.QKeyEvent) {
		if event.Key() == int(core.Qt__Key_Escape) {
			if si.capturePixmap.IsNull() ||
				si.capturePixmap.Size().Width() == 0 ||
				si.capturePixmap.Size().Height() == 0 {
				si.SetOcrData("")
				si.Close()
			} else {
				save_path := path.Join(os.TempDir(), "vtext_temp_capture.png")
				si.capturePixmap.Save(save_path, "PNG", 100)
				fmt.Println(save_path)
				defer os.Remove(save_path)
				image_bytes, _ := ioutil.ReadFile(save_path)
				go func() {
					text, err := tools.Ocr("123", image_bytes)
					if err != nil {
						fmt.Println(err)
					}
					si.SetOcrData(text)
				}()
				si.Hide()
			}
		}
	})
}

func (si *Sniper) getRect(begin *core.QPoint, end *core.QPoint) *core.QRect {
	x := begin.X()
	y := begin.Y()
	w := end.X() - begin.X()
	h := end.Y() - begin.Y()
	if w < 0 {
		x = end.X()
		w = -w
	}
	if h < 0 {
		y = end.Y()
		h = -h
	}
	rect := core.NewQRect4(x, y, w, h)
	if rect.Width() == 0 {
		rect.SetWidth(1)
	}
	if rect.Height() == 0 {
		rect.SetHeight(1)
	}
	return rect
}

func NNewSniper(parent widgets.QWidget_ITF, ff core.Qt__WindowType, ocrDataFunc func(text string)) *Sniper {
	si := NewSniper(parent, ff)
	si.init()
	fmt.Println("NewSniper")
	si.ConnectSetOcrData(func(t string) {
		ocrDataFunc(t)
		si.DeleteLater()
	})
	return si
}
