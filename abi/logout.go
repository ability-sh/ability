package abi

func Logout() {

	registry := GetRegistry()

	registry.Logout()

	Println("SUCCESS")

}
