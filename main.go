package main

import (
	"log"
	"sber_cloud/tw/cmd"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile | log.LUTC)
}

func main() {
	var err error
	if err = cmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
