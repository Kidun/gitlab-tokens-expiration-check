package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xanzy/go-gitlab"
)

func gitlabConnect(pat string, url string) (*gitlab.Client, error) {
	// Initialize the GitLab API client with your access token
	client, err := gitlab.NewClient(pat, gitlab.WithBaseURL(url))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	return client, err
}

func main() {
	const defToken = ""
	const defURL = "https://gitlab.com"

	listOptions := gitlab.ListOptions{
		PerPage: 100,
		Page:    1,
	}
	gitlabToken := defToken
	url := defURL
	cnt := 0

	if len(os.Args) > 2 {
		gitlabToken = os.Args[1]
		url = os.Args[2]
	} else {
		fmt.Printf("Usage %s token url\n", os.Args[0])
		return
	}

	gitlabClient, err := gitlabConnect(gitlabToken, url)
	if err != nil {
		log.Fatal(err)
	}

	// Projects' tokens
	for {
		projects, resp, err := gitlabClient.Projects.ListProjects(&gitlab.ListProjectsOptions{ListOptions: listOptions})
		if err != nil {
			log.Fatal(err)
		}

		for _, prj := range projects {
			fmt.Printf("\r %d", cnt) //progress counter
			cnt++

			tokens, _, err := gitlabClient.ProjectAccessTokens.ListProjectAccessTokens(prj.ID, &gitlab.ListProjectAccessTokensOptions{Page: 1, PerPage: 20})
			for _, token := range tokens {
				fmt.Printf("\rProject #%d (%s): Token #%d (%s) expires %s\n", prj.ID, prj.Name, token.ID, token.Name, token.ExpiresAt.String())
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		if resp.CurrentPage >= resp.TotalPages {
			//resp.TotalPages {
			break
		}
		listOptions.Page = resp.NextPage
	}

	// Groups' tokens
	listOptions.Page = 1
	for {
		groups, resp, err := gitlabClient.Groups.ListGroups(&gitlab.ListGroupsOptions{ListOptions: listOptions})
		if err != nil {
			log.Fatal(err)
		}

		for _, group := range groups {
			fmt.Printf("\r %d", cnt) //progress counter
			cnt++

			tokens, _, err := gitlabClient.GroupAccessTokens.ListGroupAccessTokens(group.ID, &gitlab.ListGroupAccessTokensOptions{Page: 1, PerPage: 20})
			for _, token := range tokens {
				fmt.Printf("\rGroup #%d (%s): Token #%d (%s) expires %s\n", group.ID, group.Name, token.ID, token.Name, token.ExpiresAt.String())
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		if resp.CurrentPage >= resp.TotalPages {
			//resp.TotalPages {
			break
		}
		listOptions.Page = resp.NextPage
	}
}
