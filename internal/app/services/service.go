package services

import (
	"log"
	"time"

	"github.com/lastreq/gas-price-test-task/internal/app/model"
)

type Service struct {
	prv provider
}

type provider interface {
	GetGasHistory() (model.Data, error)
}

func New(prv provider) *Service {
	svc := Service{prv: prv}
	return &svc
}

func (svc *Service) GetGasInfo() (model.ProcessedGasInfo, error) {
	data, err := svc.prv.GetGasHistory()
	if err != nil {
		return model.ProcessedGasInfo{}, err
	}

	processedGasInfo, err := calculateGasInfo(data.Ethereum.Transactions)
	if err != nil {
		return model.ProcessedGasInfo{}, err
	}

	return processedGasInfo, nil
}

func calculateGasInfo(transactions []model.Transactions) (processedGasInfo model.ProcessedGasInfo, err error) {
	monthsGasValueChan := make(chan map[string]float64)
	dayAverageGasPriceChan := make(chan float64)
	hoursAverageGasPriceChan := make(chan map[int]float64)
	allTimePaidChan := make(chan float64)

	go calculateMonthsGasValue(monthsGasValueChan, transactions)

	go calculateDayAverageGasPrice(dayAverageGasPriceChan, transactions)

	go calculateHoursAverageGasPrice(hoursAverageGasPriceChan, transactions)

	go calculateAllTimePaid(allTimePaidChan, transactions)

	processedGasInfo.MonthsGasValue = <-monthsGasValueChan
	processedGasInfo.DayAverageGasPrice = <-dayAverageGasPriceChan
	processedGasInfo.HoursAverageGasPrice = <-hoursAverageGasPriceChan
	processedGasInfo.AllTimePaid = <-allTimePaidChan

	return processedGasInfo, err
}

func calculateMonthsGasValue(monthsGasValueChan chan map[string]float64, transactions []model.Transactions) {
	monthsGasValue := make(map[string]float64)

	for i := range transactions {
		date, err := time.Parse("06-01-02 15:04", transactions[i].Time)
		if err != nil {
			log.Fatal(err)
		}

		month := date.Month().String()
		monthsGasValue[month] += transactions[i].GasValue
	}
	monthsGasValueChan <- monthsGasValue
}

func calculateDayAverageGasPrice(dayAverageGasPriceChan chan float64, transactions []model.Transactions) {
	var medianGasPriceSum, hoursSum float64

	for i := range transactions {
		hoursSum++

		medianGasPriceSum += transactions[i].MedianGasPrice
	}
	dayAverageGasPriceChan <- medianGasPriceSum / hoursSum
}

func calculateHoursAverageGasPrice(hoursAverageGasPriceChan chan map[int]float64, transactions []model.Transactions) {
	hoursAverageGasPrice := make(map[int]float64)
	hourMedianGasPrice := make(map[int][]float64)

	for i := range transactions {
		date, err := time.Parse("06-01-02 15:04", transactions[i].Time)
		if err != nil {
			log.Fatal(err)
		}

		hour := date.Hour()
		hourMedianGasPrice[hour] = append(hourMedianGasPrice[hour], transactions[i].MedianGasPrice)

		hoursAverageGasPrice[hour] = transactions[i].MedianGasPrice
	}

	hoursAverageGasPrice, err := calculateHourMedianGasPrice(hourMedianGasPrice)
	if err != nil {
		log.Fatal(err)
	}

	hoursAverageGasPriceChan <- hoursAverageGasPrice
}

func calculateAllTimePaid(allTimePaidChan chan float64, transactions []model.Transactions) {
	var allTimePaid float64
	for i := range transactions {
		allTimePaid += transactions[i].GasPrice * transactions[i].GasValue
	}
	allTimePaidChan <- allTimePaid
}

func calculateHourMedianGasPrice(hourMedianGasPrice map[int][]float64) (hoursAverageGasPrice map[int]float64, err error) {
	hoursAverageGasPrice = make(map[int]float64)

	for i := range hourMedianGasPrice {
		var hourMedianGasPriceSum, hoursSum float64
		for j := range hourMedianGasPrice[i] {
			hourMedianGasPriceSum += hourMedianGasPrice[i][j]
			hoursSum++
		}

		hoursAverageGasPrice[i] = hourMedianGasPriceSum / hoursSum
	}

	return hoursAverageGasPrice, err
}
