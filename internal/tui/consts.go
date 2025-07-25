package tui

import "time"

const (
	MinTerminalWidth       = 80
	MinTerminalHeight      = 20
	DefaultWidth           = 80
	DefaultHeight          = 24
	minLeftWidth           = 40
	maxLeftWidthMultiplier = 0.5
)

const (
	TickInterval          = 100 * time.Millisecond
	StreamTickInterval    = 50 * time.Millisecond
	ClipboardFeedbackTime = 3 * time.Second
)

const (
	OutputChannelBuffer = 100
)

const (
	ContentHeightOffset = 6
	ContentPadding      = 2
	BorderOffset        = 2
	PanelPadding        = 4
	ScrollSteps         = 3
)
