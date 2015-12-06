package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/QuentinPerez/gh-repo-contributions"
	"github.com/Sirupsen/logrus"
)

func main() {
	if len(os.Args) != 2 {
		logrus.Fatalf("Usage: %v <GithubName>", os.Args[0])
	}
	result, err := github.GetAllContibutions(os.Args[1], "golang/go")
	if err != nil {
		logrus.Fatal(err)
	}
	if err := result.GetFromRepositoryPage(github.GetLanguage, github.GetStar); err != nil {
		logrus.Fatal(err)
	}

	maxlenRepo := 0
	maxlenLanguage := 10
	maxlenStar := 5
	for _, c := range result {
		if len(c.Repo)+18 > maxlenRepo {
			maxlenRepo = len(c.Repo) + 18
		}
		if len(c.Languages[0].Name) > maxlenLanguage {
			maxlenLanguage = len(c.Languages[0].Name)
		}
		if len(strconv.Itoa(c.Star)) > maxlenStar {
			maxlenStar = len(strconv.Itoa(c.Star))
		}
	}
	fmt.Printf("%-*s | %-*s | %-*s | %s\n", maxlenRepo, "Repositories", maxlenStar, "Star", maxlenLanguage, "Languages", "Commits")
	for i := 0; i <= maxlenRepo; i++ {
		fmt.Printf("-")
	}
	fmt.Printf("|")
	for i := 0; i <= maxlenStar; i++ {
		fmt.Printf("-")
	}
	fmt.Printf("-|")
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
		fmt.Printf("%-*s | %-*d | %-*s | %d\n", maxlenRepo, "http://github.com/"+c.Repo, maxlenStar, c.Star, maxlenLanguage, c.Languages[0].Name, c.Commits)
		total += c.Commits
	}
	fmt.Printf("%-*s | %-*s | %-*s | %d\n", maxlenRepo, "", maxlenStar, "", maxlenLanguage, "", total)
	fmt.Println()
}
