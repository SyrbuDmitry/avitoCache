package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func handleInput(msg string) error {
	cmds := strings.Fields(msg)
	method := cmds[0]
	switch method {
	case "GET":
		if len(cmds) != 2 {
			return errors.New("GET must be: GET KEY")
		}
	case "SET":
		if len(cmds) != 3 && len(cmds) != 5 {
			return errors.New("SET must be: SET KEY VALUE")
		} else if len(cmds) == 5 {
			if cmds[3] != "ex" {
				return errors.New("SET must be: SET KEY VALUE ex INT")
			}
			_, err := strconv.Atoi(cmds[4])
			if err != nil {
				return errors.New("Expiration time must be Integer")
			}
		}
	case "DELETE":
		if len(cmds) != 2 {
			return errors.New("DELETE must be: DELETE KEY")
		}
	case "KEYS":
		if len(cmds) != 1 {
			return errors.New("KEY INPUT ERROR")
		}
	case "HGET":
		if len(cmds) != 3 {
			return errors.New("HGET must be: HGET KEY FIELD")
		}
	case "HGETALL":
		if len(cmds) != 2 {
			return errors.New("HGETALL must be: HGET KEY")
		}
	case "HSET":
		if len(cmds) != 4 {
			return errors.New("HSET must be: HSET KEY FIELD VALUE")
		}
	case "LPUSH":
		if len(cmds) != 3 {
			return errors.New("LPUSH must be: LPUSH KEY VALUE")
		}
	case "RPUSH":
		if len(cmds) != 3 {
			return errors.New("RPUSH must be: RPUSH KEY VALUE")
		}
	case "LRANGE":
		if len(cmds) != 4 {
			return errors.New("LRANGE must be: LRANGE KEY START_IND END_IND")
		}
		_, err := strconv.Atoi(cmds[2])
		if err != nil {
			return errors.New("Start index must be integer")
		}
		_, err = strconv.Atoi(cmds[3])
		if err != nil {
			return errors.New("End index must be integer")
		}
	case "LGET":
		if len(cmds) != 3 {
			return errors.New("LGET must be: LGET KEY INDEX")
		}
		_, err := strconv.Atoi(cmds[2])
		if err != nil {
			return errors.New("INDEX must be integer")
		}
	case "LSET":
		if len(cmds) != 4 {
			return errors.New("LSET must be: LSET KEY INDEX VALUE")
		}
		_, err := strconv.Atoi(cmds[2])
		if err != nil {
			return errors.New("INDEX must be integer")
		}
	case "EXPIRE":
		if len(cmds) != 3 {
			return errors.New("EXPIRE must be: EXPIRE KEY TIME")
		}
		_, err := strconv.Atoi(cmds[2])
		if err != nil {
			return errors.New("TIME must be integer")
		}
	case "TTL":
		if len(cmds) != 2 {
			return errors.New("TTL must be: TTL KEY")
		}

	default:
		return errors.New("No such command!")
	}
	return nil
}

func main() {
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	i := 0
	sc := bufio.NewScanner(os.Stdin)
	for {
		i++
		sc.Scan()
		msg := sc.Text()
		if err := handleInput(msg); err != nil {
			fmt.Println(err)
			continue
		}
		if n, err := conn.Write([]byte(msg)); err != nil || n == 0 {
			fmt.Println(err)
		}
		buff := make([]byte, 128)
		n, rerr := conn.Read(buff)

		if rerr != nil {
			fmt.Println("Reading error")
			fmt.Println(rerr.Error())
			continue
		}
		fmt.Println(string(buff[0:n]))
	}
}
