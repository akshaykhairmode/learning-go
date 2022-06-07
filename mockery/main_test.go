package main

import (
	"fmt"
	"testing"

	"learning-mockery/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	redisConn.On("Do", "RPUSH", q.name, pushData).Return(nil, fmt.Errorf("some-redis-error"))
	err := q.Push(pushData)
	assert.Error(t, err, "Expected Error")
}

func TestRedisPush_Success(t *testing.T) {

	redisConn := new(mocks.Queue) //create a new object of the mockery queue.

	q := RedisQueue{
		conn: redisConn,
		name: "redis-queue",
	}

	pushData := "some-data"

	//Mockery provides on and return methods.
	//The On method - In the arguments we first pass the function name which we are mocking and then the other parameters. In this case Do is the function which we are calling and other arguments are the ones which are passed when that function is called.
	//The Return method - Here we pass what we expect the Do function to return.
	redisConn.On("Do", "RPUSH", q.name, pushData).Return(nil, nil)
	err := q.Push(pushData)
	assert.NoError(t, err, "Unexpected Error")
}

func TestRedisPopWithRetry_Success(t *testing.T) {
	redisConn := new(mocks.Queue)

	q := RedisQueue{
		conn: redisConn,
		name: "redis-queue",
	}

	redisConn.On("Do", "LPOP", q.name).Return("", fmt.Errorf("some-redis-error")).Once()
	redisConn.On("Do", "LPOP", q.name).Return("", fmt.Errorf("some-redis-error")).Once()
	redisConn.On("Do", "LPOP", q.name).Return("", fmt.Errorf("some-redis-error")).Once()
	redisConn.On("Do", "LPOP", q.name).Return("", fmt.Errorf("some-redis-error")).Once()
	redisConn.On("Do", "LPOP", q.name).Return("some-data", nil).Once()

	str, err := q.PopWithRetry(5)
	assert.Equal(t, "some-data", str)
	assert.NoError(t, err, "Unexpected Error")
	mock.AssertExpectationsForObjects(t, redisConn)
}
