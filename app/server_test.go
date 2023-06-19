package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func waitForServer() {
	backoff := 50 * time.Millisecond
	log.Println("Waiting for server to start")
	for i := 0; i < 10; i++ {
		conn, err := net.DialTimeout("tcp", ":6379", 1*time.Second)
		if err != nil {
			time.Sleep(backoff)
			continue
		}
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	log.Fatal("Unable to start server")
}

func TestMain(m *testing.M) {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	go main()
	waitForServer()
	fmt.Println("Server started")
	m.Run()
	redisClient.Close()
}

func TestPing(t *testing.T) {
	resp := redisClient.Ping()
	if strings.ToUpper(resp.Val()) != "PONG" {
		t.Fail()
	}
}

func TestEcho(t *testing.T) {
	resp := redisClient.Echo("Hello")
	if strings.ToUpper(resp.Val()) != "HELLO" {
		t.Fail()
	}
}

func TestGetSet(t *testing.T) {
	setResp := redisClient.Set("Hello", "World", 0)
	if setResp.Err() != nil {
		t.Log(setResp.Err().Error())
		t.Fail()
	}
	if strings.ToUpper(setResp.Val()) != "OK" {
		t.Logf("Expected OK but received %s", setResp.Val())
		t.Fail()
	}

	getResp := redisClient.Get("Hello")
	if getResp.Err() != nil {
		t.Log(setResp.Err().Error())
		t.Fail()
	}
	if getResp.Val() != "World" {
		t.Logf("Expected World but received %s", getResp.Val())
		t.Fail()
	}
}

func TestGetSetWithExpiry(t *testing.T) {
	setResp := redisClient.Set("Hello", "World", 100*time.Millisecond)
	if setResp.Err() != nil {
		t.Log(setResp.Err().Error())
		t.Fail()
	}
	if strings.ToUpper(setResp.Val()) != "OK" {
		t.Logf("Expected OK but received %s", setResp.Val())
		t.Fail()
	}

	t.Logf("Sleeping for 100ms")
	time.Sleep(100 * time.Millisecond)

	getResp := redisClient.Get("Hello")
	if getResp.Err() != redis.Nil {
		t.Logf("Expected null string. Got %#v", getResp)
		t.Fail()
	}
}
