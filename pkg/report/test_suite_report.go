package report

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/rest-test/pkg/model"
	"gopkg.in/yaml.v3"
	"os"
)

type TestSuiteReport struct {
	ExecutionTime   float64
	TotalTime       float64
	TestSuite       *model.TestSuiteDef
	Error           error
	testCaseReports []*TestCaseReport
	Status          string
}

func (t *TestSuiteReport) GetTestDef() *model.TestSuiteDef {
	return t.TestSuite
}

func (t *TestSuiteReport) GetExecutionTime() float64 {
	return t.ExecutionTime
}

func (t *TestSuiteReport) GetTotalTime() float64 {
	return t.TotalTime
}

func (t *TestSuiteReport) GetChildren() []*TestCaseReport {
	return t.testCaseReports
}

func (t *TestSuiteReport) AddTestCaseReport(testCaseReport *TestCaseReport) {
	t.testCaseReports = append(t.testCaseReports, testCaseReport)
}

func (t *TestSuiteReport) GetError() error {
	var caseErrs []error
	if t.Error != nil {
		caseErrs = append(caseErrs, t.Error)
	}

	for _, testCaseReport := range t.testCaseReports {
		if err := testCaseReport.GetError(); err != nil {
			caseErrs = append(caseErrs, err)
		}
	}
	return errors.Join(caseErrs...)
}

func (t *TestSuiteReport) GetStatus() string {
	return t.Status
}

func (t *TestSuiteReport) HasPassed() bool {
	for _, testCaseReport := range t.testCaseReports {
		if testCaseReport.GetStatus() != Completed || testCaseReport.GetError() != nil {
			return false
		}
	}
	return true
}

func (t *TestSuiteReport) GetReportData() (map[string]any, error) {
	if t.Status == "" {
		return nil, fmt.Errorf("test suite %s name is not executed", t.TestSuite.Name)
	}
	var caseDataList []map[string]any

	for _, testCaseReport := range t.testCaseReports {
		caseData, err := testCaseReport.GetReportData()
		if err != nil {
			return nil, err
		}
		caseDataList = append(caseDataList, caseData)
	}
	return map[string]any{
		"type":          "suite",
		"name":          t.TestSuite.Name,
		"executionTime": t.ExecutionTime,
		"totalTime":     t.TotalTime,
		"error":         getErrorStr(t.Error),
		"status":        t.Status,
		"cases":         caseDataList,
	}, nil
}

func (t *TestSuiteReport) WriteReport(file string) error {
	log.Infof("write test suite report to file: %s", file)

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
