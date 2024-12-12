package influxdb

import (
	"backend/internal/netflow"
	"context"
	"fmt"
	"log"
	"net"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

// InitEnv charge les variables d'environnement depuis un fichier .env
func InitEnv() {
	// Charger le fichier .env
	err := godotenv.Load("/home/ubuntu/monitoring-ng/backend/.env")
	if err != nil {
		log.Printf("Avertissement : impossible de charger le fichier .env, utilisation des variables d'environnement du système.")
	}
}

// Configuration des variables d'environnement
var (
	influxDBURL    string
	influxDBToken  string
	influxDBOrg    string
	influxDBBucket string
)

// InitConfig initialise les variables globales depuis l'environnement
func InitConfig() {
	influxDBURL = os.Getenv("INFLUXDB_URL")
	influxDBToken = os.Getenv("INFLUXDB_API_KEY")
	influxDBOrg = os.Getenv("INFLUXDB_ORG")
	influxDBBucket = os.Getenv("INFLUXDB_BUCKET")

	// Validation des variables critiques
	if influxDBURL == "" || influxDBToken == "" || influxDBOrg == "" || influxDBBucket == "" {
		log.Fatalf("Erreur : une ou plusieurs variables d'environnement nécessaires sont manquantes.")
	}
}

// InitInfluxDB initialise le client InfluxDB
func InitInfluxDB() influxdb2.Client {
	return influxdb2.NewClient(influxDBURL, influxDBToken)
}

// WriteToInfluxDB écrit un enregistrement NetFlow dans InfluxDB
func WriteToInfluxDB(client influxdb2.Client, header netflow.NetFlowV5Header, record netflow.NetFlowV5Record) {
	writeAPI := client.WriteAPIBlocking(influxDBOrg, influxDBBucket)

	startTime := netflow.ConvertNetFlowTime(header, record.StartTime)
	endTime := netflow.ConvertNetFlowTime(header, record.EndTime)

	p := influxdb2.NewPointWithMeasurement("netflow").
		AddTag("src_ip", net.IP(record.SrcAddr[:]).String()).
		AddTag("dst_ip", net.IP(record.DstAddr[:]).String()).
		AddField("src_port", record.SrcPort).
		AddField("dst_port", record.DstPort).
		AddField("protocol", record.Proto).
		AddField("bytes", record.Bytes).
		AddField("packets", record.Packets).
		AddField("start_time", startTime.UnixNano()).
		AddField("end_time", endTime.UnixNano()).
		SetTime(startTime)

	if err := writeAPI.WritePoint(context.Background(), p); err != nil {
		fmt.Printf("Erreur lors de l'écriture dans InfluxDB : %v\n", err)
	}
}

// Query NetFlow dans InfluxDB
func QueryNetFlowData() ([]map[string]interface{}, error) {
	client := InitInfluxDB()
	defer client.Close()

	queryAPI := client.QueryAPI(influxDBOrg)
	query := `from(bucket: "ntc-bucket")
		|> range(start: -1h)
		|> filter(fn: (r) => r._measurement == "netflow")
		|> filter(fn: (r) => r.src_ip == "192.168.108.115" or r.dst_ip == "192.168.108.115")
		|> pivot(rowKey:["_time"], columnKey:["_field"], valueColumn: "_value")
		|> keep(columns: ["src_ip", "dst_ip", "src_port", "dst_port", "protocol", "bytes", "packets"])
		|> group()`
	fmt.Printf("Exécution de la requête : %s\n", query)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("Erreur lors de la requête : %v\n", err)
		return nil, err
	}

	var flows []map[string]interface{}
	for result.Next() {
		// fmt.Printf("Enregistrement trouvé : %v\n", result.Record().Values())
		flows = append(flows, result.Record().Values())
	}

	if result.Err() != nil {
		fmt.Printf("Erreur dans les résultats : %v\n", result.Err())
	}

	return flows, nil
}
