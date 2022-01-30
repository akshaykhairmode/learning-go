package main

import (
	"fmt"
	"testing"

	"learning-mockery/mocks"

	"github.com/stretchr/testify/assert"
)

func TestRedisPop_Success(t *testing.T) {
	redisConn := new(mocks.Queue)

	q := RedisQueue{
		conn: redisConn,
		name: "redis-queue",
	}

	redisConn.On("Do", "LPOP", q.name).Return("some-data", nil)
	str, err := q.Pop()
	assert.Equal(t, "some-data", str)
	assert.NoError(t, err, "Unexpected Error")
}

func TestRedisPop_Failure(t *testing.T) {
	redisConn := new(mocks.Queue)

	q := RedisQueue{
		conn: redisConn,
		name: "redis-queue",
	}

	redisConn.On("Do", "LPOP", q.name).Return("some-data", fmt.Errorf("some-redis-error"))
	str, err := q.Pop()
	assert.Equal(t, "", str)
	assert.Error(t, err, "Expected Error")
}

func TestRedisPush_Failure(t *testing.T) {
	redisConn := new(mocks.Queue)

	q := RedisQueue{
		conn: redisConn,
		name: "redis-queue",
	}

	pushData := "some-data"

	redisConn.On("Do", "RPUSH", q.name, pushData).Return(fmt.Errorf("some-redis-error"))
	err := q.Push(pushData)
	assert.Error(t, err, "Expected Error")
}

func TestRedisPush_Success(t *testing.T) {
	redisConn := new(mocks.Queue)

	q := RedisQueue{
		conn: redisConn,
		name: "redis-queue",
	}

	pushData := "some-data"

	redisConn.On("Do", "RPUSH", q.name, pushData).Return(nil, nil)
	err := q.Push(pushData)
	assert.NoError(t, err, "Unexpected Error")
}
