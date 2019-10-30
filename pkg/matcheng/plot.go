package matcheng

import (
	"bytes"
	"fmt"
)

func PlotCurves(bids []AggregatePrice, asks []AggregatePrice) string {
	var buf bytes.Buffer
	buf.WriteString("\"Ask\"\n")

	for i, entry := range asks {
		if i == 0 {
			buf.WriteString(fmt.Sprintf("%s 0\n", entry[0]))
		}
		if i > 0 {
			buf.WriteString(fmt.Sprintf("%s %s\n", entry[0], asks[i-1][1]))
		}
		buf.WriteString(fmt.Sprintf("%s %s\n", entry[0], entry[1]))
	}

	buf.WriteString("\n\n")
	buf.WriteString("\"Bid\"\n")

	for i := len(bids) - 1; i >= 0; i-- {
		entry := bids[i]
		if i == len(bids)-1 {
			buf.WriteString(fmt.Sprintf("%s 0\n", entry[0]))
		}
		if i != len(bids)-1 {
			buf.WriteString(fmt.Sprintf("%s %s\n", entry[0], bids[i+1][1]))
		}
		buf.WriteString(fmt.Sprintf("%s %s\n", entry[0], entry[1]))
		if i == 0 {
			buf.WriteString(fmt.Sprintf("0 %s\n", entry[1]))
		}
	}

	out := buf.Bytes()
	return string(out)
}
