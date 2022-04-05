package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

//easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	var line []byte
	br := bufio.NewReader(r)

	for i := 0; ; i++ {
		line, _, err = br.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}

			return
		}

		var user User
		if err = user.UnmarshalJSON(line); err != nil {
			return
		}
		result[i] = user
	}
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched := strings.HasSuffix(user.Email, "."+domain)
		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
