package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/prometheus/common/log"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"sync"
	"syscall"
	"time"
)

var astiWindow *astilectron.Window
var boundFuncs = map[string]interface{}{
	"SatHelperApp_SaveConfig":      saveConfig,
	"SatHelperApp_IsConfigLoaded":  isConfigLoaded,
	"SatHelperApp_GetConfig":       getConfig,
	"SatHelperApp_SetConfig":       setConfig,
	"SatHelperApp_LoadConfig":      loadConfig,
	"SatHelperApp_StartServer":     startServerApp,
	"SatHelperApp_StopServer":      stopServerApp,
	"SatHelperApp_ServerIsRunning": isRunning,
	"SatHelperApp_Exit":            exit,
}

func reflectCall(f interface{}, data json.RawMessage) (v interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %s", r)
			log.Errorf("error: %s", r)
			debug.PrintStack()
		}
	}()
	fType := reflect.TypeOf(f)
	if fType.Kind() != reflect.Func {
		return nil, fmt.Errorf("reflectCall called with a non function interface")
	}

	numIn := fType.NumIn()
	fmethod := reflect.Indirect(reflect.ValueOf(f))

	in := make([]reflect.Value, numIn)

	var params []interface{}

	err = json.Unmarshal(data, &params)

	if err != nil {
		return nil, err
	}

	for i := 0; i < numIn; i++ {
		t := fType.In(i)
		if len(params) > i {
			paramValue := params[i]
			if t.Kind() == reflect.Struct {
				// Re-marshal
				d, _ := json.Marshal(paramValue)
				n := reflect.New(t).Interface()
				err = json.Unmarshal(d, n)
				if err != nil {
					return nil, err
				}
				in[i] = reflect.Indirect(reflect.ValueOf(n))
			} else {
				in[i] = reflect.ValueOf(paramValue)
			}
		} else {
			in[i] = reflect.New(t).Elem()
		}
	}

	out := fmethod.Call(in)
	if len(out) == 0 {
		return nil, nil
	}

	lastParam := out[len(out)-1]

	if lastParam.Type() == reflect.TypeOf((*error)(nil)).Elem() && lastParam.Interface() != nil {
		out = out[:len(out)-1]
		if !lastParam.IsNil() {
			err = *lastParam.Interface().(*error)
		}
	}

	if len(out) == 1 {
		return out[0].Interface(), err
	}

	interfaceArray := make([]interface{}, len(out))
	for i, v := range out {
		interfaceArray[i] = v.Interface()
	}

	return interfaceArray, err
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	if function, ok := boundFuncs[m.Name]; ok {
		return reflectCall(function, m.Payload)
	}

	return
}

func exit() {
	_ = stopServerApp()
	if astiWindow != nil {
		_ = astiWindow.Close()
	}
	os.Exit(0)
}

var proc *exec.Cmd
var procMtx = sync.Mutex{}

func getAppPath() (string, error) {
	return filepath.Abs(os.Args[0])
}

func sendBuffer(stdname string, data []byte) error {
	return astiWindow.SendMessage(bootstrap.MessageOut{
		Name:    stdname,
		Payload: data,
	})
}

func stdioLoop(stdname string, stdioReader io.ReadCloser) {
	log.Infof("Starting %s readers", stdname)
	stdio := bufio.NewReader(stdioReader)
	lastFlush := time.Now()

	buffer := make([]byte, 4096)
	currentLength := 0

	for {
		b, err := stdio.ReadByte()
		if err != nil {
			if err != io.EOF {
				log.Errorf("error reading line: %s", err)
			}
			break
		}
		buffer[currentLength] = b
		currentLength++
		if currentLength == len(buffer) ||
			(time.Since(lastFlush) > time.Millisecond*10 && currentLength > 0) ||
			b == '\n' {
			err = sendBuffer(stdname, buffer[:currentLength])
			if err != nil {
				log.Errorf("error sending buffer: %s", err)
				break
			}
			currentLength = 0
			lastFlush = time.Now()
		}
	}
	log.Info("Ending readers")
}

func isRunning() bool {
	return proc != nil && proc.ProcessState == nil
}

func startServerApp() error {
	procMtx.Lock()
	defer procMtx.Unlock()
	if isRunning() {
		return fmt.Errorf("already running")
	}

	appPath, err := getAppPath()
	if err != nil {
		return err
	}

	log.Infof("Executing %s server", appPath)

	proc = exec.Command(appPath, "server")

	stderrReader, err := proc.StderrPipe()
	if err != nil {
		proc = nil
		return err
	}
	stdoutReader, err := proc.StdoutPipe()
	if err != nil {
		proc = nil
		return err
	}

	go stdioLoop("serverStdout", stdoutReader)
	go stdioLoop("serverStderr", stderrReader)

	err = proc.Start()
	if err != nil {
		proc = nil
		return err
	}

	log.Info("Process started")

	return nil
}

func stopServerApp() error {
	procMtx.Lock()
	defer procMtx.Unlock()

	if isRunning() {
		log.Info("Stopping process")
		_ = proc.Process.Signal(syscall.SIGTERM)
		_, _ = proc.Process.Wait()
		proc = nil
		log.Info("Process stopped")
	}

	return nil
}

func runAstilectron() error {
	mydir, _ := os.Getwd()
	err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:           "SatUI",
			DataDirectoryPath: mydir,
		},
		Debug: true,
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			astiWindow = ws[0]
			return nil
		},
		MenuOptions: []*astilectron.MenuItemOptions{
			{
				Label: astikit.StrPtr("File"),
				SubMenu: []*astilectron.MenuItemOptions{
					{Role: astilectron.MenuItemRoleClose},
				},
			},
			{
				Label: astikit.StrPtr("Tools"),
				SubMenu: []*astilectron.MenuItemOptions{
					{
						Accelerator: astilectron.NewAccelerator("Alt", "CommandOrControl", "I"),
						Role:        astilectron.MenuItemRoleToggleDevTools,
					},
				},
			},
			{
				Label: astikit.StrPtr("Edit"),
				SubMenu: []*astilectron.MenuItemOptions{
					{
						Label:       astikit.StrPtr("Undo"),
						Accelerator: astilectron.NewAccelerator("CmdOrCtrl", "Z"),
						Role:        astikit.StrPtr("undo:"),
					},
					{
						Label:       astikit.StrPtr("Redo"),
						Accelerator: astilectron.NewAccelerator("Shift", "CmdOrCtrl", "Z"),
						Role:        astikit.StrPtr("redo"),
					},
					{
						Type: astikit.StrPtr("separator"),
					},
					{
						Label:       astikit.StrPtr("Cut"),
						Accelerator: astilectron.NewAccelerator("CmdOrCtrl", "X"),
						Role:        astikit.StrPtr("cut"),
					},
					{
						Label:       astikit.StrPtr("Copy"),
						Accelerator: astilectron.NewAccelerator("CmdOrCtrl", "C"),
						Role:        astikit.StrPtr("copy"),
					},
					{
						Label:       astikit.StrPtr("Paste"),
						Accelerator: astilectron.NewAccelerator("CmdOrCtrl", "V"),
						Role:        astikit.StrPtr("paste"),
					},
					{
						Label:       astikit.StrPtr("Select All"),
						Accelerator: astilectron.NewAccelerator("CmdOrCtrl", "A"),
						Role:        astikit.StrPtr("selectAll"),
					},
				},
			},
		},
		//ResourcesPath: "/resources",
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#333"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(640),
				Width:           astikit.IntPtr(1000),
			},
		}},
	})

	return err
}

func startUI(pctx *kong.Context) error {
	return runAstilectron()
}
