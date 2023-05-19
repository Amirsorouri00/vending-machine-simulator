package vending_machine

type Item struct {
	name  string
	price int
}

type State int64

const (
	Idle State = iota
	SelectProduct
	Broken
)

type GetBalance func() int
type GetState func() State
type ChangeState func(st State)

type VendingMachine struct {
	id         int
	state      State
	items      []Item
	stocks     map[string]int
	balance    int
	records    map[string]int
	itemPrices map[string]int

	GetBalance  GetBalance
	GetState    GetState
	ChangeState ChangeState
}
