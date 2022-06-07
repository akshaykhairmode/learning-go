package main

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Queue interface { //Mockery will use this interface to generate mocks
	Do(string, ...interface{}) (interface{}, error)
}

type RedisQueue struct { //The struct which we will use that is an wrapper of redis connection.
	conn Queue
	name string
}

func (q RedisQueue) PopWithRetry(numRetries int) (string, error) {

	var data string
	var err error

	for i := 0; i < numRetries; i++ {
		data, err = redis.String(q.conn.Do("LPOP", q.name))
		if err == nil {
			return data, err
		}

		log.Printf("LPOP Failed, retry count : %v, sleeping for 100ms", i+1)
		time.Sleep(100 * time.Millisecond)
	}

	return data, err
}

func (q RedisQueue) Pop() (string, error) { //Our pop method. Does a left pop.

	var data string
	var err error

	data, err = redis.String(q.conn.Do("LPOP", q.name))

	return data, err
}

func (q RedisQueue) Push(data string) error { //Our push method. Does a right push.

	_, err := q.conn.Do("RPUSH", q.name, data)
	return err
}

func main() {

	conn, err := redis.Dial("tcp", "localhost:6379") //Creating connection.
	if err != nil {
		log.Fatal(err)
	}

	queue := RedisQueue{
		conn: conn,
		name: "my-test-queue",
	}

	log.Println(queue.Push("test-data"))
	log.Println(queue.Pop())
	log.Println(queue.PopWithRetry(5))

}
