package netutils

import (
	"net"
)

func Aggregate(a []net.IPNet) []net.IPNet {
	b := make([]net.IPNet, len(a))
	copy(b, a)
	SortIPNet(b)
	return aggregatePass(b)
}

func aggregatePass(p []net.IPNet) []net.IPNet {
	if len(p) <= 1 {
		return p
	}
	for i := 0; i < len(p)-1; i += 1 {
		c, merged := merge(p[i], p[i+1])
		if merged {
			q := append(p[:i], c)
			q = append(q, p[i+2:]...)
			return aggregatePass(q)
		}
	}
	return p
}

func merge(a, b net.IPNet) (c net.IPNet, merged bool) {
	a_ones, _ := a.Mask.Size()
	b_ones, _ := b.Mask.Size()
	if a_ones != b_ones {
		return
	}
	if a.IP.Equal(b.IP) {
		c = a
		merged = true
		return
	}
	c.IP = a.IP
	c.Mask = maskup(a.Mask)
	merged = c.Contains(b.IP)
	return
}

func maskup(mask net.IPMask) net.IPMask {
	m := make(net.IPMask, len(mask))
	copy(m, mask)
	for i := len(m) - 1; i >= 0; i -= 1 {
		if m[i] == 0x00 {
			continue
		}
		m[i] = m[i] ^ (^m[i] + 1)
		break
	}
	return m
}
