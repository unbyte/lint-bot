package server

import (
	"context"
	"fmt"
	"github.com/google/go-github/v35/github"
	"github.com/unbyte/lint-bot/config"
	"github.com/unbyte/lint-bot/formatter"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

type Server struct {
	secret       []byte
	githubClient *github.Client
	rules        map[string]*Rule
}

type Rule struct {
	Produce    string
	Formatters []formatter.Formatter
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, s.secret)
	if err != nil {
		return
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return
	}

	switch ev := event.(type) {
	case *github.IssueCommentEvent:
		if c, o, e := s.handleComment(ev); e == nil && c != nil {
			if o != nil {
				_ = updateComment(s.githubClient, ev.GetRepo(), o.GetID(), c)
			} else {
				_ = createComment(s.githubClient, ev.GetRepo(), ev.GetIssue().GetNumber(), c)
			}
		}
	case *github.IssuesEvent:
		if c, o, e := s.handleIssue(ev); e == nil && c != nil {
			if o != nil {
				_ = updateComment(s.githubClient, ev.GetRepo(), o.GetID(), c)
			} else {
				_ = createComment(s.githubClient, ev.GetRepo(), ev.GetIssue().GetNumber(), c)
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleIssue(event *github.IssuesEvent) (comment *github.IssueComment, oldOne *github.IssueComment, err error) {
	blocks := applyRules(getMarkedBlocks(event.GetIssue().GetBody()), s.rules)
	if len(blocks) == 0 {
		return nil, nil, nil
	}
	body := hiddenTagForIssue() + concatBlocks(blocks)
	comment = &github.IssueComment{
		Body: &body,
	}
	switch event.GetAction() {
	case "opened":
		return
	case "edited":
		oldOne, err = findIssueResponse(s.githubClient, event)
		return
	}
	return
}

func (s *Server) handleComment(event *github.IssueCommentEvent) (comment *github.IssueComment, oldOne *github.IssueComment, err error) {
	if event.GetAction() == "deleted" {
		oldOne, err = findCommentResponse(s.githubClient, event)
		if oldOne == nil {
			return
		}
		err = deleteComment(s.githubClient, event.GetRepo(), oldOne.GetID())
		return
	}
	blocks := applyRules(getMarkedBlocks(event.GetComment().GetBody()), s.rules)
	if len(blocks) == 0 {
		return nil, nil, nil
	}
	body := hiddenTagForComment(event.GetComment().GetID()) + concatBlocks(blocks)
	comment = &github.IssueComment{
		Body: &body,
	}
	switch *event.Action {
	case "created":
		return
	case "edited":
		oldOne, err = findCommentResponse(s.githubClient, event)
		return
	}
	return
}

var _ http.Handler = &Server{}

func NewServer(config *config.Config) (http.Handler, error) {
	ctx := context.TODO()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Auth.PAT},
	)

	rules, err := resolveRules(config.Rules)
	if err != nil {
		return nil, err
	}

	return &Server{
		secret:       []byte(config.Auth.Secret),
		githubClient: github.NewClient(oauth2.NewClient(ctx, ts)),
		rules:        rules,
	}, nil
}

func resolveRules(raw []config.Rule) (map[string]*Rule, error) {
	result := make(map[string]*Rule)
	for _, r := range raw {
		fs := make([]formatter.Formatter, 0, 4)
		for _, name := range r.Formatters {
			if l, ok := formatter.Formatters[strings.ToLower(name)]; ok {
				fs = append(fs, l)
			} else {
				return nil, fmt.Errorf(`no formatter of name "%s"`, name)
			}
		}
		result[r.Consume] = &Rule{
			Produce:    r.Produce,
			Formatters: fs,
		}
	}
	return result, nil
}
