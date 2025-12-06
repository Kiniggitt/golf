package registrar

import (
	"os"
	"path/filepath"
	"testing"

	"golf/internal/membership"
)

func TestReadJSON(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_users.json")

	// Write test data
	testData := `[
    {
        "Username": "testuser",
        "First": "Test",
        "Last": "User",
        "Suffix": "",
        "Address": "123 Test St",
        "Balance": 10.5,
        "Fees": [
            {
                "Amount": 5.0,
                "Reason": "Test fee",
                "Paid": false
            }
        ]
    }
]`

	err := os.WriteFile(testFile, []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test reading the file
	users, err := ReadJSON(testFile)
	if err != nil {
		t.Fatalf("ReadJSON failed: %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("Expected 1 user, got %d", len(users))
	}

	user := users[0]
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", user.Username)
	}

	if user.First != "Test" {
		t.Errorf("Expected first name 'Test', got %s", user.First)
	}

	if user.Last != "User" {
		t.Errorf("Expected last name 'User', got %s", user.Last)
	}

	if user.Address != "123 Test St" {
		t.Errorf("Expected address '123 Test St', got %s", user.Address)
	}

	if user.Balance != 10.5 {
		t.Errorf("Expected balance 10.5, got %.2f", user.Balance)
	}

	if len(user.Fees) != 1 {
		t.Fatalf("Expected 1 fee, got %d", len(user.Fees))
	}

	if user.Fees[0].Amount != 5.0 {
		t.Errorf("Expected fee amount 5.0, got %.2f", user.Fees[0].Amount)
	}

	if user.Fees[0].Reason != "Test fee" {
		t.Errorf("Expected fee reason 'Test fee', got %s", user.Fees[0].Reason)
	}

	if user.Fees[0].Paid {
		t.Error("Expected fee to be unpaid")
	}
}

func TestReadJSONNonExistentFile(t *testing.T) {
	// Reading a non-existent file should return an error
	users, err := ReadJSON("nonexistent_file.json")

	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	if users == nil {
		t.Error("Expected non-nil slice even on error")
	}

	if len(users) != 0 {
		t.Errorf("Expected empty slice, got %d users", len(users))
	}
}

func TestReadJSONEmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "empty.json")

	// Write empty array
	err := os.WriteFile(testFile, []byte("[]"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	users, err := ReadJSON(testFile)
	if err != nil {
		t.Fatalf("ReadJSON failed: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("Expected 0 users from empty file, got %d", len(users))
	}
}

func TestReadJSONInvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "invalid.json")

	// Write invalid JSON
	err := os.WriteFile(testFile, []byte("{invalid json}"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Should return error on invalid JSON
	users, err := ReadJSON(testFile)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	if len(users) != 0 {
		t.Errorf("Expected 0 users from invalid JSON, got %d", len(users))
	}
}

func TestWriteJSON(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "output.json")

	// Create test users
	users := []*membership.User{
		{
			Username: "alice",
			First:    "Alice",
			Last:     "Smith",
			Suffix:   "",
			Address:  "111 First St",
			Balance:  20.0,
			Fees: []membership.Fee{
				{Amount: 10.0, Reason: "New account fee", Paid: false},
			},
		},
		{
			Username: "bob",
			First:    "Bob",
			Last:     "Jones",
			Suffix:   "Jr.",
			Address:  "222 Second St",
			Balance:  0.0,
			Fees:     []membership.Fee{},
		},
	}

	// Write the users to file
	err := WriteJSON(testFile, users)
	if err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	// Verify the file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Expected file to be created")
	}

	// Read the file back
	readUsers, err := ReadJSON(testFile)
	if err != nil {
		t.Fatalf("ReadJSON failed: %v", err)
	}

	if len(readUsers) != 2 {
		t.Fatalf("Expected 2 users in written file, got %d", len(readUsers))
	}

	// Verify first user
	if readUsers[0].Username != "alice" {
		t.Errorf("Expected username 'alice', got %s", readUsers[0].Username)
	}

	if readUsers[0].First != "Alice" {
		t.Errorf("Expected first name 'Alice', got %s", readUsers[0].First)
	}

	if readUsers[0].Balance != 20.0 {
		t.Errorf("Expected balance 20.0, got %.2f", readUsers[0].Balance)
	}

	if len(readUsers[0].Fees) != 1 {
		t.Fatalf("Expected 1 fee for first user, got %d", len(readUsers[0].Fees))
	}

	// Verify second user
	if readUsers[1].Username != "bob" {
		t.Errorf("Expected username 'bob', got %s", readUsers[1].Username)
	}

	if readUsers[1].Suffix != "Jr." {
		t.Errorf("Expected suffix 'Jr.', got %s", readUsers[1].Suffix)
	}

	if len(readUsers[1].Fees) != 0 {
		t.Errorf("Expected 0 fees for second user, got %d", len(readUsers[1].Fees))
	}
}

func TestWriteJSONEmptySlice(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "empty_output.json")

	users := []*membership.User{}

	err := WriteJSON(testFile, users)
	if err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	// Read back and verify
	readUsers, err := ReadJSON(testFile)
	if err != nil {
		t.Fatalf("ReadJSON failed: %v", err)
	}
	if len(readUsers) != 0 {
		t.Errorf("Expected 0 users, got %d", len(readUsers))
	}
}

func TestWriteJSONFormatting(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "formatted.json")

	users := []*membership.User{
		{
			Username: "test",
			First:    "Test",
			Last:     "User",
			Address:  "123 Test St",
			Balance:  0.0,
			Fees:     []membership.Fee{},
		},
	}

	err := WriteJSON(testFile, users)
	if err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	// Read the file as raw bytes to check formatting
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Check that the file is formatted with indentation
	contentStr := string(content)
	if !containsIndentation(contentStr) {
		t.Error("Expected JSON to be formatted with indentation")
	}

	// Verify it starts with array bracket
	if contentStr[0] != '[' {
		t.Error("Expected JSON to start with '['")
	}
}

func TestRoundTrip(t *testing.T) {
	// Test writing and reading preserves data
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "roundtrip.json")

	originalUsers := []*membership.User{
		{
			Username: "user1",
			First:    "First",
			Last:     "User",
			Suffix:   "Sr.",
			Address:  "Address 1",
			Balance:  100.50,
			Fees: []membership.Fee{
				{Amount: 10.0, Reason: "Fee 1", Paid: false},
				{Amount: 5.0, Reason: "Fee 2", Paid: true},
			},
		},
	}

	// Write
	err := WriteJSON(testFile, originalUsers)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Read
	readUsers, err := ReadJSON(testFile)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	// Compare
	if len(readUsers) != len(originalUsers) {
		t.Fatalf("Expected %d users, got %d", len(originalUsers), len(readUsers))
	}

	orig := originalUsers[0]
	read := readUsers[0]

	if orig.Username != read.Username {
		t.Errorf("Username mismatch: %s != %s", orig.Username, read.Username)
	}

	if orig.First != read.First {
		t.Errorf("First name mismatch: %s != %s", orig.First, read.First)
	}

	if orig.Last != read.Last {
		t.Errorf("Last name mismatch: %s != %s", orig.Last, read.Last)
	}

	if orig.Suffix != read.Suffix {
		t.Errorf("Suffix mismatch: %s != %s", orig.Suffix, read.Suffix)
	}

	if orig.Address != read.Address {
		t.Errorf("Address mismatch: %s != %s", orig.Address, read.Address)
	}

	if orig.Balance != read.Balance {
		t.Errorf("Balance mismatch: %.2f != %.2f", orig.Balance, read.Balance)
	}

	if len(orig.Fees) != len(read.Fees) {
		t.Fatalf("Fees count mismatch: %d != %d", len(orig.Fees), len(read.Fees))
	}

	for i := range orig.Fees {
		if orig.Fees[i].Amount != read.Fees[i].Amount {
			t.Errorf("Fee %d amount mismatch: %.2f != %.2f", i, orig.Fees[i].Amount, read.Fees[i].Amount)
		}
		if orig.Fees[i].Reason != read.Fees[i].Reason {
			t.Errorf("Fee %d reason mismatch: %s != %s", i, orig.Fees[i].Reason, read.Fees[i].Reason)
		}
		if orig.Fees[i].Paid != read.Fees[i].Paid {
			t.Errorf("Fee %d paid status mismatch: %v != %v", i, orig.Fees[i].Paid, read.Fees[i].Paid)
		}
	}
}

// Helper function to check if JSON contains indentation
func containsIndentation(s string) bool {
	// Check for common indentation patterns (spaces or tabs after newline)
	return len(s) > 0 && (containsPattern(s, "\n    ") || containsPattern(s, "\n\t"))
}

func containsPattern(s, pattern string) bool {
	for i := 0; i <= len(s)-len(pattern); i++ {
		if s[i:i+len(pattern)] == pattern {
			return true
		}
	}
	return false
}

// Standard Fee Tests

func TestReadStandardFeesJSON(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_fees.json")

	testData := `[
    {
        "ID": "cart_rental",
        "Name": "Cart Rental",
        "Amount": 25.00
    },
    {
        "ID": "guest_fee",
        "Name": "Guest Fee",
        "Amount": 15.00
    }
]`

	err := os.WriteFile(testFile, []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	fees, err := ReadStandardFeesJSON(testFile)
	if err != nil {
		t.Fatalf("ReadStandardFeesJSON failed: %v", err)
	}

	if len(fees) != 2 {
		t.Fatalf("Expected 2 fees, got %d", len(fees))
	}

	// Check first fee
	if fees[0].ID != "cart_rental" {
		t.Errorf("Expected ID 'cart_rental', got %s", fees[0].ID)
	}
	if fees[0].Name != "Cart Rental" {
		t.Errorf("Expected name 'Cart Rental', got %s", fees[0].Name)
	}
	if fees[0].Amount != 25.00 {
		t.Errorf("Expected amount 25.00, got %.2f", fees[0].Amount)
	}

	// Check second fee
	if fees[1].ID != "guest_fee" {
		t.Errorf("Expected ID 'guest_fee', got %s", fees[1].ID)
	}
	if fees[1].Amount != 15.00 {
		t.Errorf("Expected amount 15.00, got %.2f", fees[1].Amount)
	}
}

func TestReadStandardFeesJSONNonExistent(t *testing.T) {
	fees, err := ReadStandardFeesJSON("nonexistent_fees.json")

	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	if fees == nil {
		t.Error("Expected non-nil slice even on error")
	}

	if len(fees) != 0 {
		t.Errorf("Expected empty slice, got %d fees", len(fees))
	}
}

func TestWriteStandardFeesJSON(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "output_fees.json")

	fees := []*membership.StandardFee{
		{ID: "fee1", Name: "Fee 1", Amount: 10.00},
		{ID: "fee2", Name: "Fee 2", Amount: 20.00},
	}

	err := WriteStandardFeesJSON(testFile, fees)
	if err != nil {
		t.Fatalf("WriteStandardFeesJSON failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("Expected file to be created")
	}

	// Read back
	readFees, err := ReadStandardFeesJSON(testFile)
	if err != nil {
		t.Fatalf("ReadStandardFeesJSON failed: %v", err)
	}

	if len(readFees) != 2 {
		t.Fatalf("Expected 2 fees in written file, got %d", len(readFees))
	}

	if readFees[0].ID != "fee1" {
		t.Errorf("Expected ID 'fee1', got %s", readFees[0].ID)
	}
	if readFees[1].Name != "Fee 2" {
		t.Errorf("Expected name 'Fee 2', got %s", readFees[1].Name)
	}
}

func TestStandardFeesRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "roundtrip_fees.json")

	originalFees := []*membership.StandardFee{
		{ID: "test_fee_1", Name: "Test Fee 1", Amount: 100.50},
		{ID: "test_fee_2", Name: "Test Fee 2", Amount: 75.25},
	}

	// Write
	err := WriteStandardFeesJSON(testFile, originalFees)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Read
	readFees, err := ReadStandardFeesJSON(testFile)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	// Compare
	if len(readFees) != len(originalFees) {
		t.Fatalf("Expected %d fees, got %d", len(originalFees), len(readFees))
	}

	for i := range originalFees {
		if originalFees[i].ID != readFees[i].ID {
			t.Errorf("Fee %d ID mismatch: %s != %s", i, originalFees[i].ID, readFees[i].ID)
		}
		if originalFees[i].Name != readFees[i].Name {
			t.Errorf("Fee %d name mismatch: %s != %s", i, originalFees[i].Name, readFees[i].Name)
		}
		if originalFees[i].Amount != readFees[i].Amount {
			t.Errorf("Fee %d amount mismatch: %.2f != %.2f", i, originalFees[i].Amount, readFees[i].Amount)
		}
	}
}
