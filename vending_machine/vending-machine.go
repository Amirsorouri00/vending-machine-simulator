package vending_machine

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

var mu sync.Mutex

func NewVendingMachine(id int) *VendingMachine {
	return newVendingMachine(id, false)
}

func NewVendingMachineForTest(id int) *VendingMachine {
	return newVendingMachine(id, true)
}

func newVendingMachine(id int, fortest bool) *VendingMachine {
	vm := &VendingMachine{
		id:    id,
		state: Idle,
		items: []Item{
			{name: "Coke", price: 100},
			{name: "Coffe", price: 50},
		},
		stocks:     make(map[string]int),
		balance:    0,
		records:    make(map[string]int),
		itemPrices: make(map[string]int),
	}
	for _, item := range vm.items {
		vm.stocks[item.name] = 10
		vm.itemPrices[item.name] = item.price
	}

	if fortest {
		vm.GetBalance = func() int {
			return vm.balance
		}
		vm.GetState = func() State {
			return vm.state
		}
		vm.ChangeState = func(st State) {
			vm.state = st
		}
	}
	return vm
}

func (vm *VendingMachine) changeState(st State) error {
	mu.Lock()

	if vm.state != Idle {
		return fmt.Errorf("the machine %v is busy with another process", vm.id)
	}
	vm.state = st

	mu.Unlock()

	return nil
}

func (vm *VendingMachine) validateState() error {
	if vm.state != SelectProduct {
		vm.changeState(Broken)
		return fmt.Errorf("state is invalid. machine %v is broken", vm.id)
	}
	return nil
}

func (vm *VendingMachine) InsertCoin(amount int) error {
	err := vm.changeState(SelectProduct)
	if err != nil {
		return err
	}
	vm.balance += amount
	log.Printf("\nInserted %d cents. Current balance: %d cents.\n", amount, vm.balance)
	return nil
}

func (vm *VendingMachine) DisplayItems() error {
	err := vm.validateState()
	if err != nil {
		vm.Refund(false)
		vm.changeState(Broken)
		return err
	}

	log.Println("--- Items Available ---")
	for _, item := range vm.items {
		log.Printf("%s (Price: %d)\n", item.name, item.price)
	}
	log.Println("-----------------------")

	return nil
}

func (vm *VendingMachine) SelectItem(name string) error {
	err := vm.validateState()
	if err != nil {
		vm.Refund(false)
		vm.changeState(Broken)
		return err
	}

	if vm.stocks[name] == 0 {
		log.Printf("Sorry, %s is out of stock.\n", name)
		vm.Refund(true)
		return errors.New("out of stock")
	}
	price, ok := vm.itemPrices[name]
	if !ok {
		log.Printf("Sorry, %s is not a valid item.\n", name)
		vm.Refund(true)
		return errors.New("invalid item")
	}
	if vm.balance < price {
		log.Printf("Insufficient balance. Please insert at least %d cents.\n", price)
		return errors.New("insufficient balance")
	}

	vm.balance -= price
	vm.stocks[name] -= 1
	vm.records[name] += 1
	log.Printf("Dispensing %s. Current balance: %d cents.\n", name, vm.balance)

	if vm.balance > 0 {
		log.Println("\nGoing to return your additional paid money")
		vm.Refund(true)
		return nil
	}
	vm.changeState(Idle)
	return nil
}

func (vm *VendingMachine) Refund(makeIdle bool) {
	log.Printf("Refunding %d cents.\n", vm.balance)
	vm.balance = 0

	if makeIdle {
		vm.changeState(Idle)
	}
}

func (vm *VendingMachine) DisplaySales() {
	log.Println("--- Sales Report ---")
	total := 0
	for name, count := range vm.records {
		price := vm.itemPrices[name]
		subtotal := price * count
		log.Printf("%s: %d x %d cents = %d cents\n", name, count, price, subtotal)
		total += subtotal
	}
	log.Printf("Total sales: %d cents\n", total)
	log.Println("--------------------")
}
