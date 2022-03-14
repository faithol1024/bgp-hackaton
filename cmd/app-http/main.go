package main

import (
	"flag"
	"net"
	"os"

	"github.com/faithol1024/bgp-hackathon/internal/config"
	"github.com/tokopedia/campaign-engine/common/util"
	"github.com/tokopedia/tdk/go/log"
)

const (
	repoName = "bgp-hackathon"
)

func main() {

	initLog()

	// init config
	cfg, err := config.New(repoName)
	if err != nil {
		log.Fatalf("failed to init the config: %v", err)
	}

	// uncomment after slack integration
	// // set panic option
	// if !env.IsDevelopment() {
	// 	svcEnv := string(env.ServiceEnv())
	// 	hostName, hostIP := getHostInfo()
	// 	panics.SetOptions(&panics.Options{
	// 		Env:                        svcEnv,
	// 		SlackWebhookURL:            "myslackwebhook",
	// 		SlackChannel:               "mychannel",
	// 		Tags:                       map[string]string{"service": repoName, "env": svcEnv, "hostname": hostName, "hostip": hostIP},
	// 		AllowPublishDumpRequestLog: true,
	// 	})
	// }

	err = startApp(cfg)
	if err != nil {
		log.Fatalf("failed to start app: %v", err)
	}
}

func initLog() {

	logDir := "/var/log/campaign-engine/"
	if util.GetEnv() == "development" {
		logDir = "log/"
	}

	var (
		debugLogPath string
		infoLogPath  string
	)

	// parse flag
	flag.StringVar(&infoLogPath, "info_log", logDir, "path for info and error log file")
	flag.StringVar(&debugLogPath, "debug_log", logDir, "path for debug log level")
	flag.Parse()

	// init log
	log.SetConfig(&log.Config{
		LogFile:   infoLogPath,
		DebugFile: debugLogPath,
	})
}

func getHostInfo() (name, ip string) {
	var err error
	name, err = os.Hostname()
	if err != nil {
		name = "-"
	}

	ip = "-"
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return name, ip
	}
	defer conn.Close()

	localAddr, ok := conn.LocalAddr().(*net.UDPAddr)
	if ok {
		ip = localAddr.IP.String()
	}
	return name, ip
}
