package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/golib/pkg/env"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"github.com/zhaojunlucky/rest-test/pkg/executor"
	"github.com/zhaojunlucky/rest-test/pkg/model"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"
)

func setupLog(ctx *core.RestTestContext, logPath, logLevel string) error {
	if runtime.GOOS == "windows" {
		panic("Windows is currently not supported.")
	}

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		return fmt.Errorf("invalid log level: %s", logLevel)
	}

	if len(logPath) <= 0 {
		logPath = "/var/log/rest_test"
	}
	t := time.Now()

	// Format the time using a layout string
	formattedTime := t.Format("2006-01-02_15-04-05")

	logPath = path.Join(logPath, formattedTime)
	log.Infof("log path: %s", logPath)

	fiInfo, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(logPath, 0755)
		if err != nil {
			panic(err)
		}
	} else if !fiInfo.IsDir() {
		return fmt.Errorf("%s must be a directory", logPath)
	}

	ctx.LogPath = logPath

	logFilePath := path.Join(logPath, "rest_test.log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Info("Failed to log to file, using default stderr")
		return nil
	}
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			//return frame.Function, fileName
			return "", fileName
		},
	})

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	return nil
}

func main() {

	planPtr := flag.String("plan", "", "a test plan")
	suitePtr := flag.String("suite", "", "a test suite")
	logLevel := flag.String("level", "info", "log level. levels: debug, info(default), warn, error")
	logPath := flag.String("logPath", "", "log path")

	flag.Parse()
	ctx := &core.RestTestContext{}

	if err := setupLog(ctx, *logPath, *logLevel); err != nil {
		log.Fatal(err)
	}

	if len(*planPtr) <= 0 && len(*suitePtr) <= 0 {
		flag.Usage()
		return
	} else if len(*planPtr) > 0 && len(*suitePtr) > 0 {
		log.Error("only one of plan or suite can be specified")
		return
	} else if len(*planPtr) > 0 {
		executePlan(ctx, *planPtr)
	} else {
		executeSuite(ctx, *suitePtr)
	}
}

func executeSuite(ctx *core.RestTestContext, s string) {
	log.Infof("execute suite: %s", s)
	testSuiteDef := model.TestSuiteDef{}
	err := testSuiteDef.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	testSuiteExecutor := executor.NewTestSuiteExecutor()

	report, err := testSuiteExecutor.ExecuteSuite(ctx, env.NewOSEnv(), &testSuiteDef)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("suite report: %s\n", report.TestSuite.Name)
	fmt.Printf("status: %s\n", report.Status)
	failed := false
	if report.Error != nil {
		fmt.Printf("error: %s\n", report.Error)
		failed = true
	}

	for i, caseReport := range report.GetChildren() {
		fmt.Printf("  %d case report: %s - %s\n", i+1, caseReport.TestCase.Description, caseReport.TestCase.Name)
		fmt.Printf("  status: %s\n", caseReport.Status)
		if caseReport.Error != nil {
			fmt.Printf("  error: %s\n", caseReport.Error)
			failed = true
		}
	}
	if failed {
		os.Exit(1)
	}
}

func executePlan(ctx *core.RestTestContext, s string) {
	log.Infof("execute plan: %s", s)
	testPlanDef := model.TestPlanDef{}
	err := testPlanDef.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	testPlanExecutor := executor.NewTestPlanExecutor()
	report, err := testPlanExecutor.ExecutePlan(ctx, env.NewOSEnv(), &testPlanDef)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("plan report: %s\n", testPlanDef.Name)
	fmt.Printf("plan status: %s\n", report.Status)
	failed := false

	for i, suiteReport := range report.GetChildren() {
		fmt.Printf("\t%d. suite report: %s\n", i+1, suiteReport.TestSuite.Name)
		fmt.Printf("\tstatus: %s\n", suiteReport.Status)
		if report.Error != nil {
			fmt.Printf("\terror: %s\n", suiteReport.Error)
			failed = true
		}

		for j, caseReport := range suiteReport.GetChildren() {
			fmt.Printf("\t\t%d case report: %s - %s\n", j+1, caseReport.TestCase.Description, caseReport.TestCase.Name)
			fmt.Printf("\t\tstatus: %s\n", caseReport.Status)
			if caseReport.Error != nil {
				fmt.Printf("\t\terror: %s\n", caseReport.Error)
				failed = true
			}
		}
	}

	if failed {
		os.Exit(1)
	}

}
