package main

import (
	"avitoRedis/server/cache"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func HandleConnection(cache *cache.Cache, conn net.Conn) {
	defer conn.Close()
	for {
		request := make([]byte, 128)
		n, err := conn.Read(request)
		if err != nil {
			fmt.Println(err)
			break
		}
		msg := string(request[0:n])
		cmds := strings.Fields(msg)
		method := cmds[0]
		switch method {
		case "GET":
			response, err := cache.Get(cmds[1])
			if err != nil {
				conn.Write([]byte(err.Error()))
			} else {
				conn.Write([]byte(response))
			}
		case "SET":
			if len(cmds) == 3 {
				cache.Set(cmds[1], 0, cmds[2])
			} else {
				expT, _ := strconv.Atoi(cmds[4])
				cache.Set(cmds[1], int64(expT), cmds[2])
			}
			conn.Write([]byte("OK!"))
		case "HSET":
			cache.HSet(cmds[1], cmds[2], cmds[3])
			conn.Write([]byte("OK!"))
		case "HGET":
			response, err := cache.HGet(cmds[1], cmds[2])
			if err != nil {
				conn.Write([]byte(err.Error()))
			} else {
				conn.Write([]byte(response))
			}
		case "HGETALL":
			keys, err := cache.HGetAll(cmds[1])
			if err != nil {
				conn.Write([]byte(err.Error()))
			} else {
				var response []byte
				for _, key := range keys {
					response = append(response, []byte(key+"\n")...)
				}
				conn.Write(response[0 : len(response)-2])
			}
		case "LPUSH":
			response, err := cache.LPush(cmds[1], cmds[2])
			if err != nil {
				conn.Write([]byte(err.Error()))
			} else {
				conn.Write([]byte(strconv.Itoa(response)))
			}
		case "RPUSH":
			response, err := cache.RPush(cmds[1], cmds[2])
			if err != nil {
				conn.Write([]byte(err.Error()))
			} else {
				conn.Write([]byte(strconv.Itoa(response)))
			}
		case "LRANGE":
			start, _ := strconv.Atoi(cmds[2])
			end, _ := strconv.Atoi(cmds[3])
			elems, err := cache.LRange(cmds[1], start, end)
			if err != nil {
				conn.Write([]byte(err.Error()))
			} else {
				var response []byte
				for _, el := range elems {
					response = append(response, []byte(el+"\n")...)
				}
				conn.Write(response[0 : len(response)-1])
			}
		case "LSET":
			ind, _ := strconv.Atoi(cmds[2])
			err := cache.LSet(cmds[1], cmds[3], ind)
			if err != nil {
				conn.Write([]byte(err.Error()))
			} else {
				conn.Write([]byte("OK!"))
			}
		case "LGET":
			ind, _ := strconv.Atoi(cmds[2])
			response, err := cache.LGet(cmds[1], ind)
			if err != nil {
				conn.Write([]byte(err.Error()))
			} else {
				conn.Write([]byte(response))
			}
		case "DEL":
			cache.Delete(cmds[1])
			conn.Write([]byte("OK!"))
		case "KEYS":
			keys := cache.Keys()
			var response []byte
			for _, key := range keys {
				response = append(response, []byte(key+" ")...)
			}
			conn.Write(response)
		case "EXPIRE":
			expt, _ := strconv.Atoi(cmds[2])
			resp := cache.Expire(cmds[1], int64(expt))
			conn.Write([]byte(resp))
		case "TTL":
			response := cache.TTL(cmds[1])
			conn.Write([]byte(strconv.Itoa(int(response))))
		default:
			fmt.Println("Unhandled command")
		}
	}
}

func main() {
	lner, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lner.Close()
	cch := cache.InitCache()
	for {
		conn, err := lner.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Handling connection...")
		go HandleConnection(cch, conn)
	}
}
