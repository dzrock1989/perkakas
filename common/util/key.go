package util

type ContextKey int

const (
	ContextClaims ContextKey = iota
	ContextClaimsBytes
)
