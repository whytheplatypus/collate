// # RCA: A tool to help us deal with rca's in github
//
// There are two problems to be solved.
// It's difficult and annoying to move an issue into
// a PR.
// It's hard to search the resulting files in the repo
//
// To help with these `rca` will attempt to transform
// an issue into a md doc that can be submitted in a PR
//
// It will also attempt to make that markdown document
// a bit easier to deal with by adding front matter
// containing the labels that had been applied to the issue
package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

var ghToken string
var templateFile string

const defaultTemplate = `
{{.Issue.GetBody}}
{{range .Comments}}
{{.GetBody}}
{{end}}
{{range .Labels}}
{{.GetName}}
{{end}}
`

func init() {
	const (
		tokenUsage = `
Use -token to supply a github token
This will try and use ~/.ghtoken as a default
`
		templateUseage = `
Use -tempalte to supply a tempalte file.
This template will have access to
.Issue
.Comments
.Labels
`
	)

	// Get a github token from a .ghtoken file
	defaultToken, err := ioutil.ReadFile(fmt.Sprintf("%s%s", os.Getenv("HOME"), "/.ghtoken"))
	if err != nil {
		panic(err)
	}

	flag.StringVar(&ghToken, "token", strings.TrimSpace(string(defaultToken)), tokenUsage)

	flag.StringVar(&templateFile, "template", "", templateUseage)
}

func main() {
	flag.Parse()
	org := flag.Arg(0)
	repo := flag.Arg(1)
	iNum, err := strconv.Atoi(flag.Arg(2))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	issue, _, err := client.Issues.Get(ctx, org, repo, iNum)
	if err != nil {
		panic(err)
	}

	comments, _, err := client.Issues.ListComments(ctx, org, repo, iNum, nil)
	if err != nil {
		panic(err)
	}

	labels, _, err := client.Issues.ListLabelsByIssue(ctx, org, repo, iNum, nil)
	if err != nil {
		panic(err)
	}

	docTmp := defaultTemplate
	if templateFile != "" {
		tb, err := ioutil.ReadFile(templateFile)
		if err != nil {
			panic(err)
		}
		docTmp = string(tb)
	}

	t := template.Must(template.New("rca-pr").Parse(docTmp))
	data := struct {
		Issue    *github.Issue
		Labels   []*github.Label
		Comments []*github.IssueComment
	}{
		issue,
		labels,
		comments,
	}
	if err := t.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
}
