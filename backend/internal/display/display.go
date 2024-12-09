package display

import (
	"fmt"
	"net"
	"os"

	"backend/internal/netflow"

	"github.com/olekukonko/tablewriter"
)

// DisplayFlowTable affiche un tableau des flux NetFlow dans le terminal
func DisplayFlowTable(records []netflow.NetFlowV5Record, header netflow.NetFlowV5Header, filterIPs []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "Type", "SrcIP", "DstIP", "SrcPort", "DstPort", "Protocol", "Bytes", "Packets"})

	colors := []tablewriter.Colors{
		{tablewriter.FgRedColor},
		{tablewriter.FgGreenColor},
		{tablewriter.FgYellowColor},
		{tablewriter.FgBlueColor},
		{tablewriter.FgCyanColor},
		{tablewriter.FgMagentaColor},
		{tablewriter.FgWhiteColor},
		{tablewriter.FgGreenColor},
	}

	// Traquer les paires Requête-Réponse
	seen := make(map[string]bool)

	for i, record := range records {
		srcIP := net.IP(record.SrcAddr[:]).String()
		dstIP := net.IP(record.DstAddr[:]).String()

		// Identifiant unique pour une paire Requête-Réponse
		flowID := fmt.Sprintf("%s:%d->%s:%d", srcIP, record.SrcPort, dstIP, record.DstPort)
		reverseFlowID := fmt.Sprintf("%s:%d->%s:%d", dstIP, record.DstPort, srcIP, record.SrcPort)

		// Vérifier si le flux est déjà traité
		if seen[flowID] || seen[reverseFlowID] {
			continue
		}
		seen[flowID] = true
		seen[reverseFlowID] = true

		// Trouver la réponse correspondante
		var response *netflow.NetFlowV5Record
		for _, potentialResponse := range records {
			if net.IP(potentialResponse.SrcAddr[:]).String() == dstIP &&
				net.IP(potentialResponse.DstAddr[:]).String() == srcIP &&
				potentialResponse.SrcPort == record.DstPort &&
				potentialResponse.DstPort == record.SrcPort {
				response = &potentialResponse
				break
			}
		}

		// Ajouter la requête au tableau
		addFlowToTable(table, i+1, "Requête", record, colors)

		// Ajouter la réponse au tableau si elle existe
		if response != nil {
			addFlowToTable(table, i+1, "Réponse", *response, colors)
		}
	}

	table.SetBorder(true)
	table.SetRowLine(true)
	table.Render()
}

// Ajoute une ligne au tableau
func addFlowToTable(table *tablewriter.Table, index int, flowType string, record netflow.NetFlowV5Record, colors []tablewriter.Colors) {
	srcIP := net.IP(record.SrcAddr[:]).String()
	dstIP := net.IP(record.DstAddr[:]).String()

	lineData := []string{
		fmt.Sprintf("%d", index),
		flowType,
		srcIP,
		dstIP,
		fmt.Sprintf("%d", record.SrcPort),
		fmt.Sprintf("%d", record.DstPort),
		fmt.Sprintf("%d", record.Proto),
		fmt.Sprintf("%d", record.Bytes),
		fmt.Sprintf("%d", record.Packets),
	}

	var rowColors []tablewriter.Colors
	for colIdx := range lineData {
		rowColors = append(rowColors, colors[colIdx%len(colors)])
	}

	table.Rich(lineData, rowColors)
}
