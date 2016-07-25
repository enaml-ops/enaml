package main

import (
	"log"

	"github.com/enaml-ops/enaml/enamlbosh"
)

func main() {
	c, err := enamlbosh.NewClient("admin", "eadxz7dh1d4e3bhgx518", "https://52.40.154.174", 25555, true)
	if err != nil {
		log.Fatalln("Error creating client", err)
	}
	log.Println("getting cloud config...")
	cc, err := c.GetCloudConfig()
	if err != nil {
		log.Fatalln("Error getting cloud config:", err)
	}
	log.Println(cc)
}
