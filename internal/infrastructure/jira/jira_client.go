package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hanif-afianto/jira-timesheet-generator/internal/domain/entity"
	"github.com/hanif-afianto/jira-timesheet-generator/internal/infrastructure/config"
)

type JiraClient struct {
	config *config.Config
	client *http.Client
}

func NewJiraClient(cfg *config.Config) *JiraClient {
	return &JiraClient{
		config: cfg,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

type jiraResponse struct {
	Total         int         `json:"total"`
	NextPageToken string      `json:"nextPageToken"`
	Issues        []jiraIssue `json:"issues"`
}

type jiraIssue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary string      `json:"summary"`
		Worklog jiraWorklog `json:"worklog"`
	} `json:"fields"`
}

type jiraWorklog struct {
	Total    int               `json:"total"`
	Worklogs []jiraWorklogItem `json:"worklogs"`
}

type jiraWorklogItem struct {
	Author struct {
		AccountID string `json:"accountId"`
	} `json:"author"`
	TimeSpent string `json:"timeSpent"`
	Comment   struct {
		Content []struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"content"`
	} `json:"comment"`
	Started string `json:"started"`
}

func (j *JiraClient) FetchIssues(ctx context.Context, jql string) ([]entity.Issue, error) {
	url := fmt.Sprintf("%s/rest/api/3/search/jql", j.config.JiraBaseURL)
	var allIssues []entity.Issue
	nextPageToken := ""

	for {
		body := map[string]interface{}{
			"fieldsByKeys": true,
			"fields":       []string{"summary", "worklog"},
			"jql":          jql,
			"maxResults":   50,
		}
		if nextPageToken != "" {
			body["nextPageToken"] = nextPageToken
		}

		payload, _ := json.Marshal(body)
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(j.config.JiraEmail, j.config.JiraAPIToken)

		res, err := j.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			data, _ := io.ReadAll(res.Body)
			return nil, fmt.Errorf("jira API error: %s", string(data))
		}

		var result jiraResponse
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			return nil, err
		}

		for _, issue := range result.Issues {
			allIssues = append(allIssues, entity.Issue{
				Key:     issue.Key,
				Summary: issue.Fields.Summary,
				Worklogs: j.mapWorklogs(issue.Fields.Worklog.Worklogs),
			})
		}

		if result.NextPageToken == "" {
			break
		}
		nextPageToken = result.NextPageToken
	}

	return allIssues, nil
}

func (j *JiraClient) FetchWorklogs(ctx context.Context, issueKey string) ([]entity.Worklog, error) {
	var allWorklogs []entity.Worklog
	startAt := 0
	maxResults := 50

	for {
		url := fmt.Sprintf("%s/rest/api/3/issue/%s/worklog?startAt=%d&maxResults=%d", 
			j.config.JiraBaseURL, issueKey, startAt, maxResults)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Accept", "application/json")
		req.SetBasicAuth(j.config.JiraEmail, j.config.JiraAPIToken)

		res, err := j.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			data, _ := io.ReadAll(res.Body)
			return nil, fmt.Errorf("jira API error for worklogs (%s): %s", issueKey, string(data))
		}

		var result struct {
			Total    int               `json:"total"`
			Worklogs []jiraWorklogItem `json:"worklogs"`
		}
		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			return nil, err
		}

		allWorklogs = append(allWorklogs, j.mapWorklogs(result.Worklogs)...)

		if len(allWorklogs) >= result.Total {
			break
		}
		startAt += len(result.Worklogs)
	}

	return allWorklogs, nil
}

func (j *JiraClient) mapWorklogs(items []jiraWorklogItem) []entity.Worklog {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	var worklogs []entity.Worklog
	for _, item := range items {
		worklogs = append(worklogs, entity.Worklog{
			AuthorAccountID: item.Author.AccountID,
			TimeSpent:       item.TimeSpent,
			Comment:         j.extractComment(item),
			Started:         j.parseDate(item.Started, loc),
		})
	}
	return worklogs
}

func (j *JiraClient) parseDate(s string, loc *time.Location) time.Time {
	layouts := []string{
		"2006-01-02T15:04:05.000-0700",
		"2006-01-02T15:04:05.000Z",
		time.RFC3339,
	}
	for _, l := range layouts {
		if t, err := time.ParseInLocation(l, s, loc); err == nil {
			return t.In(loc)
		}
	}
	return time.Time{}
}

func (j *JiraClient) extractComment(item jiraWorklogItem) string {
	if len(item.Comment.Content) == 0 {
		return ""
	}
	first := item.Comment.Content[0]
	if len(first.Content) == 0 {
		return ""
	}
	return first.Content[0].Text
}
