package report

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/rest-test/pkg/model"
	"gopkg.in/yaml.v3"
	"os"
)

type TestPlanReport struct {
	ExecutionTime     float64
	TotalTime         float64
	TestPlan          *model.TestPlanDef
	testSuitesReports []*TestSuiteReport
	Error             error
	Status            string
}

func (t *TestPlanReport) GetTestDef() *model.TestPlanDef {
	return t.TestPlan
}

func (t *TestPlanReport) GetExecutionTime() float64 {
	return t.ExecutionTime
}

func (t *TestPlanReport) GetTotalTime() float64 {
	return t.TotalTime
}

func (t *TestPlanReport) GetChildren() []*TestSuiteReport {
	return t.testSuitesReports
}

func (t *TestPlanReport) AddTestSuiteReport(testSuiteReport *TestSuiteReport) {
	t.testSuitesReports = append(t.testSuitesReports, testSuiteReport)
}

func (t *TestPlanReport) GetError() error {
	var suiteErrs []error
	if t.Error != nil {
		suiteErrs = append(suiteErrs, t.Error)
	}

	for _, testSuiteReport := range t.testSuitesReports {
		if err := testSuiteReport.GetError(); err != nil {
			suiteErrs = append(suiteErrs, err)
		}
	}
	return errors.Join(suiteErrs...)
}

func (t *TestPlanReport) GetStatus() string {
	return t.Status
}

func (t *TestPlanReport) GetReportData() (map[string]any, error) {
	if t.Status == "" {
		return nil, fmt.Errorf("test plan %s name is not executed", t.TestPlan.Name)
	}
	var suiteDataList []map[string]any

	for _, suiteReport := range t.testSuitesReports {
		suiteData, err := suiteReport.GetReportData()
		if err != nil {
			return nil, err
		}
		suiteDataList = append(suiteDataList, suiteData)
	}
	return map[string]any{
		"type":          "plan",
		"name":          t.TestPlan.Name,
		"executionTime": t.ExecutionTime,
		"totalTime":     t.TotalTime,
		"error":         getErrorStr(t.Error),
		"status":        t.Status,
		"suites":        suiteDataList,
	}, nil
}

func (t *TestPlanReport) WriteReport(file string) error {
	log.Infof("write test plan report to file: %s", file)

	data, err := t.GetReportData()
	if err != nil {
		return err
	}
	out, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}
	fi, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func(fi *os.File) {
		err := fi.Close()
		if err != nil {
			log.Errorf("close file %s error: %s", file, err.Error())
		}
	}(fi)
	_, err = fi.Write(out)
	fmt.Println(string(out))

	return err
}
