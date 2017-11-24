package gor

// Next go to next handler
type Next struct {
	count int
}

// Next next
func (n *Next) Next() {
	n.count++
}
