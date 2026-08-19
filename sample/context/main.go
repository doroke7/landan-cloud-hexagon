package main

import (
	"context"
	"fmt"
)

func main() {
	oContext, fCancel := context.WithCancel(context.Background())

	oAuthorization := new(string)
	oContext = context.WithValue(oContext, "a", oAuthorization)

	test(oContext)

	sAuthorization := *oAuthorization

	fmt.Println("sAuthorization 18 =", sAuthorization)

	defer fCancel()
}

// context 本身是
func test(oContext context.Context) {

	if oAuthorization, bOk := oContext.Value("a").(*string); bOk {
		*oAuthorization = "12345678"
	}

}
