package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Bank struct {
	Name    string
	BinFrom int
	BinTo   int
}

func extractBIN(cardNumber string) int {
	bin, _ := strconv.Atoi(cardNumber[:6])
	return bin
}

func identifyBank(bin int, banks []Bank) string {
	for _, bank := range banks {
		if bin >= bank.BinFrom && bin <= bank.BinTo {
			return bank.Name
		}
	}
	return "Неизвестный банк"
}

func loadBankData(path string) ([]Bank, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	banks := []Bank{}
	i := 0
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		binFrom, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		binTo, err := strconv.Atoi(parts[2])
		if err != nil {
			continue
		}
		bank := Bank{Name: string(parts[0]), BinFrom: binFrom, BinTo: binTo}
		banks = append(banks, bank)
		i++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return banks, nil
}

func validateLuhn(cardNumber string) bool {
	parity := len(cardNumber) % 2
	sum := 0

	for i := 0; i < len(cardNumber); i++ {
		num := int(cardNumber[i] - '0')
		if (i)%2 == parity {
			num *= 2
			if num > 9 {
				num -= 9
			}
		}
		sum += num
	}
	
	return sum%10 == 0
}

func validateInput(input string) bool {
	if len(input) < 13 || len(input) > 19 {
		return false
	}
	for _, r := range input {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func getUserInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите номер карты> ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	cleanInput := strings.TrimSpace(input)
	return cleanInput, nil
}

func main() {
	fmt.Println("Credit card validator")
	banks, err := loadBankData("banks.txt")
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return
	}

	for {
		cardNumber, err := getUserInput()
		if err != nil {
			fmt.Println("Номер карты недействителен")
			continue
		}
		if len(cardNumber) == 0 {
			fmt.Println("Завершение программы")
			return
		}
		if !validateInput(cardNumber) || !validateLuhn(cardNumber) {
			fmt.Println("Номер карты недействителен")
			continue
		}
		bin := extractBIN(cardNumber)
		bank := identifyBank(bin, banks)
		fmt.Println("Банк:", bank)
	}
}
