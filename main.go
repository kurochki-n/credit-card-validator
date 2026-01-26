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

func extractBIN(cardNumber string) (int, error) {
	bin, err := strconv.Atoi(cardNumber[:6])
	if err != nil {
		fmt.Println("Ошибка извлечения BIN:", err)
		return -1, err
	}
	return bin, nil
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
		fmt.Println("Ошибка открытия:", err)
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
			fmt.Println(i, "Ошибка чтения BinFrom:", err)
			continue
		}
		binTo, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Println(i, "Ошибка чтения BinTo:", err)
			continue
		}
		bank := Bank{Name: string(parts[0]), BinFrom: binFrom, BinTo: binTo}
		banks = append(banks, bank)
		i++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения:", err)
		return nil, err
	}

	return banks, nil
}

func validateLuhn(cardNumber string) bool {
	sum := 0
	for i := len(cardNumber) - 1; i >= 0; i-- {
		intN := int(cardNumber[i])
		if intN%2 == 1 {
			intN *= 2
		}
		if intN > 9 {
			intN -= 9
		}
		sum += intN
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
		fmt.Println("Ошибка ввода:", err)
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
		bin, err := extractBIN(cardNumber)
		if err != nil {
			continue
		}
		bank := identifyBank(bin, banks)
		fmt.Println("Банк:", bank)
	}
}
