package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/k0kubun/pp"
)

// 1. Создаем структуру данных Bank
// 2. Создаем базу данных банков в txt
// 3. Создаем func для чтения и парсинга данных из файла banks.txt
// 4. Создаем func, которая определяет банк по первым цифрам номера
// 5. Создаем func алгоритма Луна
// 6. Создаем func для получения номера карты от пользователя
// 7. Создаем func для проверки введенного номера карты
// 8. Собираем воедино

type Bank struct {
	// Название банка
	Name string
	// Начало диапазона номеров
	BinFrom int
	// Конец диапазона номеров
	BinTo int
}

func loadBankData(path string) ([]Bank, error) {
	// Открываем файл и обрабатывает ошибки
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return nil, err
	}
	defer file.Close()

	// Инициализируем пустой слайс структуры
	banks := make([]Bank, 0)

	// Читаем файл построчно
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// парсинг строки файла
		str := scanner.Text()
		parts := strings.SplitN(str, ",", 3)
		if len(parts) != 3 {
			fmt.Println("Неправильный формат строки", str)
			continue
		}

		// Преобразование типов
		name := parts[0]

		binFrom, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Не получилось преобразовать binFrom в int")
			continue
		}

		binTo, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Println("Не получилось преобразовать binTo в int")
			continue
		}

		// заполняем структуру
		bank := Bank{
			Name:    name,
			BinFrom: binFrom,
			BinTo:   binTo,
		}

		// добавляем в слайс
		banks = append(banks, bank)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения файла", err)
		return nil, err
	}

	return banks, nil
}

func extractBIN(cardNumber string) int {
	// 1. Принимает номер карты
	// 2. Возвращает первые 6 цифр

	// Валидируем
	if len(cardNumber) < 6 {
		fmt.Println("Некорректная длина номера")
		return 0
	}

	// Из строки переводим в int
	binStr := cardNumber[:6]

	bin, err := strconv.Atoi(binStr)
	if err != nil {
		fmt.Println("Не удалось конвертировать BIN")
		return 0
	}

	return bin
}

func identifyBank(bin int, banks []Bank) string {
	// 1. Принимает bin и слайс банков
	// 2. Возвращает название банка

	for _, bank := range banks {
		if bin >= bank.BinFrom && bin <= bank.BinTo {
			return bank.Name
		}
	}

	return "Неизвестный банк"
}

func validateLuhn(cardNumber string) bool {
	sum := 0
	// флаг для удвоения цифры на чётных позициях с конца строки
	isEven := false

	for i := len(cardNumber) - 1; i >= 0; i-- {

		digit := int(cardNumber[i] - '0')
		// если не цифра возвращаем false
		if digit < 0 || digit > 9 {
			return false
		}

		if isEven {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isEven = !isEven
	}

	return sum%10 == 0
}

func getUserInput() string {
	r := bufio.NewReader(os.Stdin)

	pp.Println("Введите данные карты:")

	str, err := r.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении:", err)
		return ""
	}

	pp.Println("Ваши данные:")
	fmt.Println(strings.TrimSpace(str))
	fmt.Println("")

	return strings.TrimSpace(str)
}

func validateInput(cardNumber string) bool {

	// Убираем пробелы в начале и в конце
	cardNumber = strings.TrimSpace(cardNumber)

	// проверяем длину строки (13-19 символов)
	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		fmt.Println("Неправильное количество цифр номера")
		return false
	}

	for i := 0; i < len(cardNumber); i++ {
		// преобразуем в int
		digit := int(cardNumber[i] - '0')
		// проверяем, чтобы все данные были цифрами
		if digit < 0 || digit > 9 {
			fmt.Println("Данные не являются цифрами")
			return false
		}
	}

	return true
}

func main() {
	fmt.Println("Добро пожаловать в программу валидации карт!")

	banks, err := loadBankData("banks.txt")
	if err != nil {
		fmt.Println("Ошибка", err)
		return
	}

	//pp.Println(validateLuhn("4532015112830366"))
	//pp.Println(validateLuhn("1234567890123456"))
	//fmt.Println("")
	//pp.Println(identifyBank(800456, banks))

	for {
		cardNumber := getUserInput()
		//fmt.Println(validateInput(cardNumber))

		if len(cardNumber) == 0 {
			fmt.Println("Пустая строка")
			pp.Println("Программа завершена")
			break
		}

		if validateInput(cardNumber) != true {
			fmt.Println("Ошибка формата")
			continue
		}

		if validateLuhn(cardNumber) != true {
			fmt.Println("Невалидный номер")
			continue
		}

		bin := extractBIN(cardNumber)
		fmt.Println(identifyBank(bin, banks))
	}

}
