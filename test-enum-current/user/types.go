// Package user provides primitives for OpenAPI-based schema validation.
//
// Code generated from OpenAPI specification. DO NOT EDIT manually.
package user


// UpdateEmailRequest 
type UpdateEmailRequest struct {
	// Valid email address
	NewEmail EmailAddress `json:"newEmail"`
}

// UpdateStatusRequest 
type UpdateStatusRequest struct {
	// User account status
	NewStatus UserStatus `json:"newStatus"`
	// Reason for status change
	Reason string `json:"reason,omitempty"`
}

// UserProfile 
type UserProfile struct {
	// Account creation timestamp
	CreatedAt string `json:"createdAt"`
	// Valid email address
	Email EmailAddress `json:"email"`
	// Unique user identifier
	Id UserId `json:"id"`
	// Last login timestamp
	LastLogin string `json:"lastLogin,omitempty"`
	// Full name
	Name string `json:"name"`
	// User account status
	Status UserStatus `json:"status"`
}



// EmailAddress Valid email address
type EmailAddress = string

// UserId Unique user identifier
type UserId = string

// UserStatus User account status
type UserStatus = string
