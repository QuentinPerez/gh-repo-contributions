package github

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Contribution struct {
	Repo    string
	Commits int
}

func getBody(url string) ([]string, error) {
	for {
		resp, err := http.Get(url)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		tab := strings.Split(string(b), "\n")
		nbOfLines := len(tab)
		wait := false
		for i := 0; i < nbOfLines; i++ {
			if strings.Contains(tab[i], "You have triggered an abuse detection mechanism.") {
				time.Sleep(1 * time.Minute)
				wait = true
				break
			}
		}
		if !wait {
			// Bypass GitHub abuse detection
			value := rand.Intn(2000-1000) + 1000
			time.Sleep(time.Duration(value) * time.Millisecond)
			return tab, nil
		}
	}
}

func GetAllContibutions(pseudo string) ([]Contribution, error) {
	result := make(map[string]int)

	url := fmt.Sprintf("https://github.com/%s", pseudo)
	tab, err := getBody(url)
	if err != nil {
		return nil, err
	}
	nbOfLines := len(tab)
	joinYear := 0
	for i := 0; i < nbOfLines; i++ {
		if strings.Contains(tab[i], "join-date") {
			joinYear, _ = strconv.Atoi(tab[i][len(tab[i])-16 : len(tab[i])-12])
		}
	}
	if joinYear == 0 {
		return nil, fmt.Errorf("Coudln't find join-date")
	}

	year := time.Now().Year() + 1
	for y := joinYear; y < year; y++ {
		for m := 1; m < 13; m++ {
			if m == 12 {
				url = fmt.Sprintf("https://github.com/%s?tab=contributions&from=%d-%02d-01&to=%d-%02d-31", pseudo, y, m, y, m)
			} else {
				url = fmt.Sprintf("https://github.com/%s?tab=contributions&from=%d-%02d-01&to=%d-%02d-01", pseudo, y, m, y, m+1)
			}
			tab, err := getBody(url)
			if err != nil {
				return nil, err
			}
			nbOfLines := len(tab)
			for i := 0; i < nbOfLines; i++ {
				if strings.Contains(tab[i], "octicon octicon-git-commit") {
					for ; i < nbOfLines; i++ {
						if strings.Contains(tab[i], "href=") {
							i++
							contrib := strings.Fields(tab[i][:len(tab[i])-4])
							nb, _ := strconv.Atoi(contrib[1])
							if _, ok := result[contrib[4]]; !ok {
								result[contrib[4]] = nb
							} else {
								result[contrib[4]] = result[contrib[4]] + nb
							}
						}
						if strings.Contains(tab[i], "</ul>") {
							break
						}
					}
					break
				}
			}
		}
	}
	var keys []int
	findKeyByValue := func(m map[string]int, value int) string {
		for k, v := range m {
			if v == value {
				return k
			}
		}
		return ""
	}

	for _, v := range result {
		keys = append(keys, v)
	}
	sort.Ints(keys)
	ret := make([]Contribution, len(keys))
	inv := len(ret) - 1
	for i := range ret {
		k := findKeyByValue(result, keys[inv])
		ret[i] = Contribution{
			Repo:    k,
			Commits: keys[inv],
		}
		delete(result, k)
		inv--
	}
	return ret, nil
}
