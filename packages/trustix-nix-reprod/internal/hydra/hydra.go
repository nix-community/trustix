package hydra

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"strings"
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

type HydraEvalResponse struct {
	Next  string       `json:"next"`
	First string       `json:"first"`
	Evals []*HydraEval `json:"evals"`

	baseURL string
	project string
	jobset  string
}

// Get NIX_PATH from a Hydra evaluation
func (e *HydraEval) NixPath() (string, error) {

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
				return "", fmt.Errorf("unsupported hydra input type: %s", v.Type)
			}

			// Sanity check URI to prevent command injections
			_, err := url.Parse(v.URI)
			if err != nil {
				return "", fmt.Errorf("error parsing URL: %w", err)
			}

			// Sanity check rev to prevent command injections
			_, err = hex.DecodeString(v.Revision)
			if err != nil {
				return "", fmt.Errorf("error parsing revision: %w", err)
			}

			b.WriteString(fmt.Sprintf("  %s = builtins.fetchGit { url = \"%s\"; rev = \"%s\"; allRefs = true; };", fmt.Sprintf("%q", k), v.URI, v.Revision))
		}

		b.WriteRune('}')

		expr = b.String()
	}

	var nixpath map[string]string
	{
		var stdout bytes.Buffer

		cmd := exec.Command("nix-instantiate", "--eval", "--expr", expr)
		cmd.Stdout = &stdout

		err := cmd.Run()
		if err != nil {
			return "", fmt.Errorf("error downloading path: %w", err)
		}

		s := ""
		err = json.Unmarshal(stdout.Bytes(), &s)
		if err != nil {
			return "", fmt.Errorf("error decoding wrapping string for nixpath map: %w", err)
		}

		err = json.Unmarshal([]byte(s), &nixpath)
		if err != nil {
			return "", fmt.Errorf("error decoding nixpath map: %w", err)
		}
	}

	var b strings.Builder

	for k, v := range nixpath {
		b.WriteString(k)
		b.WriteRune('=')
		b.WriteString(v)
	}

	return b.String(), nil
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

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

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
