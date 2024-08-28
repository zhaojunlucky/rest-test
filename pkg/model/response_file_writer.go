package model

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func writeBodyRaw(ctx *core.RestTestContext, caseDef *TestCaseDef, resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("failed to read response body for test case %s-%s: %s", caseDef.GetID(), caseDef.GetFullDesc(), err.Error())
	}
	return writeBody(ctx, caseDef, resp, string(body))
}

func writeBody(ctx *core.RestTestContext, caseDef *TestCaseDef, resp *http.Response, str string) error {
	bodyFile := caseDef.GetFullDesc()
	bodyFile = fmt.Sprintf("%s_%s_response.txt", caseDef.GetID(), bodyFile)
	bodyFile = filepath.Join(ctx.LogPath, bodyFile)
	bodyFile = filepath.Clean(bodyFile)
	log.Infof("write test case response body to file: %s", bodyFile)

	fi, err := os.Create(bodyFile)
	if err != nil {
		log.Errorf("create file %s error: %s", bodyFile, err.Error())
		return err
	}
	defer func(fi *os.File) {
		err := fi.Close()
		if err != nil {
			log.Errorf("close file %s error: %s", bodyFile, err.Error())
		}
	}(fi)

	_, err = io.WriteString(fi, fmt.Sprintf("Status: %d\n", resp.StatusCode))

	_, err = io.WriteString(fi, "\nHeaders:\n")
	if err != nil {
		return err
	}

	for k, v := range resp.Header {
		_, err = io.WriteString(fi, fmt.Sprintf("%s: %s\n", k, strings.Join(v, ",")))
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(fi, "\nBody:\n")

	_, err = io.WriteString(fi, str)
	if err != nil {
		log.Errorf("write file %s error: %s", bodyFile, err.Error())
		return err
	}
	return nil
}
