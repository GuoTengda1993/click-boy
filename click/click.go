package click

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/danielpaulus/go-ios/ios"
	"github.com/electricbubble/gadb"
	"github.com/sunls24/gwda"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Click struct {
	ctx    context.Context
	device *Device
	status ClickStatus
	lock   *sync.RWMutex
	resume chan bool
}

type Device struct {
	adb       gadb.Client
	adbDev    gadb.Device
	iosDriver gwda.WebDriver
	Status    bool   `json:"status"`
	Type      string `json:"type"`
	Serial    string `json:"serial"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type Point struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Event string `json:"event"`
	Edit  bool   `json:"edit"`
}

type Params struct {
	Times    int    `json:"times"`
	Interval int    `json:"interval"`
	Duration int    `json:"duration"`
	Package  string `json:"package"`  // android for appPackage, ios for bundleId
	Activity string `json:"activity"` // just for android
}

type Message struct {
	Msg  string  `json:"msg"`
	Type MsgType `json:"type"`
}

var iosKeyMap map[int]gwda.DeviceButton = map[int]gwda.DeviceButton{
	3:  gwda.DeviceButtonHome,
	24: gwda.DeviceButtonVolumeUp,
	25: gwda.DeviceButtonVolumeDown,
}

func NewClick(ctx context.Context) *Click {
	c := &Click{ctx: ctx, lock: &sync.RWMutex{}, resume: make(chan bool, 1), device: &Device{}}
	return c
}

func (c *Click) detectDevice() {
	for {
		if adb, err := gadb.NewClient(); err == nil {
			if devList, err := adb.DeviceList(); err == nil && len(devList) > 0 {
				c.device.adb = adb
				c.device.adbDev = devList[0]
				c.device.Serial = devList[0].Serial()
				c.device.Type = DeviceTypeAndroid
				c.device.Status = true
				runtime.EventsEmit(c.ctx, EventDetectDevice, c.device)
				return
			}
		}
		if iosList, err := ios.ListDevices(); err == nil && len(iosList.DeviceList) > 0 {
			c.device.Type = DeviceTypeIOS
			c.device.Serial = iosList.DeviceList[0].Properties.SerialNumber
			c.device.Status = false
			if driver, err := gwda.NewUSBDriver(nil); err == nil {
				c.device.iosDriver = driver
				if devSize, err := driver.WindowSize(); err == nil && devSize.Width > 0 {
					c.device.Status = true
					c.device.Width = devSize.Width
					c.device.Height = devSize.Height
					runtime.EventsEmit(c.ctx, EventDetectDevice, c.device)
					return
				}
			}
			runtime.EventsEmit(c.ctx, EventDetectDevice, c.device)
		}

		time.Sleep(1 * time.Second)
	}
}

func (c *Click) ConnectDevice() {
	go c.detectDevice()
}

func (c *Click) ScreenShot() string {
	if !c.device.Status {
		c.sendMsg("device not connected", MsgError)
		return ""
	}
	switch c.device.Type {
	case DeviceTypeAndroid:
		return c.screenShotAndroid()
	case DeviceTypeIOS:
		return c.screenShotIOS()
	default:
		c.sendMsg("device not support", MsgError)
		return ""
	}
}

func (c *Click) screenShotAndroid() string {
	picName := fmt.Sprintf("screeshot-%d.png", time.Now().Unix())
	picDevPath := fmt.Sprintf("/storage/emulated/0/DCIM/%s", picName)
	cmd := fmt.Sprintf("screencap -p %s", picDevPath)
	if _, err := c.device.adbDev.RunShellCommand(cmd); err != nil {
		c.sendMsg(fmt.Sprintf("android screen shot error: %s", err.Error()), MsgError)
		return ""
	}
	var buf bytes.Buffer
	if err := c.device.adbDev.Pull(picDevPath, &buf); err != nil {
		c.sendMsg(fmt.Sprintf("android screen shot error: %s", err.Error()), MsgError)
		return ""
	}
	_, _ = c.device.adbDev.RunShellCommand(fmt.Sprintf("rm -f %s", picDevPath))
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

func (c *Click) screenShotIOS() string {
	screenshot, err := c.device.iosDriver.Screenshot()
	if err != nil {
		c.sendMsg(fmt.Sprintf("iPhone screen shot error: %s", err.Error()), MsgError)
		return ""
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(screenshot.Bytes())
}

func (c *Click) GetStatus() ClickStatus {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.status
}

func (c *Click) setStatus(s ClickStatus) ClickStatus {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.status = s
	return s
}

func (c *Click) Pause() ClickStatus {
	return c.setStatus(ClickStatusPause)
}

func (c *Click) Resume() ClickStatus {
	c.resume <- true
	return c.setStatus(ClickStatusRunning)
}

func (c *Click) Stop() ClickStatus {
	if c.GetStatus() == ClickStatusPause {
		c.resume <- true
	}
	return c.setStatus(ClickStatusStop)
}

func (c *Click) sendMsg(msg string, msgType MsgType) {
	runtime.EventsEmit(c.ctx, EventMessage, &Message{Msg: msg, Type: msgType})
}

func (c *Click) StartClick(points []*Point, params *Params) {
	if len(points) == 0 {
		c.sendMsg("no points", MsgError)
		return
	}
	if params.Times <= 0 {
		c.sendMsg("loop times is 0", MsgError)
		return
	}
	if !c.device.Status {
		c.sendMsg("device invalid, please check", MsgError)
		return
	}
	if c.GetStatus() != ClickStatusStop {
		c.sendMsg("already running", MsgError)
		return
	}
	c.setStatus(ClickStatusRunning)
	defer func() {
		c.setStatus(ClickStatusStop)
		runtime.EventsEmit(c.ctx, EventLoopNum, -1)
		c.sendMsg("Click finish ~", MsgSuccess)
	}()
	c.sendMsg("Click start ~", MsgSuccess)
	for i := 0; i < params.Times; i++ {
		switch c.GetStatus() {
		case ClickStatusPause:
			<-c.resume
		case ClickStatusStop:
			return
		}
		runtime.EventsEmit(c.ctx, EventLoopNum, i+1)
		var err error
		switch c.device.Type {
		case DeviceTypeAndroid:
			err = c.clickAndroid(points, params)
		case DeviceTypeIOS:
			err = c.clickIOS(points, params)
		default:
			err = fmt.Errorf("device not support")
		}
		if err != nil {
			c.sendMsg(fmt.Sprintf("Run click error: %s", err.Error()), MsgError)
		}
	}
}

func (c *Click) clickAndroid(points []*Point, params *Params) error {
	var err error
	var cmd string
	lastSwipeX := 0
	lastSwipeY := 0
	dbClick := false

	for _, p := range points {
		switch p.Event {
		case "Key":
			if p.X != -1 {
				continue
			}
			// Key Press
			// Power-26 Menu-82 HOME-3 Back-4 Vol+-24 Vol--25 Mute-164
			cmd = fmt.Sprintf("input keyevent %d", p.Y)
		case "Click", "DbClick":
			cmd = fmt.Sprintf("input tap %d %d", p.X, p.Y)
			if p.Event == "DbClick" {
				dbClick = true
			}
		case "LongPress":
			cmd = fmt.Sprintf("input swipe %d %d %d %d %d", p.X, p.Y, p.X+1, p.Y+1, params.Duration)
		case "Swipe", "QuickSwipe", "SlowSwipe":
			if lastSwipeX <= 0 || lastSwipeY <= 0 {
				lastSwipeX = p.X
				lastSwipeY = p.Y
				continue
			}
			dms := params.Duration
			switch p.Event {
			case "QuickSwipe":
				dms = 100
			case "SlowSwipe":
				dms = 1000
			}
			cmd = fmt.Sprintf("input swipe %d %d %d %d %d", lastSwipeX, lastSwipeY, p.X, p.Y, dms)
		case "StopApp":
			if p.X != -9 {
				continue
			}
			cmd = fmt.Sprintf("am force-stop %s", params.Package)
		case "StartApp":
			if p.X != -9 {
				continue
			}
			cmd = fmt.Sprintf("am start -n %s/%s", params.Package, params.Activity)
		default:
			if p.X == -1 {
				cmd = fmt.Sprintf("input keyevent %d", p.Y)
			} else {
				cmd = fmt.Sprintf("input tap %d %d", p.X, p.Y)
			}
		}

		_, err = c.device.adbDev.RunShellCommand(cmd)
		if dbClick {
			time.Sleep(50 * time.Millisecond)
			_, err = c.device.adbDev.RunShellCommand(cmd)
		}
		// reset
		dbClick = false
		lastSwipeX = 0
		lastSwipeY = 0
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(params.Interval) * time.Millisecond)
	}
	return nil
}

func (c *Click) clickIOS(points []*Point, params *Params) error {
	var err error
	lastSwipeX := 0
	lastSwipeY := 0
	for _, p := range points {
		switch p.Event {
		case "Key":
			if p.X != -1 {
				continue
			}
			key, ok := iosKeyMap[p.Y]
			if !ok {
				return fmt.Errorf("key not support")
			}
			err = c.device.iosDriver.PressButton(key)
		case "Click":
			err = c.device.iosDriver.Tap(p.X, p.Y)
		case "DbClick":
			err = c.device.iosDriver.DoubleTap(p.X, p.Y)
		case "LongPress":
			err = c.device.iosDriver.TouchAndHold(p.X, p.Y)
		case "Swipe":
			if lastSwipeX <= 0 || lastSwipeY <= 0 {
				lastSwipeX = p.X
				lastSwipeY = p.Y
				continue
			}
			err = c.device.iosDriver.Swipe(lastSwipeX, lastSwipeY, p.X, p.Y)
		case "Drag":
			if lastSwipeX <= 0 || lastSwipeY <= 0 {
				lastSwipeX = p.X
				lastSwipeY = p.Y
				continue
			}
			err = c.device.iosDriver.Drag(lastSwipeX, lastSwipeY, p.X, p.Y, 0.3)
		case "StopApp":
			if p.X != -9 {
				continue
			}
			_, err = c.device.iosDriver.AppTerminate(params.Package)
		case "StartApp":
			if p.X != -9 {
				continue
			}
			err = c.device.iosDriver.AppActivate(params.Package)
		default:
			err = c.device.iosDriver.Tap(p.X, p.Y)
		}

		// 重置临时变量
		lastSwipeX = 0
		lastSwipeY = 0
		if err != nil {
			return err
		}
		time.Sleep(time.Duration(params.Interval) * time.Millisecond)
	}
	return nil
}
