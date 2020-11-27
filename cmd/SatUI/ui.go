package main

import (
	"encoding/json"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
	"github.com/gorilla/mux"
	"github.com/opensatelliteproject/SatHelperApp"
	"github.com/webview/webview"
	"net"
	"net/http"
	"time"
)

var astiWindow *astilectron.Window
var webviewWindow webview.WebView
var tcpPort int
var httpListen net.Listener

var boundFuncs = map[string]interface{}{
	"SatHelperApp_SaveConfig":     saveConfig,
	"SatHelperApp_IsConfigLoaded": isConfigLoaded,
	"SatHelperApp_GetConfig":      getConfig,
	"SatHelperApp_SetConfig":      setConfig,
	"SatHelperApp_LoadConfig":     loadConfig,
	"SatHelperApp_StartServer":    startServerApp,
	"SatHelperApp_StopServer":     stopServerApp,
	"SatHelperApp_Exit":           exit,
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "explore":
		// Unmarshal payload
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}
	}
	return
}

func exit() {
	_ = stopServerApp()
	if httpListen != nil {
		_ = httpListen.Close()
	}
}

func startServerApp(function string) error {
	go func() {
		time.Sleep(time.Second * 5)
		webviewWindow.Dispatch(func() {
			webviewWindow.Eval(fmt.Sprintf("%s()", function))
		})
	}()
	return nil
}

func stopServerApp() error {
	time.Sleep(time.Second * 2)
	return nil
}

func initView() error {
	var err error
	webviewWindow.SetTitle(fmt.Sprintf("SatUI - %s", SatHelperApp.GetVersion()))
	webviewWindow.SetSize(1000, 600, webview.HintNone)
	webviewWindow.Navigate(fmt.Sprintf("http://127.0.0.1:%d/index.html", tcpPort))
	for funcName, function := range boundFuncs {
		fmt.Printf("Binding function %s\n", funcName)
		err = webviewWindow.Bind(funcName, function)
		if err != nil {
			return err
		}
	}
	return err
}

func webserv(c chan bool) error {
	var err error
	defer func() {
		c <- false
	}()

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./site"))
	r.PathPrefix("/").Handler(fs)

	httpListen, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}

	tcpPort = httpListen.Addr().(*net.TCPAddr).Port
	c <- true
	fmt.Printf("Webserver running at 127.0.0.1:%d\n", tcpPort)
	return http.Serve(httpListen, r)
}

func runWebkit(c chan bool) {
	fmt.Println("Starting webkit")
	c <- true
	webviewWindow.Run()
	c <- false
}

func runAstilectron() error {
	err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName: "SatUI",
		},
		Debug: true,
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			astiWindow = ws[0]
			go func() {
				time.Sleep(5 * time.Second)
				if err := bootstrap.SendMessage(astiWindow, "FODACI", "FODACI O ASTIELETROCON"); err != nil {
					fmt.Println(fmt.Errorf("sending check.out.menu event failed: %w", err))
				}
			}()
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
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#333"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(1000),
				Width:           astikit.IntPtr(800),
			},
		}},
	})

	return err
}

func startUI(pctx *kong.Context) error {
	webChan := make(chan bool, 1)
	//kitChan := make(chan bool, 1)
	// Starting server
	go webserv(webChan)

	<-webChan // Wait server to start

	//debug := true
	//webviewWindow = webview.New(debug)
	//defer webviewWindow.Destroy()
	//err := initView()
	//if err != nil {
	//	return err
	//}

	go func() {
		fmt.Println("Running...")

		select {
		case <-webChan:
			fmt.Println("HTTP Server Closed")
		}

		if httpListen != nil {
			_ = httpListen.Close()
		}
		if webviewWindow != nil {
			webviewWindow.Terminate()
		}
	}()

	err := runAstilectron()

	if err != nil {
		if httpListen != nil {
			_ = httpListen.Close()
		}
		return err
	}

	return nil
}
