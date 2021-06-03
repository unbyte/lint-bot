package server

import (
	"context"
	"fmt"
	"github.com/google/go-github/v35/github"
	"regexp"
	"strings"
)

func findCommentResponse(client *github.Client, event *github.IssueCommentEvent) (*github.IssueComment, error) {
	list, _, err := client.Issues.ListComments(context.TODO(),
		event.GetRepo().GetOwner().GetLogin(),
		event.GetRepo().GetName(),
		event.GetIssue().GetNumber(),
		&github.IssueListCommentsOptions{
			Since: event.GetComment().CreatedAt,
		})
	if err != nil {
		return nil, err
	}
	for _, comment := range list {
		if strings.HasPrefix(comment.GetBody(), hiddenTagForComment(event.GetComment().GetID())) {
			return comment, nil
		}
	}
	return nil, nil
}

func findIssueResponse(client *github.Client, event *github.IssuesEvent) (*github.IssueComment, error) {
	list, _, err := client.Issues.ListComments(context.TODO(),
		event.GetRepo().GetOwner().GetLogin(),
		event.GetRepo().GetName(),
		event.GetIssue().GetNumber(),
		&github.IssueListCommentsOptions{},
	)
	if err != nil {
		return nil, err
	}
	for _, comment := range list {
		if strings.HasPrefix(comment.GetBody(), hiddenTagForIssue()) {
			return comment, nil
		}
	}
	return nil, nil
}

func hiddenTagForIssue() string {
	return "<!--formatted-for-issue-->"
}

func hiddenTagForComment(id int64) string {
	return fmt.Sprintf("<!--formatted-for-comment#%d-->", id)
}

func createComment(client *github.Client, repo *github.Repository, issueID int, comment *github.IssueComment) error {
	_, _, err := client.Issues.CreateComment(context.TODO(), repo.GetOwner().GetLogin(), repo.GetName(), int(issueID), comment)
	return err
}

func updateComment(client *github.Client, repo *github.Repository, commentID int64, comment *github.IssueComment) error {
	_, _, err := client.Issues.EditComment(context.TODO(), repo.GetOwner().GetLogin(), repo.GetName(), commentID, comment)
	return err
}

func deleteComment(client *github.Client, repo *github.Repository, commentID int64) error {
	_, err := client.Issues.DeleteComment(context.TODO(), repo.GetOwner().GetLogin(), repo.GetName(), commentID)
	return err
}

type Block struct {
	Mark string
	Body string
}

var blocksRegexp = regexp.MustCompile("```(?P<mark>\\S+)\\r?\\n(?P<body>[\\s\\S]*?)```")

func getMarkedBlocks(body string) []*Block {
	matched := blocksRegexp.FindAllStringSubmatch(body, -1)
	records := make([]*Block, 0, len(matched))
	for _, m := range matched {
		record := make(map[string]string)
		for i, name := range blocksRegexp.SubexpNames() {
			if i != 0 && name != "" {
				record[name] = m[i]
			}
		}
		records = append(records, &Block{
			Mark: record["mark"],
			Body: record["body"],
		})
	}
	return records
}

func applyRules(blocks []*Block, rules map[string]*Rule) []*Block {
	result := make([]*Block, 0, len(blocks))
	for _, block := range blocks {
		r, ok := rules[block.Mark]
		if !ok {
			continue
		}
		body := block.Body
		for _, f := range r.Formatters {
			body = f.Handle(body)
		}
		if body == block.Body {
			continue
		}
		result = append(result, &Block{
			Mark: r.Produce,
			Body: body,
		})
	}
	return result
}

func concatBlocks(blocks []*Block) string {
	var sb strings.Builder
	for _, block := range blocks {
		sb.WriteString(fmt.Sprintf(
			"<details>\n<summary>Click to see formatted blocks</summary>\n\n```%s\n%s```\n\n</details>",
			block.Mark,
			block.Body,
		))
	}
	return sb.String()
}
