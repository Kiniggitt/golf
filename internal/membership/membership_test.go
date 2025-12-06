package membership

import (
	"strings"
	"testing"
	"time"
)

// Helper function to reset the idMap between tests
func resetUsers() {
	idMap = make(map[string]*User)
}

func TestNewFee(t *testing.T) {
	fee := NewFee(25.50, "Cart rental")

	if fee.Amount != 25.50 {
		t.Errorf("Expected amount 25.50, got %.2f", fee.Amount)
	}
	if fee.Reason != "Cart rental" {
		t.Errorf("Expected reason 'Cart rental', got %s", fee.Reason)
	}
	if fee.Paid {
		t.Error("Expected fee.Paid to be false for new fee")
	}
}

func TestUserGetBalance(t *testing.T) {
	resetUsers()

	user := User{
		Username: "testuser",
		Balance:  10.00,
		Fees: []Fee{
			{Amount: 5.00, Reason: "Fee 1", Paid: false},
			{Amount: 3.00, Reason: "Fee 2", Paid: false},
			{Amount: 2.00, Reason: "Fee 3", Paid: true}, // This should not be counted
		},
	}

	balance := user.GetBalance()
	expected := "$18.00" // 10.00 + 5.00 + 3.00

	if balance != expected {
		t.Errorf("Expected balance %s, got %s", expected, balance)
	}
}

func TestUserGetBalanceNoFees(t *testing.T) {
	resetUsers()

	user := User{
		Username: "testuser",
		Balance:  15.00,
		Fees:     []Fee{},
	}

	balance := user.GetBalance()
	expected := "$15.00"

	if balance != expected {
		t.Errorf("Expected balance %s, got %s", expected, balance)
	}
}

func TestUserGetFees(t *testing.T) {
	resetUsers()

	user := User{
		Fees: []Fee{
			{Amount: 10.00, Reason: "New account fee", Paid: false},
			{Amount: 5.00, Reason: "Late payment", Paid: false},
		},
	}

	fees := user.GetFees()
	expected := "New account fee - $10.00\nLate payment - $5.00"

	if fees != expected {
		t.Errorf("Expected fees:\n%s\nGot:\n%s", expected, fees)
	}
}

func TestUserGetFeesEmpty(t *testing.T) {
	resetUsers()

	user := User{
		Fees: []Fee{},
	}

	fees := user.GetFees()
	expected := ""

	if fees != expected {
		t.Errorf("Expected empty string, got %s", fees)
	}
}

func TestUserGetName(t *testing.T) {
	resetUsers()

	tests := []struct {
		name     string
		user     User
		expected string
	}{
		{
			name:     "Full name with suffix",
			user:     User{First: "John", Last: "Doe", Suffix: "Jr."},
			expected: "John Doe Jr.",
		},
		{
			name:     "Name without suffix",
			user:     User{First: "Jane", Last: "Smith", Suffix: ""},
			expected: "Jane Smith",
		},
		{
			name:     "Name with extra spaces",
			user:     User{First: " Bob ", Last: " Johnson ", Suffix: " Sr. "},
			expected: "Bob Johnson Sr.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.GetName()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestUserAddFee(t *testing.T) {
	resetUsers()

	user := User{
		Fees: []Fee{},
	}

	user.AddFee(15.00, "Test fee")

	if len(user.Fees) != 1 {
		t.Fatalf("Expected 1 fee, got %d", len(user.Fees))
	}

	if user.Fees[0].Amount != 15.00 {
		t.Errorf("Expected fee amount 15.00, got %.2f", user.Fees[0].Amount)
	}

	if user.Fees[0].Reason != "Test fee" {
		t.Errorf("Expected fee reason 'Test fee', got %s", user.Fees[0].Reason)
	}

	if user.Fees[0].Paid {
		t.Error("Expected new fee to be unpaid")
	}
}

func TestNewUser(t *testing.T) {
	resetUsers()

	user, err := NewUser("John Doe", "123 Main St")
	if err != nil {
		t.Fatalf("NewUser failed: %v", err)
	}

	if user.Username == "" {
		t.Error("Expected username to be set")
	}

	if user.First != "John" {
		t.Errorf("Expected first name 'John', got %s", user.First)
	}

	if user.Last != "Doe" {
		t.Errorf("Expected last name 'Doe', got %s", user.Last)
	}

	if user.Address != "123 Main St" {
		t.Errorf("Expected address '123 Main St', got %s", user.Address)
	}

	if user.Balance != 0 {
		t.Errorf("Expected balance 0, got %.2f", user.Balance)
	}

	if len(user.Fees) != 1 {
		t.Fatalf("Expected 1 fee (new account fee), got %d", len(user.Fees))
	}

	if user.Fees[0].Amount != 10.00 {
		t.Errorf("Expected new account fee of 10.00, got %.2f", user.Fees[0].Amount)
	}

	if user.Fees[0].Reason != "New account fee" {
		t.Errorf("Expected fee reason 'New account fee', got %s", user.Fees[0].Reason)
	}
}

func TestNewUserUsername(t *testing.T) {
	resetUsers()

	// Test username generation
	user1, err := NewUser("John Doe", "123 Main St")
	if err != nil {
		t.Fatalf("NewUser failed: %v", err)
	}
	if user1.Username != "johndoe" {
		t.Errorf("Expected username 'johndoe', got %s", user1.Username)
	}

	// Test duplicate username handling
	user2, err := NewUser("John Doe", "456 Oak Ave")
	if err != nil {
		t.Fatalf("NewUser failed: %v", err)
	}
	if user2.Username == "johndoe" {
		t.Error("Expected different username for duplicate name")
	}
}

func TestUsernameExists(t *testing.T) {
	resetUsers()

	if UsernameExists("testuser") {
		t.Error("Expected UsernameExists to return false for non-existent user")
	}

	user, err := NewUser("Test User", "123 Main St")
	if err != nil {
		t.Fatalf("NewUser failed: %v", err)
	}

	if !UsernameExists(user.Username) {
		t.Errorf("Expected UsernameExists to return true for user %s", user.Username)
	}
}

func TestGetUser(t *testing.T) {
	resetUsers()

	user, err := NewUser("Test User", "123 Main St")
	if err != nil {
		t.Fatalf("NewUser failed: %v", err)
	}

	retrieved, ok := GetUser(user.Username)
	if !ok {
		t.Errorf("Expected to find user %s", user.Username)
	}

	if retrieved.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrieved.Username)
	}

	_, ok = GetUser("nonexistent")
	if ok {
		t.Error("Expected GetUser to return false for non-existent user")
	}
}

func TestGetUsers(t *testing.T) {
	resetUsers()

	// Test with no users
	users := GetUsers()
	if len(users) != 0 {
		t.Errorf("Expected 0 users, got %d", len(users))
	}

	// Add some users
	_, err := NewUser("Alice Smith", "111 First St")
	if err != nil {
		t.Fatalf("NewUser failed: %v", err)
	}
	_, err = NewUser("Bob Jones", "222 Second St")
	if err != nil {
		t.Fatalf("NewUser failed: %v", err)
	}
	_, err = NewUser("Charlie Brown", "333 Third St")
	if err != nil {
		t.Fatalf("NewUser failed: %v", err)
	}

	users = GetUsers()
	if len(users) != 3 {
		t.Fatalf("Expected 3 users, got %d", len(users))
	}

	// Verify they are sorted by username
	for i := 0; i < len(users)-1; i++ {
		if users[i].Username >= users[i+1].Username {
			t.Errorf("Users not sorted correctly: %s should come before %s",
				users[i].Username, users[i+1].Username)
		}
	}
}

func TestAddUsers(t *testing.T) {
	resetUsers()

	newUsers := []*User{
		{Username: "alice", First: "Alice", Last: "Smith", Address: "111 First St"},
		{Username: "bob", First: "Bob", Last: "Jones", Address: "222 Second St"},
	}

	AddUsers(newUsers)

	if !UsernameExists("alice") {
		t.Error("Expected user 'alice' to exist")
	}

	if !UsernameExists("bob") {
		t.Error("Expected user 'bob' to exist")
	}

	// Test adding duplicate user (should be skipped)
	duplicateUsers := []*User{
		{Username: "alice", First: "Alice", Last: "Duplicate", Address: "999 Dup St"},
	}

	AddUsers(duplicateUsers)

	user, _ := GetUser("alice")
	if user.Last == "Duplicate" {
		t.Error("Expected duplicate user to be rejected")
	}
}

func TestUpdateUser(t *testing.T) {
	resetUsers()

	user, err := NewUser("John Doe", "123 Main St")
	if err != nil {
		t.Fatalf("NewUser failed: %v", err)
	}

	err = UpdateUser(user.Username, "Jane Smith Jr.", "456 Oak Ave")
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	updated, ok := GetUser(user.Username)
	if !ok {
		t.Fatal("Expected to find updated user")
	}

	if updated.First != "Jane" {
		t.Errorf("Expected first name 'Jane', got %s", updated.First)
	}

	if updated.Last != "Smith" {
		t.Errorf("Expected last name 'Smith', got %s", updated.Last)
	}

	if updated.Suffix != "Jr." {
		t.Errorf("Expected suffix 'Jr.', got %s", updated.Suffix)
	}

	if updated.Address != "456 Oak Ave" {
		t.Errorf("Expected address '456 Oak Ave', got %s", updated.Address)
	}
}

func TestUpdateUserNonExistent(t *testing.T) {
	resetUsers()

	// This should return an error
	err := UpdateUser("nonexistent", "Test Name", "Test Address")
	if err == nil {
		t.Error("Expected error when updating non-existent user")
	}

	if UsernameExists("nonexistent") {
		t.Error("Expected non-existent user to remain non-existent")
	}
}

func TestUserString(t *testing.T) {
	resetUsers()

	user := User{
		Username: "johndoe",
		First:    "John",
		Last:     "Doe",
		Balance:  25.50,
	}

	str := user.String()

	if !strings.Contains(str, "johndoe") {
		t.Error("Expected string representation to contain username")
	}

	if !strings.Contains(str, "John") {
		t.Error("Expected string representation to contain first name")
	}

	if !strings.Contains(str, "Doe") {
		t.Error("Expected string representation to contain last name")
	}

	if !strings.Contains(str, "$25.50") {
		t.Error("Expected string representation to contain balance")
	}
}

func TestNewUserEmptyName(t *testing.T) {
	resetUsers()

	// Test with empty name
	_, err := NewUser("", "123 Main St")
	if err == nil {
		t.Error("Expected error for empty name")
	}

	// Test with whitespace-only name
	_, err = NewUser("   ", "123 Main St")
	if err == nil {
		t.Error("Expected error for whitespace-only name")
	}
}

func TestUpdateUserEmptyName(t *testing.T) {
	resetUsers()

	user, err := NewUser("John Doe", "123 Main St")
	if err != nil {
		t.Fatalf("NewUser failed: %v", err)
	}

	// Test with empty name
	err = UpdateUser(user.Username, "", "456 Oak Ave")
	if err == nil {
		t.Error("Expected error for empty name")
	}

	// Test with whitespace-only name
	err = UpdateUser(user.Username, "   ", "456 Oak Ave")
	if err == nil {
		t.Error("Expected error for whitespace-only name")
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{" test ", "test"},
		{"test", "test"},
		{"  test  ", "test"},
		{"", ""},
	}

	for _, tt := range tests {
		result := strings.TrimSpace(tt.input)
		if result != tt.expected {
			t.Errorf("trim(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestCombineStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "Single string",
			input:    []string{"test"},
			expected: "test",
		},
		{
			name:     "Two strings",
			input:    []string{"hello", "world"},
			expected: "helloworld",
		},
		{
			name:     "Multiple strings",
			input:    []string{"a", "b", "c", "d"},
			expected: "abcd",
		},
		{
			name:     "Empty strings",
			input:    []string{"", "", ""},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := combineStrings(tt.input...)
			if result != tt.expected {
				t.Errorf("combineStrings(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetDaysSinceLastFee(t *testing.T) {
	resetUsers()

	t.Run("User with no fees", func(t *testing.T) {
		user := &User{
			Username: "testuser",
			Fees:     []Fee{},
		}
		result := user.GetDaysSinceLastFee()
		if result != "N/A" {
			t.Errorf("Expected 'N/A' for user with no fees, got %s", result)
		}
	})

	t.Run("User with new fee", func(t *testing.T) {
		user := &User{
			Username: "testuser",
			Fees:     []Fee{NewFee(10.00, "Test fee")},
		}
		result := user.GetDaysSinceLastFee()
		if result != "0" {
			t.Errorf("Expected '0' for fee created today, got %s", result)
		}
	})

	t.Run("User with legacy fee (no timestamp)", func(t *testing.T) {
		user := &User{
			Username: "testuser",
			Fees: []Fee{
				{Amount: 10.00, Reason: "Legacy fee", Paid: false}, // CreatedAt will be zero time
			},
		}
		result := user.GetDaysSinceLastFee()
		if result != "Unknown" {
			t.Errorf("Expected 'Unknown' for legacy fee without timestamp, got %s", result)
		}
	})

	t.Run("User with old fee", func(t *testing.T) {
		oldFee := NewFee(10.00, "Old fee")
		oldFee.CreatedAt = time.Now().AddDate(0, 0, -30) // 30 days ago
		user := &User{
			Username: "testuser",
			Fees:     []Fee{oldFee},
		}
		result := user.GetDaysSinceLastFee()
		if result != "30" {
			t.Errorf("Expected '30' for fee created 30 days ago, got %s", result)
		}
	})

	t.Run("User with multiple fees", func(t *testing.T) {
		oldFee := NewFee(10.00, "Old fee")
		oldFee.CreatedAt = time.Now().AddDate(0, 0, -30) // 30 days ago
		recentFee := NewFee(20.00, "Recent fee")
		recentFee.CreatedAt = time.Now().AddDate(0, 0, -5) // 5 days ago
		user := &User{
			Username: "testuser",
			Fees:     []Fee{oldFee, recentFee},
		}
		result := user.GetDaysSinceLastFee()
		if result != "5" {
			t.Errorf("Expected '5' for most recent fee, got %s", result)
		}
	})

	t.Run("User with mix of legacy and new fees", func(t *testing.T) {
		legacyFee := Fee{Amount: 10.00, Reason: "Legacy fee", Paid: false} // No timestamp
		newFee := NewFee(20.00, "New fee")
		newFee.CreatedAt = time.Now().AddDate(0, 0, -7) // 7 days ago
		user := &User{
			Username: "testuser",
			Fees:     []Fee{legacyFee, newFee},
		}
		result := user.GetDaysSinceLastFee()
		if result != "7" {
			t.Errorf("Expected '7' for most recent fee with timestamp, got %s", result)
		}
	})
}
