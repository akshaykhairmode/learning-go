package main

import (
	"log"
	"redis/rgo"
)

func Expiry() {

	//SETEX
	err := rgo.SetEx("Name", "10", "John")
	log.Println("REDIGO SETEX", "Name", err)
	ttl("Name")

	//SETWITHEX
	err = rgo.SetWithExpiry("Age", "15", "34")
	log.Println("REDIGO SET with EX", "Age", err)
	ttl("Age")

	//EXPIRE
	err = rgo.Expire("Age", "30")
	log.Println("REDIGO EXPIRE", "Age", err)
	ttl("Age")
}

func ttl(key string) {
	ttl, err := rgo.TTL(key)
	log.Println("REDIGO TTL", key, ttl, err)
}
