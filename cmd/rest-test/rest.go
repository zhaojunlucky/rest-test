package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/golib/pkg/env"
	"github.com/zhaojunlucky/rest-test/pkg/executor"
	"github.com/zhaojunlucky/rest-test/pkg/model"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
)

func setupLog() error {
	if runtime.GOOS == "windows" {
		panic("Windows is currently not supported.")
	}
	logPath := "/var/log/rest_test"
	fiInfo, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(logPath, 0755)
		if err != nil {
			panic(err)
		}
	} else if !fiInfo.IsDir() {
		panic(fmt.Sprintf("%s must be a directory.", logPath))
	}

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
	if err := setupLog(); err != nil {
		panic(err)
	}

	planPtr := flag.String("plan", "", "a test plan")
	suitePtr := flag.String("suite", "", "a test suite")

	flag.Parse()

	if len(*planPtr) <= 0 && len(*suitePtr) <= 0 {
		flag.Usage()
		return
	} else if len(*planPtr) > 0 && len(*suitePtr) > 0 {
		log.Error("only one of plan or suite can be specified")
		return
	} else if len(*planPtr) > 0 {
		executePlan(*planPtr)
	} else {
		executeSuite(*suitePtr)
	}
}

func executeSuite(s string) {
	testSuiteDef := model.TestSuiteDef{}
	err := testSuiteDef.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	testSuiteExecutor := executor.NewTestSuiteExecutor()

	report, err := testSuiteExecutor.ExecuteSuite(env.NewOSEnv(), &testSuiteDef)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("suite report: %s\n", report.TestSuite.Name)
	fmt.Printf("status: %s\n", report.Status)
	if report.Error != nil {
		fmt.Printf("error: %s\n", report.Error)
	}

	for i, caseReport := range report.GetChildren() {
		fmt.Printf("  %d case report: %s\n", i+1, caseReport.TestCase.Name)
		fmt.Printf("  status: %s\n", caseReport.Status)
		if caseReport.Error != nil {
			fmt.Printf("  error: %s\n", caseReport.Error)
		}
	}
}

func executePlan(s string) {
	testPlanDef := model.TestPlanDef{}
	err := testPlanDef.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	testPlanExecutor := executor.NewTestPlanExecutor()
	_, err = testPlanExecutor.ExecutePlan(env.NewOSEnv(), &testPlanDef)
	if err != nil {
		log.Fatal(err)
	}
}
