package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

// NetFlowV5Header représente l'entête du format NetFlow v5
type NetFlowV5Header struct {
	Version      uint16
	Count        uint16
	SysUptime    uint32
	UnixSecs     uint32
	UnixNsecs    uint32
	FlowSequence uint32
	EngineType   uint8
	EngineID     uint8
	Sampling     uint16
}

// NetFlowV5Record représente un enregistrement unique de NetFlow v5
type NetFlowV5Record struct {
	SrcAddr   [4]byte
	DstAddr   [4]byte
	NextHop   [4]byte
	Input     uint16
	Output    uint16
	Packets   uint32
	Bytes     uint32
	StartTime uint32
	EndTime   uint32
	SrcPort   uint16
	DstPort   uint16
	Pad1      uint8
	TCPFlags  uint8
	Proto     uint8
	Tos       uint8
	SrcAS     uint16
	DstAS     uint16
	SrcMask   uint8
	DstMask   uint8
	Pad2      uint16
}

// Mapping pour identifier les applications courantes en fonction des ports
var portToApplication = map[uint16]string{
	80:   "HTTP",
	443:  "HTTPS",
	22:   "SSH",
	53:   "DNS",
	9790: "NI",
}

// Fonction pour obtenir le nom de l'application en fonction du port
func getApplicationName(port uint16) string {
	if app, exists := portToApplication[port]; exists {
		return app
	}
	return "Unknown"
}

// Fonction utilitaire pour définir une couleur différente par ligne
func getColor(index int) string {
	colors := []string{
		"\033[31m", // Rouge
		"\033[32m", // Vert
		"\033[33m", // Jaune
		"\033[34m", // Bleu
		"\033[35m", // Magenta
		"\033[36m", // Cyan
	}
	return colors[index%len(colors)]
}

// Fonction pour afficher un flux NetFlow avec des couleurs
func printFlow(i int, record NetFlowV5Record, srcApp, dstApp string) {
	color := getColor(i) // Couleur basée sur l'index
	reset := "\033[0m"   // Réinitialisation des couleurs

	// Affichage dans une ligne avec les colonnes alignées
	fmt.Printf("%sFlow #%d: SrcIP=%v, DstIP=%v, SrcPort=%d(%s), DstPort=%d(%s), Protocol=%d, Bytes=%d, Packets=%d%s\n",
		color,                     // Couleur de début
		i+1,                       // Numéro du flux
		net.IP(record.SrcAddr[:]), // Adresse IP source
		net.IP(record.DstAddr[:]), // Adresse IP destination
		record.SrcPort, srcApp,    // Port source + Application
		record.DstPort, dstApp, // Port destination + Application
		record.Proto,   // Protocole
		record.Bytes,   // Nombre d'octets
		record.Packets, // Nombre de paquets
		reset,          // Réinitialisation de la couleur
	)
}

func main() {
	address := ":2055" // Écoute sur le port UDP 2055 (port NetFlow par défaut)

	// Crée un écouteur UDP
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		fmt.Printf("Erreur lors de la création de l'écoute UDP : %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Écoute des données NetFlow sur %s...\n", address)

	buf := make([]byte, 1500) // Alloue un tampon pour les paquets entrants

	for {
		// Lit les données de la connexion
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Printf("Erreur lors de la lecture de la connexion : %v\n", err)
			continue
		}

		fmt.Printf("Reçu %d octets de %s\n", n, addr)

		// Analyse l'entête NetFlow
		header := NetFlowV5Header{}
		reader := bytes.NewReader(buf[:n])
		if err := binary.Read(reader, binary.BigEndian, &header); err != nil {
			fmt.Printf("Erreur lors de l'analyse de l'entête : %v\n", err)
			continue
		}

		fmt.Printf("Version NetFlow : %d, Nombre d'enregistrements : %d\n", header.Version, header.Count)

		// Analyse les enregistrements NetFlow
		for i := 0; i < int(header.Count); i++ {
			record := NetFlowV5Record{}
			if err := binary.Read(reader, binary.BigEndian, &record); err != nil {
				fmt.Printf("Erreur lors de l'analyse de l'enregistrement : %v\n", err)
				break
			}

			// Obtenez les noms des applications
			srcApp := getApplicationName(record.SrcPort)
			dstApp := getApplicationName(record.DstPort)

			// Affichez le flux avec la fonction améliorée
			printFlow(i, record, srcApp, dstApp)
		}
	}
}
