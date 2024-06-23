package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

// DataPoint representa un punto de datos con latitud, longitud, tipo de delito y fecha/hora.
type DataPoint struct {
	Latitude  float64 `json:"latitud"`
	Longitude float64 `json:"longitud"`
	CrimeType string  `json:"tipo_delito"`
	DateTime  string  `json:"fecha_hora"`
}

// KMeansResult representa el resultado de la agrupación k-means.
type KMeansResult struct {
	Centroids []DataPoint `json:"centroids"`
}

func loadDataFromURL(url string) ([]DataPoint, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []DataPoint
	reader := csv.NewReader(resp.Body)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		lat, _ := strconv.ParseFloat(record[0], 64)
		lon, _ := strconv.ParseFloat(record[1], 64)
		data = append(data, DataPoint{
			Latitude:  lat,
			Longitude: lon,
			CrimeType: record[2],
			DateTime:  record[3],
		})
	}
	return data, nil
}

func sendDataToServer(data []DataPoint, address string, wg *sync.WaitGroup, results chan<- KMeansResult) {
	defer wg.Done()

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(data); err != nil {
		log.Println("Error encoding data:", err)
		return
	}

	var result KMeansResult
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&result); err != nil {
		log.Println("Error decoding result:", err)
		return
	}

	results <- result
}

func main() {
	url := "https://raw.githubusercontent.com/tu_usuario/tu_repositorio/main/datos_delitos.csv"
	data, err := loadDataFromURL(url)
	if err != nil {
		log.Fatal(err)
	}

	numServers := 4
	chunkSize := len(data) / numServers
	var wg sync.WaitGroup
	results := make(chan KMeansResult, numServers)

	for i := 0; i < numServers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numServers-1 {
			end = len(data)
		}
		wg.Add(1)
		go sendDataToServer(data[start:end], "localhost:8080", &wg, results)
	}

	wg.Wait()
	close(results)

	var allResults []KMeansResult
	for result := range results {
		allResults = append(allResults, result)
	}

	// Procesar y consolidar los resultados aquí
	log.Println("All results:", allResults)
}
