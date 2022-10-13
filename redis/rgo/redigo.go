package rgo

import (
	"fmt"
	"log"

	rgo "github.com/gomodule/redigo/redis"
)

var conn, connErr = rgo.Dial("tcp", "localhost:6379")

func init() {
	if connErr != nil {
		log.Fatal(connErr)
	}
}

func Set(key, value string) error {

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return fmt.Errorf("error while doing SET command : %v", err)
	}

	return err

}

func Get(key string) (string, error) {

	//String is a helper method
	val, err := rgo.String(conn.Do("GET", key))
	if err != nil {
		return "", fmt.Errorf("error while doing GET command : %v", err)
	}

	return val, err
}

func Hset(hash, key, value string) error {

	_, err := conn.Do("HSET", hash, key, value)
	if err != nil {
		return fmt.Errorf("error while doing HSET command : %v", err)
	}
	return err
}

func Hget(hash, key string) (string, error) {

	val, err := rgo.String(conn.Do("HGET", hash, key))
	if err != nil {
		return "", fmt.Errorf("error while doing HGET command : %v", err)
	}
	return val, err
}

func Hgetall(hash string) (map[string]string, error) {
	val, err := rgo.StringMap(conn.Do("HGETALL", hash))
	if err != nil {
		return nil, fmt.Errorf("error while doing HGETALL command : %v", err)
	}
	return val, err
}
