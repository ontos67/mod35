package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"strings"
	"time"
)

func main() {
	// Загрузка поговорок из файла
	gpv, err := loadGPV("GoProverbs.txt")
	if err != nil {
		fmt.Println("Ошибка при загрузке поговорок:", err)
		return
	}

	// Запуск службы
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка при запуске службы:", err)
		return
	}
	defer listen.Close()

	fmt.Println("Служба запущена. Ожидание подключений...")

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Ошибка при принятии подключения:", err)
			continue
		}

		go handleClient(conn, gpv)
	}
}

func handleClient(conn net.Conn, gpv []string) {
	defer conn.Close()

	fmt.Printf("Новое подключение от %s\n", conn.RemoteAddr())

	for {
		// Отправляем случайную поговорку клиенту
		randomGPV := getRandomGPV(gpv)
		_, err := conn.Write([]byte(randomGPV + "\n"))
		if err != nil {
			fmt.Println("Ошибка при отправке данных клиенту:", err)
			return
		}

		// Ждем 3 секунды перед отправкой следующей поговорки
		time.Sleep(3 * time.Second)
	}
}

func loadGPV(filename string) ([]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	gpv := strings.Split(string(data), "\n")
	return gpv, nil
}

func getRandomGPV(gpv []string) string {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(gpv))
	return gpv[randomIndex]
}
