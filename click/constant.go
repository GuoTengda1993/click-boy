package click

type ClickStatus int
type MsgType string

const (
	DeviceTypeAndroid = "android"
	DeviceTypeIOS     = "ios"
)

const (
	ClickStatusStop    ClickStatus = 0
	ClickStatusRunning ClickStatus = 1
	ClickStatusPause   ClickStatus = 2
)

const (
	EventDetectDevice = "detect-device"
	EventMessage      = "message"
	EventLoopNum      = "loop-num"
)

const (
	MsgSuccess  MsgType = "success"
	MsgInfo     MsgType = "info"
	MsgWarnning MsgType = "warming"
	MsgError    MsgType = "error"
)
