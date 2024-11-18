package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	// Parâmetros CLI
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 100, "Número total de requests")
	concurrency := flag.Int("concurrency", 10, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" {
		fmt.Println("Erro: A URL é obrigatória.")
		return
	}

	// Validações de entrada
	if *requests <= 0 || *concurrency <= 0 {
		fmt.Println("Erro: O número de requests e de concorrência deve ser maior que zero.")
		return
	}

	fmt.Printf("Iniciando teste de carga na URL: %s\n", *url)
	fmt.Printf("Total de Requests: %d | Concorrência: %d\n", *requests, *concurrency)

	startTime := time.Now()

	// Realizar requests
	statusCounts := make(map[int]int)
	var totalRequests int
	var mu sync.Mutex

	// Worker pool
	wg := &sync.WaitGroup{}
	sem := make(chan struct{}, *concurrency)

	for i := 0; i < *requests; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(requestNum int) {
			defer wg.Done()
			resp, err := http.Get(*url)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				statusCounts[0]++ // 0 para erros
				fmt.Printf("Request #%d: Erro de conexão\n", requestNum+1)
			} else {
				statusCounts[resp.StatusCode]++
				fmt.Printf("Request #%d: Status %d\n", requestNum+1, resp.StatusCode)
				resp.Body.Close()
			}
			totalRequests++

			<-sem
		}(i)
	}

	wg.Wait()
	elapsedTime := time.Since(startTime)

	// Geração do relatório
	fmt.Println("\n--- Relatório ---")
	fmt.Printf("Tempo total gasto: %.2fs\n", elapsedTime.Seconds())
	fmt.Printf("Total de requests realizados: %d\n", totalRequests)

	for status, count := range statusCounts {
		if status == 0 {
			fmt.Printf("Erros de conexão: %d\n", count)
		} else {
			fmt.Printf("Status %d: %d\n", status, count)
		}
	}
}
