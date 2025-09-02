package interfaces

// Factory creates instances that implement the core interfaces
type Factory struct{}

// NewFactory creates a new factory instance
func NewFactory() *Factory {
	return &Factory{}
}

// TODO: Implement factory methods when concrete implementations 
// are refactored to use interfaces
// This file serves as a placeholder for future SOLID refactoring