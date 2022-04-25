package capture

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type Toolbar struct {
	*widgets.QWidget
}

func NewToolBar() *Toolbar {
	tl := &Toolbar{widgets.NewQWidget(nil, 0)}
	tl.SetWindowFlags(core.Qt__FramelessWindowHint | core.Qt__WindowStaysOnTopHint)
	tl.SetAttribute(core.Qt__WA_TranslucentBackground, true)
	return tl
}
