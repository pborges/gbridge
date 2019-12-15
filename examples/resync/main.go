package main

import "github.com/pborges/gbridge/homegraph"

func main() {
	client, err := homegraph.Dial("gbridge-58f4081c3cea.json")
	if err != nil {
		panic(err)
	}
	err = client.RequestResync("pborges")
	if err != nil {
		panic(err)
	}
}
