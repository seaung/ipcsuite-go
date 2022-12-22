package camera

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	hk "github.com/seaung/ipcsuite-go/internal/http"
	"github.com/seaung/ipcsuite-go/internal/protos"
	"github.com/seaung/ipcsuite-go/pkg/utils"
)

type Auditor struct {
	Requests *http.Request
	NsePocs  *NsePoc
}

func auditVuln(auditors []Auditor, ticker *time.Ticker) <-chan Auditor {
	var wg sync.WaitGroup
	res := make(chan Auditor)

	for _, audit := range auditors {
		wg.Add(1)

		go func(audit Auditor) {
			defer wg.Done()
			<-ticker.C
			ok, err := runPoc(audit.Requests, audit.NsePocs)
			if err != nil {
				return
			}

			if ok {
				res <- audit
			}
		}(audit)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}

func runPoc(request *http.Request, poc *NsePoc) (bool, error) {
	celOptions := NewCelEnvOptions()
	celOptions.UpdateCompileOptions(poc.Set)

	env, err := NewEnvCelOption(&celOptions)
	if err != nil {
		return false, err
	}

	paramsMap := make(map[string]interface{})

	req, err := hk.ParseRequest(request)
	if err != nil {
		return false, err
	}

	paramsMap["request"] = req

	keys := make([]string, 0)

	for key := range poc.Set {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		expression := poc.Set[key]

		if key != "payload" {
			out, err := EvalExpression(env, expression, paramsMap)
			if err != nil {
				continue
			}

			switch value := out.Value().(type) {
			case *protos.UrlType:
				paramsMap[key] = UrlType2String(value)
			case int64:
				paramsMap[key] = int(value)
			default:
				paramsMap[key] = fmt.Sprintf("%v", out)
			}
		}
	}

	if poc.Set["payload"] != "" {
		out, err := EvalExpression(env, poc.Set["payload"], paramsMap)
		if err != nil {
			return false, err
		}

		paramsMap["payload"] = fmt.Sprintf("%v", out)
	}

	success := false

	for _, rule := range poc.Rule {
		for k1, v1 := range paramsMap {
			_, ok := v1.(map[string]string)

			if ok {
				continue
			}

			value := fmt.Sprintf("%v", v1)

			for k2, v2 := range rule.Headers {
				rule.Headers[k2] = strings.ReplaceAll(v2, "{{"+k1+"}}", value)
			}

			rule.Path = strings.ReplaceAll(strings.TrimSpace(rule.Path), "{{"+k1+"}}", value)
			rule.Body = strings.ReplaceAll(strings.TrimSpace(rule.Body), "{{"+k1+"}}", value)
		}

		if request.URL.Path != "" && request.URL.Path != "/" {
			req.Url.Path = fmt.Sprint(request.URL.Path, rule.Path)
		} else {
			req.Url.Path = rule.Path
		}

		req.Url.Path = strings.ReplaceAll(req.Url.Path, " ", "%20")
		req.Url.Path = strings.ReplaceAll(req.Url.Path, "+", "%20")

		client, _ := http.NewRequest(rule.Method, fmt.Sprintf("%s://%s%s", req.Url.Scheme, req.Url.Host, req.Url.Path), strings.NewReader(rule.Body))
		client.Header = request.Header.Clone()

		for key, value := range rule.Headers {
			client.Header.Set(key, value)
		}

		response, err := hk.SendRequest(client, rule.AllowRedirect)
		if err != nil {
			return false, err
		}

		paramsMap["response"] = response

		if rule.Matchs != "" {
			res := utils.FindMatch(strings.TrimSpace(rule.Matchs), string(response.Body))
			if res != nil && len(res) > 0 {
				for key, value := range res {
					paramsMap[key] = value
				}
			} else {
				return false, nil
			}
		}

		out, err := EvalExpression(env, rule.Expression, paramsMap)
		if err != nil {
			return false, err
		}

		if fmt.Sprintf("%v", out) == "false" {
			success = false
			break
		}

		success = true

	}

	return success, nil
}
