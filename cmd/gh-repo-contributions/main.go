package main

import (
	"fmt"
	"os"

	"github.com/QuentinPerez/gh-repo-contributions"
	"github.com/Sirupsen/logrus"
)

func main() {
	if len(os.Args) != 2 {
		logrus.Fatalf("Usage: %v <GithubName>", os.Args[0])
	}
	result, err := github.GetAllContibutions(os.Args[1])
	if err != nil {
		logrus.Fatal(err)
	}
	if err := result.GetLanguages(); err != nil {
		logrus.Fatal(err)
	}

	maxlenRepo := 0
	maxlenLanguage := 10
	for _, c := range result {
		if len(c.Repo)+18 > maxlenRepo {
			maxlenRepo = len(c.Repo) + 18
		}
		if len(c.Languages[0].Name) > maxlenLanguage {
			maxlenLanguage = len(c.Languages[0].Name)
		}
	}
	fmt.Printf("%-*s | %-*s | %s\n", maxlenRepo, "Repositories", maxlenLanguage, "Languages", "Commits")
	for i := 0; i <= maxlenRepo; i++ {
		fmt.Printf("-")
	}
	fmt.Printf("|")
	for i := 0; i <= maxlenLanguage; i++ {
		fmt.Printf("-")
	}
	fmt.Printf("-|")
	for i := 0; i < 10; i++ {
		fmt.Printf("-")
	}
	fmt.Println()
	total := 0
	for _, c := range result {
		// FIXME Don't use + operator to concat strings
		fmt.Printf("%-*s | %-*s | %d\n", maxlenRepo, "http://github.com/"+c.Repo, maxlenLanguage, c.Languages[0].Name, c.Commits)
		total += c.Commits
	}
	fmt.Printf("%-*s | %-*s | %d\n", maxlenRepo, "", maxlenLanguage, "", total)
	fmt.Println()
}
