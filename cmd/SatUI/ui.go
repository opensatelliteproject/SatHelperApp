package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/gorilla/mux"
	"github.com/opensatelliteproject/SatHelperApp"
	"github.com/webview/webview"
	"net"
	"net/http"
)

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
}

func startServerApp() error {
	return nil
}

func stopServerApp() error {
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

func startUI(pctx *kong.Context) error {
	webChan := make(chan bool, 1)
	kitChan := make(chan bool, 1)
	// Starting server
	go webserv(webChan)

	<-webChan // Wait server to start

	debug := true
	webviewWindow = webview.New(debug)
	defer webviewWindow.Destroy()
	err := initView()
	if err != nil {
		return err
	}

	go func() {
		fmt.Println("Waiting webkit to start...")
		<-kitChan // Wait webview to start
		fmt.Println("Running...")

		select {
		case <-kitChan:
			fmt.Println("Webkit closed")
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

	runWebkit(kitChan)

	return nil
}
