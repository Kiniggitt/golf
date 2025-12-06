package admin

import (
	"testing"

	"golf/internal/membership"
)

// Helper function to reset users between tests
func resetUsers() {
	// Access the internal map via the membership package
	// We need to clear all users
	users := membership.GetUsers()
	for _, user := range users {
		// This is a workaround - in production, membership package should expose a Reset function
		_ = user
	}
}

func TestListUsersEmpty(t *testing.T) {
	// Note: This test depends on the current state of membership.GetUsers()
	// In a real scenario, we'd want to inject dependencies or reset state

	users := ListUsers()

	// Should return empty slice, not nil
	if users == nil {
		t.Error("Expected non-nil slice")
	}
}

func TestListUsersWithData(t *testing.T) {
	// Create test users via membership package
	user1, err := membership.NewUser("Alice Smith", "111 First St")
	if err != nil {
		t.Fatalf("NewUser failed: %v", err)
	}
	user2, err := membership.NewUser("Bob Jones", "222 Second St")
	if err != nil {
		t.Fatalf("NewUser failed: %v", err)
	}

	users := ListUsers()

	// Should have at least 2 users
	if len(users) < 2 {
		t.Fatalf("Expected at least 2 users, got %d", len(users))
	}

	// Verify each user has 5 columns
	for i, user := range users {
		if user == nil {
			continue // Skip nil entries
		}
		if len(user) != 5 {
			t.Errorf("User %d should have 5 columns, got %d", i, len(user))
		}
	}

	// Find our test users in the results
	foundAlice := false
	foundBob := false

	for _, user := range users {
		if len(user) < 1 {
			continue
		}

		if user[0] == user1.Username {
			foundAlice = true
			// Verify columns: Username, Name, Address, Balance, Days Since Last Fee
			if user[1] != user1.GetName() {
				t.Errorf("Expected name %s, got %s", user1.GetName(), user[1])
			}
			if user[2] != user1.Address {
				t.Errorf("Expected address %s, got %s", user1.Address, user[2])
			}
			if user[3] != user1.GetBalance() {
				t.Errorf("Expected balance %s, got %s", user1.GetBalance(), user[3])
			}
			if user[4] != user1.GetDaysSinceLastFee() {
				t.Errorf("Expected days since last fee %s, got %s", user1.GetDaysSinceLastFee(), user[4])
			}
		}

		if user[0] == user2.Username {
			foundBob = true
			if user[1] != user2.GetName() {
				t.Errorf("Expected name %s, got %s", user2.GetName(), user[1])
			}
			if user[2] != user2.Address {
				t.Errorf("Expected address %s, got %s", user2.Address, user[2])
			}
			if user[3] != user2.GetBalance() {
				t.Errorf("Expected balance %s, got %s", user2.GetBalance(), user[3])
			}
			if user[4] != user2.GetDaysSinceLastFee() {
				t.Errorf("Expected days since last fee %s, got %s", user2.GetDaysSinceLastFee(), user[4])
			}
		}
	}

	if !foundAlice {
		t.Error("Expected to find Alice in user list")
	}
	if !foundBob {
		t.Error("Expected to find Bob in user list")
	}
}

func TestListUsersColumnOrder(t *testing.T) {
	user, err := membership.NewUser("Test User", "123 Test St")
	if err != nil {
		t.Fatalf("NewUser failed: %v", err)
	}

	users := ListUsers()

	// Find our test user
	var testUser []string
	for _, u := range users {
		if len(u) > 0 && u[0] == user.Username {
			testUser = u
			break
		}
	}

	if testUser == nil {
		t.Fatal("Could not find test user in results")
	}

	// Verify column order: Username, Name, Address, Balance, Days Since Last Fee
	if len(testUser) != 5 {
		t.Fatalf("Expected 5 columns, got %d", len(testUser))
	}

	// Column 0: Username
	if testUser[0] != user.Username {
		t.Errorf("Column 0 should be username, got %s", testUser[0])
	}

	// Column 1: Name
	if testUser[1] != user.GetName() {
		t.Errorf("Column 1 should be name, got %s", testUser[1])
	}

	// Column 2: Address
	if testUser[2] != user.Address {
		t.Errorf("Column 2 should be address, got %s", testUser[2])
	}

	// Column 3: Balance
	if testUser[3] != user.GetBalance() {
		t.Errorf("Column 3 should be balance, got %s", testUser[3])
	}

	// Column 4: Days Since Last Fee
	if testUser[4] != user.GetDaysSinceLastFee() {
		t.Errorf("Column 4 should be days since last fee, got %s", testUser[4])
	}
}

func TestGetUserTable(t *testing.T) {
	// Skip this test in non-GUI environment
	// Fyne widgets require an app context which we can't create in tests
	t.Skip("Skipping Fyne widget test - requires GUI context")
}

func TestGetUserTableWithUsers(t *testing.T) {
	t.Skip("Skipping Fyne widget test - requires GUI context")
}

func TestGetUserTableColumnWidths(t *testing.T) {
	t.Skip("Skipping Fyne widget test - requires GUI context")
}

func TestListUsersNilHandling(t *testing.T) {
	// Test that ListUsers handles nil users gracefully
	// This is more of a defensive test

	users := ListUsers()

	// Verify no panics and result is valid
	if users == nil {
		t.Error("Expected non-nil result")
	}

	// All entries should be valid (non-nil or properly initialized)
	for i, user := range users {
		if user == nil {
			// Skip nil entries - this is acceptable
			continue
		}

		// If not nil, should have proper structure
		if len(user) != 5 {
			t.Errorf("User %d has invalid structure: %d columns", i, len(user))
		}
	}
}

func TestListUsersDataIntegrity(t *testing.T) {
	// Create a user with specific data
	testName := "Data Integrity Test"
	testAddress := "999 Integrity Lane"

	user, err := membership.NewUser(testName, testAddress)
	if err != nil {
		t.Fatalf("NewUser failed: %v", err)
	}

	// Add a fee to test balance calculation
	user.AddFee(25.00, "Test fee")

	users := ListUsers()

	// Find our user
	var found []string
	for _, u := range users {
		if len(u) > 0 && u[0] == user.Username {
			found = u
			break
		}
	}

	if found == nil {
		t.Fatal("Could not find test user")
	}

	// Verify the data matches
	expectedName := user.GetName()
	if found[1] != expectedName {
		t.Errorf("Name mismatch: expected %s, got %s", expectedName, found[1])
	}

	if found[2] != testAddress {
		t.Errorf("Address mismatch: expected %s, got %s", testAddress, found[2])
	}

	// Balance should include the fee
	expectedBalance := user.GetBalance()
	if found[3] != expectedBalance {
		t.Errorf("Balance mismatch: expected %s, got %s", expectedBalance, found[3])
	}
}
