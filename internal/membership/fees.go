package membership

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Fee represents a fee assigned to a user's account.
// Fees track charges such as cart rentals, guest fees, or penalties.
// Each fee has an amount, reason, paid status, and creation timestamp.
type Fee struct {
	// Amount is the fee charge in dollars.
	Amount float64

	// Reason describes what the fee is for (e.g., "Cart Rental").
	Reason string

	// Paid indicates whether this fee has been paid.
	Paid bool

	// CreatedAt is the timestamp when this fee was created.
	CreatedAt time.Time
}

// NewFee creates a new unpaid fee with the specified amount and reason.
// The fee is automatically marked as unpaid (Paid = false) and timestamped
// with the current time.
func NewFee(amt float64, reason string) Fee {
	return Fee{
		Amount:    amt,
		Reason:    reason,
		Paid:      false,
		CreatedAt: time.Now(),
	}
}

// StandardFee represents a reusable fee template.
// Standard fees can be quickly assigned to users without re-entering
// the fee name and amount each time. Examples include "Cart Rental",
// "Guest Fee", or "Late Payment".
type StandardFee struct {
	// ID is the unique identifier for this fee template, automatically
	// generated from the name (lowercase with underscores).
	ID string

	// Name is the display name for this fee (e.g., "Cart Rental").
	Name string

	// Amount is the standard charge for this fee in dollars.
	Amount float64
}

// Standard fee storage
var standardFees = make(map[string]*StandardFee)

// NewStandardFee creates a new standard fee template with a unique ID.
// The ID is automatically generated from the name (lowercase, spaces replaced
// with underscores). If a fee with the same ID exists, a number suffix is added.
// Returns an error if the name is empty or the amount is not positive.
func NewStandardFee(name string, amount float64) (*StandardFee, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("fee name cannot be empty")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("fee amount must be positive")
	}

	// Generate ID from name (lowercase, no spaces)
	id := strings.ToLower(strings.ReplaceAll(name, " ", "_"))

	// Handle duplicates by adding a number
	baseID := id
	counter := 1
	for StandardFeeExists(id) {
		id = fmt.Sprintf("%s_%d", baseID, counter)
		counter++
	}

	fee := &StandardFee{
		ID:     id,
		Name:   name,
		Amount: amount,
	}
	standardFees[id] = fee
	return fee, nil
}

// AddStandardFee adds a single standard fee to the system.
// This is typically used when loading fees from persistent storage.
// If the fee is nil or has an empty ID, it is ignored.
func AddStandardFee(fee *StandardFee) {
	if fee != nil && fee.ID != "" {
		standardFees[fee.ID] = fee
	}
}

// AddStandardFees adds multiple standard fees to the system from a slice.
// This is typically used when loading fees from persistent storage.
// Fees that already exist (by ID) are skipped. Nil fees are ignored.
func AddStandardFees(fees []*StandardFee) {
	for _, fee := range fees {
		if fee != nil && !StandardFeeExists(fee.ID) {
			standardFees[fee.ID] = fee
		}
	}
}

// GetStandardFee retrieves a standard fee by its ID.
// Returns the StandardFee and true if found, or nil and false if not found.
func GetStandardFee(id string) (*StandardFee, bool) {
	fee, ok := standardFees[id]
	return fee, ok
}

// GetStandardFees returns all standard fees sorted alphabetically by name.
// This provides a consistent ordering for display in menus and lists.
func GetStandardFees() []*StandardFee {
	fees := make([]*StandardFee, 0)
	for _, fee := range standardFees {
		fees = append(fees, fee)
	}
	// Sort by name for consistent ordering
	sort.Slice(fees, func(i, j int) bool {
		return fees[i].Name < fees[j].Name
	})
	return fees
}

// StandardFeeExists checks if a standard fee with the given ID exists.
// Returns true if the fee exists, false otherwise.
func StandardFeeExists(id string) bool {
	_, ok := standardFees[id]
	return ok
}

// UpdateStandardFee updates the name and amount of an existing standard fee.
// The ID cannot be changed. Returns true if the fee was updated, false if
// the fee ID was not found.
func UpdateStandardFee(id, name string, amount float64) bool {
	fee, ok := standardFees[id]
	if ok {
		fee.Name = name
		fee.Amount = amount
		return true
	}
	return false
}

// DeleteStandardFee removes a standard fee from the system by its ID.
// Returns true if the fee was deleted, false if the fee ID was not found.
func DeleteStandardFee(id string) bool {
	if StandardFeeExists(id) {
		delete(standardFees, id)
		return true
	}
	return false
}
