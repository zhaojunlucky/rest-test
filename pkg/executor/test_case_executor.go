package executor

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/golib/pkg/env"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"github.com/zhaojunlucky/rest-test/pkg/execution"
	"github.com/zhaojunlucky/rest-test/pkg/model"
	"github.com/zhaojunlucky/rest-test/pkg/report"
	"net/http"
	"strings"
	"time"
)

type TestCaseExecutor struct {
}

func (t *TestCaseExecutor) Execute(ctx *core.RestTestContext, env env.Env, global *model.GlobalSetting, testCaseExecResult *execution.TestCaseExecutionResult,
	testSuiteCase *TestSuiteCaseContext) *report.TestCaseReport {
	defer func() {
		log.Infof("[Case] end execute test case: %s-%s", testCaseExecResult.TestCaseDef.Description, testCaseExecResult.TestCaseDef.Name)
	}()
	log.Infof("[Case] start execute test case: %s-%s", testCaseExecResult.TestCaseDef.Description, testCaseExecResult.TestCaseDef.Name)
	testCaseExecResult.Executed = true

	testCaseReport := &report.TestCaseReport{
		TestCase: testCaseExecResult.TestCaseDef,
	}
	testCaseExecResult.TestCaseReport = testCaseReport

	testCaseDef := testCaseExecResult.TestCaseDef

	if !testCaseDef.Enabled {
		testCaseReport.Status = report.Skipped
		return testCaseReport
	}

	env.SetAll(testCaseDef.Environment)

	js, err := NewJSScriptler(env, testSuiteCase)
	if err != nil {
		log.Errorf("failed to create js scriptler: %v", err)
		testCaseReport.Status = report.InitError
		testCaseReport.Error = err
		return testCaseReport
	}

	start := time.Now()
	global, err = global.Expand(js)
	if err != nil {
		log.Errorf("failed to expand global: %v", err)
		testCaseReport.Error = err
		testCaseReport.Status = report.InitError
		return testCaseReport
	}

	resp, err := t.performHTTPRequest(ctx, global, testCaseDef, js)
	if err != nil {
		log.Errorf("failed to perform http request: %v", err)
		testCaseReport.Error = err
		testCaseReport.Status = report.ExecutionError
		return testCaseReport
	}
	body, err := testCaseExecResult.TestCaseDef.Response.Validate(ctx, resp, js)
	if err != nil {
		log.Errorf("failed to validate response: %v", err)
		testCaseReport.Error = err
		testCaseReport.Status = report.ExecutionError
		return testCaseReport
	}
	testCaseReport.TotalTime = time.Since(start).Seconds()
	testCaseReport.Status = report.Completed
	testCaseReport.Error = nil

	testCaseReport.ExecutionTime = time.Since(start).Seconds()

	err = testSuiteCase.Add(testCaseExecResult, body)
	if err != nil {
		log.Errorf("failed to add test case result: %v", err)
		log.Error(err)
	}
	return testCaseReport
}

func (t *TestCaseExecutor) Prepare(ctx *execution.TestSuiteExecutionResult, def *model.TestCaseDef) error {
	defer func() {
		log.Infof("[Case] end prepare test case: %s-%s", def.Description, def.Name)
	}()
	log.Infof("[Case] prepare test case: %s-%s", def.Description, def.Name)
	if ctx.HasNamed(def.Name) {
		return fmt.Errorf("duplicated named test case %s", def.Name)
	}

	testCaseExecResult := &execution.TestCaseExecutionResult{
		TestCaseDef:              def,
		TestSuiteExecutionResult: ctx,
	}
	ctx.AddTestCaseExecResults(testCaseExecResult)
	return nil
}

func (t *TestCaseExecutor) Validate(result *execution.TestCaseExecutionResult) error {

	return nil
}

func (t *TestCaseExecutor) performHTTPRequest(ctx *core.RestTestContext, global *model.GlobalSetting, def *model.TestCaseDef, js *JSScriptler) (*http.Response, error) {
	url, err := js.Expand(def.Request.URL)
	if err != nil {
		log.Infof("failed to expand url: %v", err)
		return nil, err
	}

	if !strings.HasPrefix(url, "http") {
		if !strings.HasPrefix(global.APIPrefix, "http") {
			return nil, fmt.Errorf("apiPrefix or url must be http or https")
		}
		url = fmt.Sprintf("%s/%s", global.APIPrefix, url)
	}

	bodyReader, err := def.Request.Body.GetBody(global.DataDir, js)
	if err != nil {
		log.Errorf("failed to get body: %v", err)
		return nil, err
	}
	var req *http.Request
	req, err = http.NewRequest(def.Request.Method, url, bodyReader)
	if err != nil {
		log.Errorf("failed to create request: %v", err)
		return nil, err
	}

	for k, v := range global.Headers {
		req.Header.Add(k, v)
	}
	for k, v := range def.Request.Parameters {
		req.URL.Query().Add(k, v)
	}
	for k, v := range def.Request.Headers {
		req.Header.Add(k, v)
	}
	return http.DefaultClient.Do(req)
}

func NewTestCaseExecutor() *TestCaseExecutor {
	return &TestCaseExecutor{}
}
