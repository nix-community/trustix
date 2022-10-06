package hydra

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type JobsetEvalInput struct {
	Revision string `json:"revision"`
}

type HydraEval struct {
	ID         int64                       `json:"id"`
	Timestamp  int64                       `json:"timestamp"`
	EvalInputs map[string]*JobsetEvalInput `json:"evalinputs"`
}

type HydraEvalResponse struct {
	Next  string       `json:"next"`
	First string       `json:"first"`
	Evals []*HydraEval `json:"evals"`

	baseURL string
	project string
	jobset  string
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
