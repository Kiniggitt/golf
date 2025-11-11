package membership

import (
	"fmt"
	"strings"
)

type User struct {
	Username string
	First    string
	Last     string
	Suffix   string
	Address  string
	balance  float64
	fees     float64
}

func (u *User) GetBalance() string {
	return fmt.Sprintf("$%.2f", u.balance)
}

func (u *User) GetFees() string {
	return fmt.Sprintf("$%.2f", u.fees)
}

func (u *User) GetName() string {
	return fmt.Sprintf("%s %s %s", trim(u.First), trim(u.Last), trim(u.Suffix))
}

func trim(s string) string {
	return strings.Trim(s, " ")
}

var idMap map[string]*User = make(map[string]*User)

func (user *User) setName(name string) {
	names := strings.Split(name, " ")
	for i, name := range names {
		switch i {
		case 0:
			user.First = name
		case 1:
			user.Last = name
		case 2:
			user.Suffix = name
		}
	}
}

func UpdateUser(username, name, address string) {
	user, ok := GetUser(username)
	if ok {
		user.setName(name)
		user.Address = address
	}
}

func (user *User) getNewUsername() string {
	username := combineStrings(strings.ToLower(user.First), strings.ToLower(user.Last))
	if len(username) > len(user.First) {
		for i := range user.First {
			if !UsernameExists(username) {
				return username
			}
			username = combineStrings(strings.ToLower(user.First[:i+1]), strings.ToLower(user.Last))
		}
	}
	for count := 1; UsernameExists(username); count++ {
		username = combineStrings(username, string(rune(count)))
	}
	return username
}

func (user User) String() string {
	return fmt.Sprintf("username: %s\nfirstname: %s\nlastname: %s\nbalance: $%.2f", user.Username, user.First, user.Last, user.balance)
}

func NewUser(name, address string) (user User) {
	user.setName(name)
	username := user.getNewUsername()
	user.Username = username
	idMap[username] = &user
	user.Address = address
	user.fees = 0
	user.balance = 0
	return user
}

func AddUsers(newUsers []*User) {
	for _, user := range newUsers {
		if !UsernameExists(user.Username) {
			idMap[user.Username] = user
			fmt.Println("Added User", user)
		}
	}
}

func combineStrings(str ...string) (newStr string) {
	if len(str) == 1 {
		return str[0]
	}
	for _, s := range str {
		newStr = fmt.Sprintf("%s%s", newStr, s)
	}
	return
}

func UsernameExists(str string) bool {
	_, ok := idMap[str]
	return ok
}

func GetUser(str string) (*User, bool) {
	fmt.Println("Get user", str)
	user, ok := idMap[str]
	return user, ok
}

func GetUsers() []*User {
	users := make([]*User, 0)
	for _, user := range idMap {
		users = append(users, user)
	}
	return users
}
