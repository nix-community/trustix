// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package hydra

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type JobsetEvalInput struct {
	Revision string `json:"revision"`
	URI      string `json:"uri"`
	Type     string `json:"type"`
}

type HydraEval struct {
	ID         int64                       `json:"id"`
	Timestamp  int64                       `json:"timestamp"`
	EvalInputs map[string]*JobsetEvalInput `json:"jobsetevalinputs"`
}

type NixPath map[string]string

func (n *NixPath) String() string {
	var b strings.Builder

	nInputs := len(*n)

	i := 0
	for k, v := range *n {
		if i >= 1 && i <= nInputs-1 {
			b.WriteRune(':')
		}

		b.WriteString(k)
		b.WriteRune('=')
		b.WriteString(v)

		i++
	}

	return b.String()
}

// Get NIX_PATH from a Hydra evaluation
func (e *HydraEval) NixPath() (NixPath, error) {

	// We need resolved local paths, and the only reasonable way to download them is to feed them to a nix expression
	var expr string
	{
		var b strings.Builder

		b.WriteString("builtins.toJSON {")

		for k, v := range e.EvalInputs {
			switch v.Type {
			case "boolean":
				continue
			case "git":
			default:
				return nil, fmt.Errorf("unsupported hydra input type: %s", v.Type)
			}

			// Sanity check URI to prevent command injections
			_, err := url.Parse(v.URI)
			if err != nil {
				return nil, fmt.Errorf("error parsing URL: %w", err)
			}

			// Sanity check rev to prevent command injections
			_, err = hex.DecodeString(v.Revision)
			if err != nil {
				return nil, fmt.Errorf("error parsing revision: %w", err)
			}

			b.WriteString(fmt.Sprintf("  %s = builtins.fetchGit { url = \"%s\"; rev = \"%s\"; allRefs = true; };", fmt.Sprintf("%q", k), v.URI, v.Revision))
		}

		b.WriteRune('}')

		expr = b.String()
	}

	var nixpath NixPath
	{
		var stdout bytes.Buffer

		cmd := exec.Command("nix-instantiate", "--eval", "--expr", expr)
		cmd.Stdout = &stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			return nil, fmt.Errorf("error downloading path: %w", err)
		}

		s := ""
		err = json.Unmarshal(stdout.Bytes(), &s)
		if err != nil {
			return nil, fmt.Errorf("error decoding wrapping string for nixpath map: %w", err)
		}

		err = json.Unmarshal([]byte(s), &nixpath)
		if err != nil {
			return nil, fmt.Errorf("error decoding nixpath map: %w", err)
		}
	}

	return nixpath, nil
}

type HydraEvalResponse struct {
	Next  string       `json:"next"`
	First string       `json:"first"`
	Evals []*HydraEval `json:"evals"`

	baseURL string
	project string
	jobset  string
}

type HydraJobset struct {
	NixExprInput string `json:"nixexprinput"`
	NixExprPath  string `json:"nixexprpath"`
}

func (h *HydraEvalResponse) NextPage() (*HydraEvalResponse, error) {
	queryParams := strings.TrimPrefix(h.Next, "?")

	if queryParams == "" {
		return nil, io.EOF
	}

	return getEvaluations(h.baseURL, h.project, h.jobset, queryParams)
}

func getEvaluations(baseURL string, project string, jobset string, urlParamsQuery string) (*HydraEvalResponse, error) {
	extraQuery, err := url.ParseQuery(urlParamsQuery)
	if err != nil {
		return nil, fmt.Errorf("error parsing extra params: %w", err)
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing baseURL: %w", err)
	}

	u.Path = path.Join("jobset", project, jobset, "evals")
	u.RawQuery = extraQuery.Encode()

	client := &http.Client{}

	urlString := u.String()

	l := log.WithFields(log.Fields{
		"url": urlString,
	})

	l.Info("Requesting hydra evaluations")

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("accept", "application/json")

	start := time.Now().UTC()

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	l.WithFields(log.Fields{
		"duration": time.Since(start),
	}).Info("finished request")

	var r *HydraEvalResponse

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}

	r.baseURL = baseURL
	r.project = project
	r.jobset = jobset

	return r, nil
}

// Get evaluation metadata from hydra jobset
func GetEvaluations(baseURL string, project string, jobset string) (*HydraEvalResponse, error) {
	return getEvaluations(baseURL, project, jobset, "")
}

// Get a hydra project
func GetJobset(baseURL string, project string, jobset string) (*HydraJobset, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing baseURL: %w", err)
	}

	u.Path = path.Join("jobset", project, jobset)

	l := log.WithFields(log.Fields{
		"url": u.String(),
	})

	l.Info("Requesting hydra jobset")

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	var r *HydraJobset

	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}

	return r, nil
}
