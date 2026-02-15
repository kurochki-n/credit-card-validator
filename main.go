package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Bank представляет банк с диапазоном BIN кодов.
type Bank struct {
	Name    string // Название банка
	BinFrom int    // Начало диапазона BIN
	BinTo   int    // Конец диапазона BIN
}

// extractBIN извлекает первые 6 цифр номера карты (BIN - Bank Identification Number).
// Возвращает BIN как целое число.
func extractBIN(cardNumber string) int {
	bin, _ := strconv.Atoi(cardNumber[:6])
	return bin
}

// identifyBank определяет банк-эмитент по BIN коду.
// Принимает BIN код и список банков, возвращает название банка или "Неизвестный банк".
func identifyBank(bin int, banks []Bank) string {
	for _, bank := range banks {
		if bin >= bank.BinFrom && bin <= bank.BinTo {
			return bank.Name
		}
	}
	return "Неизвестный банк"
}

// loadBankData загружает данные о банках из файла.
// Ожидает формат CSV: Название,BINОт,BINДо
// Возвращает список банков или ошибку при проблемах чтения файла.
func loadBankData(path string) ([]Bank, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	banks := []Bank{}
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
		bank := Bank{Name: parts[0], BinFrom: binFrom, BinTo: binTo}
		banks = append(banks, bank)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return banks, nil
}

// validateLuhn проверяет номер карты по алгоритму Луна (Luhn algorithm).
// Возвращает true, если номер карты прошел валидацию контрольной суммы.
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

// validateInput проверяет формат входной строки номера карты.
// Номер должен содержать от 13 до 19 цифр согласно стандарту ISO/IEC 7812.
// Возвращает true, если формат корректен.
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

// getUserInput читает номер карты из стандартного ввода.
// Возвращает очищенную от пробельных символов строку или ошибку чтения.
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
	// Определение флагов командной строки
	banksFile := flag.String("banks", "banks.csv", "Путь к файлу с данными о банках")
	flag.Parse()

	fmt.Println("Credit card validator")
	banks, err := loadBankData(*banksFile)
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
