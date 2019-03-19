package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opensatelliteproject/SatHelperApp"
	"github.com/opensatelliteproject/SatHelperApp/RPC/sathelperapp"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

func GetStatistics(client sathelperapp.InformationClient) {
	stats, err := client.GetStatistics(context.Background(), &empty.Empty{})
	if err != nil {
		fmt.Printf("Error fetching statistics: %s", err)
		os.Exit(1)
	}

	data, err := json.MarshalIndent(stats, "", "   ")
	if err != nil {
		fmt.Printf("Error serialzing data: %s", err)
	}

	fmt.Println(string(data))
}

func GetConsoleLines(client sathelperapp.InformationClient) {
	lines, err := client.GetConsoleLines(context.Background(), &empty.Empty{})
	if err != nil {
		fmt.Printf("Error fetching statistics: %s", err)
		os.Exit(1)
	}

	for _, v := range lines.ConsoleLines {
		fmt.Println(v)
	}
}

func main() {
	kingpin.Version(SatHelperApp.GetVersion())

	server := kingpin.Flag("server", "Server Address").Default("127.0.0.1:5500").String()
	method := kingpin.Arg("method", "Method").Required().String()
	kingpin.Parse()

	conn, err := grpc.Dial(*server, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Error connectint to %s: %s", *server, err)
		os.Exit(1)
	}
	defer conn.Close()

	client := sathelperapp.NewInformationClient(conn)

	switch *method {
	case "GetStatistics":
		GetStatistics(client)
	case "GetConsoleLines":
		GetConsoleLines(client)
	default:
		fmt.Printf("Unknown method: %s\n", *method)
		os.Exit(1)
	}
}
