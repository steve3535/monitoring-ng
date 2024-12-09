package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"backend/internal/display"
	"backend/internal/influxdb"
	"backend/internal/netflow"
)

func main() {
	influxdb.InitEnv()
	influxdb.InitConfig()

	client := influxdb.InitInfluxDB()
	defer client.Close()

	address := ":2055" // Port UDP 2055
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		fmt.Printf("Erreur lors de l'écoute UDP : %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Écoute des données NetFlow v5 sur %s...\n", address)

	buf := make([]byte, 1500)                                // Tampon pour les paquets entrants
	filterIPs := []string{"192.168.108.115", "10.12.129.73"} // Adresses IP à filtrer

	for {
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Printf("Erreur lors de la réception de données : %v\n", err)
			continue
		}

		reader := bytes.NewReader(buf[:n])
		header := netflow.NetFlowV5Header{}
		if err := binary.Read(reader, binary.BigEndian, &header); err != nil {
			continue
		}

		var requests, responses []netflow.NetFlowV5Record
		for i := 0; i < int(header.Count); i++ {
			record := netflow.NetFlowV5Record{}
			if err := binary.Read(reader, binary.BigEndian, &record); err != nil {
				fmt.Printf("Erreur lors de l'analyse de l'enregistrement : %v\n", err)
				break
			}

			srcIP := net.IP(record.SrcAddr[:]).String()
			dstIP := net.IP(record.DstAddr[:]).String()

			// Séparer les flux en requêtes et réponses
			if contains(filterIPs, dstIP) {
				requests = append(requests, record) // Requête : DstIP correspond
			} else if contains(filterIPs, srcIP) {
				responses = append(responses, record) // Réponse : SrcIP correspond
			}
		}

		// Afficher les flux : Requêtes suivies des Réponses
		if len(requests) > 0 || len(responses) > 0 {
			fmt.Println("Nouveau tableau des flux filtrés :")
			display.DisplayFlowTable(append(requests, responses...), header, filterIPs)
		}
	}
}

// Vérifier si une adresse IP est dans la liste des adresses à filtrer
func contains(filterIPs []string, ip string) bool {
	for _, filterIP := range filterIPs {
		if filterIP == ip {
			return true
		}
	}
	return false
}
