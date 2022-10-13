package main

import (
	"log"
	"redis/gredis"
	"redis/rgo"
)

func main() {

	redigoKey, redigoValue := "REDIGO", "https://github.com/gomodule/redigo"
	goRedisKey, goredisValue := "GOREDIS", "https://github.com/go-redis/redis"

	//SET
	log.Println("REDIGO SET", rgo.Set(redigoKey, redigoValue))
	log.Println("GOREDIS SET", gredis.Set(goRedisKey, goredisValue))

	//GET
	val, err := rgo.Get(redigoKey)
	log.Println("REDIGO GET", val, err)
	val, err = gredis.Get(goRedisKey)
	log.Println("GOREDIS GET", val, err)

	hash := "MyTestHash"

	//HSET
	log.Println("REDIGO HSET", rgo.Hset(hash, redigoKey, redigoValue))
	log.Println("GOREDIS HSET", gredis.Hset(hash, goRedisKey, goredisValue))

	//HGET
	val, err = rgo.Hget(hash, redigoKey)
	log.Println("REDIGO GET", val, err)
	val, err = gredis.Hget(hash, goRedisKey)
	log.Println("GOREDIS GET", val, err)

	//HGETALL
	valMap, err := rgo.Hgetall(hash)
	log.Println("REDIGO HGETALL", valMap, err)
	valMap, err = gredis.Hgetall(hash)
	log.Println("REDIGO HGETALL", valMap, err)

}
