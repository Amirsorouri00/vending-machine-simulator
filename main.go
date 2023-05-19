package main

import (
	"fmt"

	vmachine "vm_go/vending_machine"
)

const MACHINE_NUM = 3

func server(vms []*vmachine.VendingMachine, i int, j int) {
	vms[j].InsertCoin(i)
	vms[j].DisplayItems()

	fmt.Print("If you want Coke input 1 and if you want Coffe input 2:")
	var s int
	fmt.Scanf("%v", &s)
	switch s {
	case 1:
		vms[j].SelectItem("Coke")
	case 2:
		vms[j].SelectItem("Coffe")
	default:
		fmt.Println("Wrong input!")
		vms[j].Refund(true)
	}

	fmt.Println("\n🤚 Good Luck!")
}

func main() {
	vms := make([]*vmachine.VendingMachine, 0)
	for i := 1; i <= MACHINE_NUM; i++ {
		vms = append(vms, vmachine.NewVendingMachine(i))
	}
	fmt.Println(" 👋 Welcome to the Vending Machine System!")
	for i := 0; true; i++ {
		var i, j int
		fmt.Print("\n Type/Insert the amount of your coin and the ID of the Vending Machine with space(Ex: 1300 1):")
		fmt.Scanf("%v %v", &i, &j)
		if j >= MACHINE_NUM {
			fmt.Println("wrong input. Try again")
			continue
		}
		server(vms, i, j)
		break
	}
}
