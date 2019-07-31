package main

// Hit the Github API to list private repositories
// within a private organization

import (
	"fmt"
	"context"
	"os"
	"flag"
	"strings"
	"io/ioutil"

	"github.com/mitchellh/go-homedir"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Required override to print our help message instead of
// something that isn't useful.
func customUsage() {
	fmt.Printf("List repositories that exist within a private organization on Github.\n\n")
	fmt.Printf("Requires a personal access token be placed in $HOME/.gittoken\n")
	fmt.Printf("\thttps://github.com/settings/tokens\n\n")

	fmt.Println("Required token scope:")
	fmt.Printf("\trepo:status, repo_deployment, public_repo, repo:invite, read:org\n\n")
	fmt.Printf("Usage: %s [OPTIONS] orgname\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	version := "v1.0.1"

	// Parse commandline argument(s)
	limitsPtr := flag.Bool("L", false, "Print your Github API limits.")
	versionPtr := flag.Bool("v", false, "Display the version of this utility.")
	/* revist this feature
	scopesPtr := flag.Bool("s", false, "Print your Github token scope.")
	*/
	flag.Usage = customUsage
	flag.Parse()



	// Get the user home directory and expand it if the parse returns
	// the '~' character.
	homeDir, _ := homedir.Dir()
	homeDirPath, _ := homedir.Expand(homeDir)

	// this comes in as a list of bytes and followed by a return/newline
	// character that needs to be stripped.
	data, _ := ioutil.ReadFile(homeDirPath + "/.gittoken")
	userToken := strings.TrimSpace(string(data))
	if len(userToken) == 0 {
		// this implies an error so we need to bail.
		fmt.Printf("Token file not found at %v.\n", homeDirPath + "/.gittoken")
		os.Exit(3)
	}

	// Establish authentication context with the token.
	context := context.Background()
	tokenService := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: userToken},
	)
	tokenClient := oauth2.NewClient(context, tokenService)

	g := github.NewClient(tokenClient)

	// If the user requests the version, do that and quit
	if *versionPtr == true {
		fmt.Printf("%s %s\n", os.Args[0], version)
		os.Exit(0)
	}

	// Check rate limits for authenticated user. If this fails,
	// we assume that your authentication has failed.
	rateLimit, _, err := g.RateLimits(context)
	if err != nil {
		fmt.Printf("Problem in getting rate limit information %v\n", err)
		os.Exit(4)
	} else if *limitsPtr == true {
		fmt.Print("Github API Limits:\n")
		fmt.Printf("Limit: %d \nRemaining %d \n",
					rateLimit.Core.Limit,
					rateLimit.Core.Remaining )
		os.Exit(0)
	}

	// we only take one positional argument, otherwise fail.
	if flag.NArg() > 1 || flag.NArg() == 0 {
		fmt.Printf("ERR This command takes exactly one argument\n")
		fmt.Println()
		flag.Usage()
		os.Exit(2)
	}

	// the Github org we'll be searching through.
	organization := flag.Args()[0]

	/* TODO Revisit this feature
	if *scopesPtr == true {
		scopes, _, _ := g.Authorizations.List(context, nil)
		if err != nil {
			fmt.Printf("Problem in getting token scopes %v.\n", err)
			os.Exit(5)
		} else {
			fmt.Println(scopes)
			os.Exit(0)
		}
	}
	*/ 

	// This needs to be set, but doesn't seem to handle paging properly.
	// TODO: figure out what the hell is happening here.
	opt := &github.RepositoryListByOrgOptions{}

	// Estabilsh the allRepos variable as a list of github.Repository objects.
	var allRepos []*github.Repository
	for {
		// iterate until we have all pages.
		repos, resp, err := g.Repositories.ListByOrg(context, organization, opt)
		if err != nil {
			fmt.Println(err)
			os.Exit(6)
		}

		// TODO: figure out this ellipsis
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// The first item returned is an index which we really don't need here.
	// for this reason we _
	divider := strings.Repeat("-", 35)
	for _, element := range allRepos {
		fmt.Println(divider)
		fmt.Printf("%14v %v\n", "Name:", element.GetName())
		fmt.Printf("%14v %v\n", "Link:", element.GetHTMLURL())
		fmt.Printf("%14v %v\n", "Clone (SSH):", element.GetSSHURL())
		fmt.Printf("%14v %v\n", "Clone (HTTPS):", element.GetCloneURL())
	}
	fmt.Println(divider)
}