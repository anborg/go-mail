package main

import (
	"bufio"
	"fmt"
	"strings"
)

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func main1() {
	var x string = ` (I)8751                       PB20034 EST2 ASPHALT            31/07/20 $469974.60 001000006
		(I)8753                       PB20034 EST3 ASPHALT            15/08/20 $583537.82 001000006
		(I)8682                       PB20034 EST20101 ASPHALT        30/06/20$1353810.33 001000006
		(I)8685                       PB20034 JN20 ASPHALT INDEX      30/06/20 $125440.53 001000006
		`
	result := cleanInvoiceBlob(x)
	fmt.Println(result)
}

func cleanInvoiceBlob(multiLineBlob string) (cleanInvoice string) {
	scanner := bufio.NewScanner(strings.NewReader(multiLineBlob))
	cleanInvoice = ""
	for scanner.Scan() {
		var s = strings.TrimSpace(scanner.Text())
		cleanInvoice = cleanInvoice + substr(s, 30, 52) + "\n"
	}
	return
}

//(I)8751
//(I)8751                       PB20034 EST2 ASPHALT            31/07/20 $469974.60
//PB20034 EST2 ASPHALT            31/07/20 $469974.60
