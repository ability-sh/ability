package abi

func Login() {

	registry := GetRegistry()

	err := registry.Auth()

	if err != nil {
		Panicln(err)
	}

	Println("SUCCESS")

}
