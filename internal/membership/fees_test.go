package membership

import (
	"testing"
)

// Helper to reset standard fees between tests
func resetStandardFees() {
	standardFees = make(map[string]*StandardFee)
}

func TestNewStandardFee(t *testing.T) {
	resetStandardFees()

	fee, err := NewStandardFee("Cart Rental", 25.00)
	if err != nil {
		t.Fatalf("NewStandardFee failed: %v", err)
	}

	if fee == nil {
		t.Fatal("Expected non-nil fee")
	}

	if fee.Name != "Cart Rental" {
		t.Errorf("Expected name 'Cart Rental', got %s", fee.Name)
	}

	if fee.Amount != 25.00 {
		t.Errorf("Expected amount 25.00, got %.2f", fee.Amount)
	}

	if fee.ID != "cart_rental" {
		t.Errorf("Expected ID 'cart_rental', got %s", fee.ID)
	}

	// Verify it was added to the map
	if !StandardFeeExists("cart_rental") {
		t.Error("Expected fee to exist in map")
	}
}

func TestNewStandardFeeDuplicateID(t *testing.T) {
	resetStandardFees()

	fee1, err := NewStandardFee("Cart Rental", 25.00)
	if err != nil {
		t.Fatalf("NewStandardFee failed: %v", err)
	}
	fee2, err := NewStandardFee("Cart Rental", 30.00)
	if err != nil {
		t.Fatalf("NewStandardFee failed: %v", err)
	}

	if fee1.ID == fee2.ID {
		t.Error("Expected different IDs for duplicate names")
	}

	if fee2.ID != "cart_rental_1" {
		t.Errorf("Expected ID 'cart_rental_1' for duplicate, got %s", fee2.ID)
	}
}

func TestAddStandardFee(t *testing.T) {
	resetStandardFees()

	fee := &StandardFee{
		ID:     "test_fee",
		Name:   "Test Fee",
		Amount: 15.00,
	}

	AddStandardFee(fee)

	if !StandardFeeExists("test_fee") {
		t.Error("Expected fee to exist after AddStandardFee")
	}

	retrieved, ok := GetStandardFee("test_fee")
	if !ok {
		t.Fatal("Expected to retrieve fee")
	}

	if retrieved.Name != "Test Fee" {
		t.Errorf("Expected name 'Test Fee', got %s", retrieved.Name)
	}
}

func TestAddStandardFeeNil(t *testing.T) {
	resetStandardFees()

	// Should not panic
	AddStandardFee(nil)

	// Should not panic with empty ID
	AddStandardFee(&StandardFee{ID: "", Name: "Test", Amount: 10})
}

func TestAddStandardFees(t *testing.T) {
	resetStandardFees()

	fees := []*StandardFee{
		{ID: "fee1", Name: "Fee 1", Amount: 10.00},
		{ID: "fee2", Name: "Fee 2", Amount: 20.00},
		{ID: "fee3", Name: "Fee 3", Amount: 30.00},
	}

	AddStandardFees(fees)

	if !StandardFeeExists("fee1") {
		t.Error("Expected fee1 to exist")
	}
	if !StandardFeeExists("fee2") {
		t.Error("Expected fee2 to exist")
	}
	if !StandardFeeExists("fee3") {
		t.Error("Expected fee3 to exist")
	}
}

func TestAddStandardFeesDuplicateSkipped(t *testing.T) {
	resetStandardFees()

	// Add initial fee
	_, err := NewStandardFee("Original", 10.00)
	if err != nil {
		t.Fatalf("NewStandardFee failed: %v", err)
	}

	// Try to add duplicate
	fees := []*StandardFee{
		{ID: "original", Name: "Duplicate", Amount: 20.00},
	}

	AddStandardFees(fees)

	// Original should remain
	fee, _ := GetStandardFee("original")
	if fee.Name != "Original" {
		t.Errorf("Expected original name 'Original', got %s", fee.Name)
	}
	if fee.Amount != 10.00 {
		t.Errorf("Expected original amount 10.00, got %.2f", fee.Amount)
	}
}

func TestGetStandardFee(t *testing.T) {
	resetStandardFees()

	_, err := NewStandardFee("Test Fee", 15.00)
	if err != nil {
		t.Fatalf("NewStandardFee failed: %v", err)
	}

	fee, ok := GetStandardFee("test_fee")
	if !ok {
		t.Fatal("Expected to find fee")
	}

	if fee.Name != "Test Fee" {
		t.Errorf("Expected name 'Test Fee', got %s", fee.Name)
	}

	_, ok = GetStandardFee("nonexistent")
	if ok {
		t.Error("Expected false for non-existent fee")
	}
}

func TestGetStandardFees(t *testing.T) {
	resetStandardFees()

	// Test with no fees
	fees := GetStandardFees()
	if len(fees) != 0 {
		t.Errorf("Expected 0 fees, got %d", len(fees))
	}

	// Add some fees
	_, err := NewStandardFee("Zebra Fee", 10.00)
	if err != nil {
		t.Fatalf("NewStandardFee failed: %v", err)
	}
	_, err = NewStandardFee("Apple Fee", 20.00)
	if err != nil {
		t.Fatalf("NewStandardFee failed: %v", err)
	}
	_, err = NewStandardFee("Middle Fee", 15.00)
	if err != nil {
		t.Fatalf("NewStandardFee failed: %v", err)
	}

	fees = GetStandardFees()
	if len(fees) != 3 {
		t.Fatalf("Expected 3 fees, got %d", len(fees))
	}

	// Verify they are sorted by name
	if fees[0].Name != "Apple Fee" {
		t.Errorf("Expected first fee to be 'Apple Fee', got %s", fees[0].Name)
	}
	if fees[1].Name != "Middle Fee" {
		t.Errorf("Expected second fee to be 'Middle Fee', got %s", fees[1].Name)
	}
	if fees[2].Name != "Zebra Fee" {
		t.Errorf("Expected third fee to be 'Zebra Fee', got %s", fees[2].Name)
	}
}

func TestStandardFeeExists(t *testing.T) {
	resetStandardFees()

	if StandardFeeExists("test_fee") {
		t.Error("Expected false for non-existent fee")
	}

	_, err := NewStandardFee("Test Fee", 10.00)
	if err != nil {
		t.Fatalf("NewStandardFee failed: %v", err)
	}

	if !StandardFeeExists("test_fee") {
		t.Error("Expected true for existing fee")
	}
}

func TestUpdateStandardFee(t *testing.T) {
	resetStandardFees()

	_, err := NewStandardFee("Original", 10.00)
	if err != nil {
		t.Fatalf("NewStandardFee failed: %v", err)
	}

	success := UpdateStandardFee("original", "Updated Name", 25.00)
	if !success {
		t.Fatal("Expected update to succeed")
	}

	fee, _ := GetStandardFee("original")
	if fee.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %s", fee.Name)
	}
	if fee.Amount != 25.00 {
		t.Errorf("Expected amount 25.00, got %.2f", fee.Amount)
	}

	// ID should remain the same
	if fee.ID != "original" {
		t.Errorf("Expected ID 'original', got %s", fee.ID)
	}
}

func TestUpdateStandardFeeNonExistent(t *testing.T) {
	resetStandardFees()

	success := UpdateStandardFee("nonexistent", "Test", 10.00)
	if success {
		t.Error("Expected update to fail for non-existent fee")
	}
}

func TestDeleteStandardFee(t *testing.T) {
	resetStandardFees()

	_, err := NewStandardFee("To Delete", 10.00)
	if err != nil {
		t.Fatalf("NewStandardFee failed: %v", err)
	}

	if !StandardFeeExists("to_delete") {
		t.Fatal("Expected fee to exist before deletion")
	}

	success := DeleteStandardFee("to_delete")
	if !success {
		t.Fatal("Expected deletion to succeed")
	}

	if StandardFeeExists("to_delete") {
		t.Error("Expected fee to not exist after deletion")
	}
}

func TestDeleteStandardFeeNonExistent(t *testing.T) {
	resetStandardFees()

	success := DeleteStandardFee("nonexistent")
	if success {
		t.Error("Expected deletion to fail for non-existent fee")
	}
}

func TestNewStandardFeeEmptyName(t *testing.T) {
	resetStandardFees()

	// Test with empty name
	_, err := NewStandardFee("", 25.00)
	if err == nil {
		t.Error("Expected error for empty fee name")
	}

	// Test with whitespace-only name
	_, err = NewStandardFee("   ", 25.00)
	if err == nil {
		t.Error("Expected error for whitespace-only fee name")
	}
}

func TestNewStandardFeeInvalidAmount(t *testing.T) {
	resetStandardFees()

	// Test with zero amount
	_, err := NewStandardFee("Test Fee", 0)
	if err == nil {
		t.Error("Expected error for zero amount")
	}

	// Test with negative amount
	_, err = NewStandardFee("Test Fee", -10.00)
	if err == nil {
		t.Error("Expected error for negative amount")
	}
}

func TestStandardFeeIDGeneration(t *testing.T) {
	resetStandardFees()

	tests := []struct {
		name       string
		expected   string
	}{
		{"Simple Name", "simple_name"},
		{"Multiple   Spaces", "multiple___spaces"},
		{"UPPERCASE", "uppercase"},
		{"Mixed-Case-Name", "mixed-case-name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetStandardFees()
			fee, err := NewStandardFee(tt.name, 10.00)
			if err != nil {
				t.Fatalf("NewStandardFee failed: %v", err)
			}
			if fee.ID != tt.expected {
				t.Errorf("Expected ID %s, got %s", tt.expected, fee.ID)
			}
		})
	}
}
