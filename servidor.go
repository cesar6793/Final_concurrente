package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"
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

// KMeans algoritmo de agrupación
func KMeans(data []DataPoint, k int) KMeansResult {
	// Implementación simplificada del algoritmo k-means
	centroids := make([]DataPoint, k)
	for i := 0; i < k; i++ {
		centroids[i] = data[rand.Intn(len(data))]
	}
	return KMeansResult{Centroids: centroids}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var data []DataPoint
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&data); err != nil {
		log.Println("Error decoding data:", err)
		return
	}

	result := KMeans(data, 4) // Asumiendo k = 4

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(result); err != nil {
		log.Println("Error encoding result:", err)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}
