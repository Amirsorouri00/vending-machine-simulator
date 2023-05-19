package vending_machine

import (
	"io/ioutil"
	"log"
	rand "math/rand"
	"os"
	"testing"
	"time"

	vmachine "vm_go/vending_machine"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	rand.Seed(time.Now().Unix())
	code := m.Run()
	os.Exit(code)
}

func makeMachine(num int) []*vmachine.VendingMachine {
	var vms = make([]*vmachine.VendingMachine, 0)
	for i := 1; i <= num; i++ {
		vms = append(vms, vmachine.NewVendingMachineForTest(i))
	}
	return vms
}

func TestInsertCoinValid(t *testing.T) {
	vm := makeMachine(1)
	vm[0].InsertCoin(1200)
	want := 1200

	if vm[0].GetBalance() != want {
		t.Errorf("got %q, wanted %q", vm[0].GetBalance(), want)
	}
}

func TestInsertCoinInvalid(t *testing.T) {
	vm := makeMachine(1)
	states := []vmachine.State{vmachine.SelectProduct, vmachine.Broken}
	vm[0].ChangeState(states[rand.Intn(len(states))])

	err := vm[0].InsertCoin(1200)
	if nil == err {
		t.Errorf("got nil, wanted err to not be nil")
	}
}

func TestDisplayItemsValid1st(t *testing.T) {
	vm := makeMachine(1)
	err := vm[0].DisplayItems()
	if nil == err {
		t.Errorf("got nil, wanted err to not be nil")
	}
}

func TestDisplayItemsValid2nd(t *testing.T) {
	vm := makeMachine(1)
	vm[0].ChangeState(vmachine.SelectProduct)
	err := vm[0].DisplayItems()
	if nil != err {
		t.Errorf("error should be nil. error is %v", err)
	}
}

func TestDisplayItemsInvalid(t *testing.T) {
	vm := makeMachine(1)
	states := []vmachine.State{vmachine.Idle, vmachine.Broken}
	vm[0].ChangeState(states[rand.Intn(len(states))])

	err := vm[0].DisplayItems()
	if nil == err {
		t.Errorf("error must not be nil when state is Idle or Broken")
	}
}

func TestProcessInputTabularMethod(t *testing.T) {
	cases := []struct {
		name        string
		changeState bool
		coin        int
		itemName    string
		expectError bool
	}{
		{
			name:        "Case 1-valid coke",
			coin:        100,
			itemName:    "Coke",
			expectError: false,
		},
		{
			name:        "Case 2-Invalid coint amount for coke",
			coin:        99,
			itemName:    "Coke",
			expectError: true,
		},
		{
			name:        "Case 3-valid coke",
			coin:        1200,
			itemName:    "Coke",
			expectError: false,
		},
		{
			name:        "Case 4-valid coffe",
			coin:        50,
			itemName:    "Coffe",
			expectError: false,
		},
		{
			name:        "Case 5-Invalid coint amount for coffe",
			coin:        49,
			itemName:    "Coffe",
			expectError: true,
		},
		{
			name:        "Case 6-valid coffe",
			coin:        1300,
			itemName:    "Coffe",
			expectError: false,
		},
		{
			name:        "Case 7-invalid state",
			changeState: true,
			coin:        1300,
			itemName:    "Coffe",
			expectError: true,
		},
		{
			name:        "Case 8-Invalid name",
			changeState: false,
			coin:        1300,
			itemName:    "Coffe2",
			expectError: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			vm := makeMachine(1)
			vm[0].InsertCoin(tc.coin)
			var err error = nil
			if tc.changeState {
				states := []vmachine.State{vmachine.Idle, vmachine.Broken}
				vm[0].ChangeState(states[rand.Intn(len(states))])
				err = vm[0].SelectItem(tc.itemName)
				if err == nil && !tc.expectError {
					t.Fatal(err)
				}
			} else {
				err = vm[0].SelectItem(tc.itemName)
			}

			if err != nil && !tc.expectError {
				t.Fatal(err)
			} else if err == nil && tc.expectError {
				t.Fatal(err)
			}
		})
	}
}
