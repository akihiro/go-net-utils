package netutils

import (
	"bytes"
	"net"
	"sort"
)

func SortIPNet(p []net.IPNet) {
	sort.Slice(p, func(i, j int) bool {
		return bytes.Compare(p[i].IP, p[j].IP) < 0
	})
}
