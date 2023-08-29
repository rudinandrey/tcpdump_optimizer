package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var buffer map[string]int
var currentDate string

func ConvertIpWithPort(ip string) (string, error) {
	fields := strings.Split(ip, ".")
	if len(fields) != 5 {
		return "", errors.New("Не правильный IP адрес с портом")
	}
	return fields[0]+"."+fields[1]+"."+fields[2]+"."+fields[3]+":"+fields[4], nil
}

func AddPacket(line string) error {

	fields := strings.Fields(string(line))
	currentDate = fields[0]
	
	if(len(fields) == 8 && fields[7] != "0" && fields[6] == "tcp") {
		// tcp
		date := fields[0]
		time := strings.Split(fields[1], ".")[0]
		src, _ := ConvertIpWithPort(fields[3])
		dst, _ := ConvertIpWithPort(strings.Replace(fields[5], ":", "", 1))
		proto := fields[6]
		length, err := strconv.Atoi(fields[7])

		if err != nil {
			return err
		}

		key := date+" "+time+" "+src+" "+dst+" "+proto

		_, ok := buffer[key]

		if ok {
			buffer[key] += length
		} else {
			buffer[key] = length
		}
	}

	if(len(fields) == 9 && fields[6] == "UDP,") {
		date := fields[0]
		time := strings.Split(fields[1], ".")[0]
		src := fields[3]
		dst := strings.Replace(fields[5], ":", "", 1)
		proto := strings.Replace(fields[6], "UDP,", "udp", 1)
		length, err := strconv.Atoi(fields[8])

		if err != nil {
			return err
		}

		key := date+" "+time+" "+src+" "+dst+" "+proto

		_, ok := buffer[key]

		if ok {
			buffer[key] += length
		} else {
			buffer[key] = length
		}
	}
	return nil
}

func WriteFile(data *map[string]int) error {
	fileName := "/var/www/traffic/logs/"+currentDate+".log"
	// fileName := currentDate+".log"
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) 
	if err != nil {
		return err
	}
	defer f.Close()

	var sb strings.Builder

	for key, value := range *data {
		text := key + " " + fmt.Sprint(value) + "\n"
		sb.WriteString(text)		
	}

	_, err = f.WriteString(sb.String()) 

	if err != nil {
		return err
	}

	return nil
}

func main() {


	buffer = make(map[string]int)

	reader := bufio.NewReader(os.Stdin)

	for {		
		line, _, err := reader.ReadLine()

		if err != nil {
			log.Fatal(err)
		}

		err = AddPacket(string(line))

		if err != nil {
			log.Fatal(err)
		}

		if len(buffer) > 1024 {
			copyBuffer := make(map[string]int)
			for k, v := range(buffer) {
				copyBuffer[k] = v
			}

			go WriteFile(&copyBuffer)

			// err := WriteFile(&buffer)
			// if err != nil {
			// 	log.Panicln(err)
			// }
			buffer = make(map[string]int)
		}
	}
}