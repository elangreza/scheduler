package main

import (
	"fmt"
	"log"

	"github.com/elangreza/scheduler/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg)

	// to := []string{"babehracing14@gmail.com"}
	// cc := []string{}
	// subject := "Test mail"
	// message := "Hello"

	// err = sendMail(cfg, to, cc, subject, message)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

}
