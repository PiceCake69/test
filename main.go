package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const apiURL = "https://www.cbr.ru/scripts/XML_daily.asp"

type ValCurs struct {
	Date      string `xml:"Date,attr"`
	ValuteArr []Valute
}

type Valute struct {
	NumCode    string `xml:"NumCode"`
	CharCode   string `xml:"CharCode"`
	Nominal    int    `xml:"Nominal"`
	Name       string `xml:"Name"`
	Value      string `xml:"Value"`
}

func getExchangeRates(date string) (ValCurs, error) {
	response, err := http.Get(fmt.Sprintf("%s?date_req=%s", apiURL, date))
	if err != nil {
		return ValCurs{}, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ValCurs{}, err
	}

	var valCurs ValCurs
	err = xml.Unmarshal(body, &valCurs)
	if err != nil {
		return ValCurs{}, err
	}

	return valCurs, nil
}

func main() {
	// Чтение даты с консоли
	fmt.Print("Введите дату (в формате ДД.ММ.ГГГГ): ")
	var inputDate string
	fmt.Scanln(&inputDate)

	// Преобразование в формат, поддерживаемый API
	date, err := time.Parse("02.01.2006", inputDate)
	if err != nil {
		fmt.Println("Ошибка: неверный формат даты.")
		os.Exit(1)
	}
	apiDate := date.Format("02/01/2006")

	// Получение курсов валют
	valCurs, err := getExchangeRates(apiDate)
	if err != nil {
		fmt.Printf("Ошибка при получении курсов валют: %s\n", err.Error())
		os.Exit(1)
	}

	// Вывод результатов
	fmt.Printf("Курсы доллара ЦБ РФ на %s:\n", valCurs.Date)
	for _, valute := range valCurs.ValuteArr {
		if valute.CharCode == "USD" {
			fmt.Printf("USD (Доллар США) = %s рублей\n", valute.Value)
		}
	}
}