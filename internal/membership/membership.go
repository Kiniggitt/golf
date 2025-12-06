// Package membership provides golf course member management functionality.
// It handles user registration, username generation, and fee tracking.
package membership

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// User represents a golf course member with their account information.
// Each user has a unique username, personal information, account balance,
// and a list of assigned fees. The balance calculation automatically
// includes all unpaid fees.
type User struct {
	// Username is the unique identifier for this user, automatically
	// generated from their name during registration.
	Username string

	// First is the user's first name.
	First string

	// Last is the user's last name.
	Last string

	// Suffix is an optional name suffix (Jr., Sr., III, etc.).
	Suffix string

	// Address is the user's mailing address.
	Address string

	// Balance is the base account balance before fees are applied.
	Balance float64

	// Fees is the list of all fees assigned to this user.
	Fees []Fee
}

// GetBalance returns the total account balance including all unpaid fees,
// formatted as a currency string (e.g., "$25.00").
func (u *User) GetBalance() string {
	balance := u.Balance
	for _, fee := range u.Fees {
		if !fee.Paid {
			balance += fee.Amount
		}
	}
	return fmt.Sprintf("$%.2f", balance)
}

// GetFees returns a newline-separated string of all fees assigned to this user.
// Each line contains the fee reason and amount (e.g., "Cart Rental - $25.00").
func (u *User) GetFees() string {
	var fees []string
	for _, fee := range u.Fees {
		fees = append(fees, fmt.Sprintf("%s - $%.2f", fee.Reason, fee.Amount))
	}
	return strings.Join(fees, "\n")
}

// GetName returns the user's full name by combining First, Last, and Suffix.
// Extra whitespace is trimmed from each component and the final result.
func (u *User) GetName() string {
	name := fmt.Sprintf("%s %s %s", strings.TrimSpace(u.First), strings.TrimSpace(u.Last), strings.TrimSpace(u.Suffix))
	return strings.TrimSpace(name)
}

// AddFee adds a new unpaid fee to the user's account.
// The fee is automatically marked as unpaid and added to the Fees slice.
func (u *User) AddFee(amt float64, reason string) {
	u.Fees = append(u.Fees, NewFee(amt, reason))
}

// GetDaysSinceLastFee returns the number of days since the most recent fee was created.
// Returns "N/A" if the user has no fees.
// For legacy fees without timestamps (loaded from old data), returns "Unknown".
func (u *User) GetDaysSinceLastFee() string {
	if len(u.Fees) == 0 {
		return "N/A"
	}

	// Find the most recent fee with a valid timestamp
	var mostRecent time.Time
	hasValidTimestamp := false

	for _, fee := range u.Fees {
		if !fee.CreatedAt.IsZero() {
			if fee.CreatedAt.After(mostRecent) {
				mostRecent = fee.CreatedAt
				hasValidTimestamp = true
			}
		}
	}

	// If no valid timestamp found (legacy fees from old data), return "Unknown"
	if !hasValidTimestamp {
		return "Unknown"
	}

	days := int(time.Since(mostRecent).Hours() / 24)
	return fmt.Sprintf("%d", days)
}

var idMap = make(map[string]*User)

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

// UpdateUser updates an existing user's name and address.
// The name is parsed into First, Last, and Suffix components.
// Returns an error if the username doesn't exist or the name is empty.
func UpdateUser(username, name, address string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	user, ok := GetUser(username)
	if !ok {
		return fmt.Errorf("user %s not found", username)
	}
	user.setName(name)
	user.Address = address
	return nil
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

		username = combineStrings(username, strconv.Itoa(count))
	}
	return username
}

// String returns a string representation of the User for debugging purposes.
// It includes the username, first name, last name, and balance.
func (user User) String() string {
	return fmt.Sprintf("username: %s\nfirstname: %s\nlastname: %s\nbalance: $%.2f",
		user.Username, user.First, user.Last, user.Balance)
}

// NewUser creates a new golf course member with an automatically generated username.
// The username is derived from the user's name and guaranteed to be unique.
// A default "New account fee" of $10.00 is automatically added.
// Returns an error if the name is empty or username generation fails.
func NewUser(name, address string) (*User, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	user := &User{}
	user.setName(name)
	username := user.getNewUsername()
	if username == "" {
		return nil, fmt.Errorf("failed to generate username")
	}
	user.Username = username
	idMap[username] = user
	user.Address = address
	user.Fees = make([]Fee, 0)
	user.AddFee(10.00, "New account fee")
	user.Balance = 0
	return user, nil
}

// AddUsers adds multiple users to the system from a slice.
// This is typically used when loading users from persistent storage.
// Users with usernames that already exist are skipped.
func AddUsers(newUsers []*User) {
	for _, user := range newUsers {
		if !UsernameExists(user.Username) {
			idMap[user.Username] = user
			fmt.Println("Added User", user)
		}
	}
}

func combineStrings(str ...string) string {
	return strings.Join(str, "")
}

// UsernameExists checks if a username is already registered in the system.
// Returns true if the username exists, false otherwise.
func UsernameExists(str string) bool {
	_, ok := idMap[str]
	return ok
}

// GetUser retrieves a user by their username.
// Returns the User and true if found, or nil and false if not found.
func GetUser(str string) (*User, bool) {
	user, ok := idMap[str]
	return user, ok
}

// GetUsers returns all registered users sorted alphabetically by username.
// This provides a consistent ordering for display purposes.
func GetUsers() []*User {
	users := make([]*User, 0)
	for _, user := range idMap {
		users = append(users, user)
	}
	// Sort users by username for consistent ordering
	sort.Slice(users, func(i, j int) bool {
		return users[i].Username < users[j].Username
	})
	return users
}
