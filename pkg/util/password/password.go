package password

import (
	"math/rand"
	"time"
)

func Generate(length int, hasSpecial bool) string {
	rand.Seed(time.Now().UnixNano())

	upperChar := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowChar := "abcdefghijklmnopqrstuvwxyz"
	digits := "0123456789"
	all := upperChar + lowChar + digits
	buf := make([]byte, length)
	buf[0] = upperChar[rand.Intn(len(upperChar))]
	buf[1] = lowChar[rand.Intn(len(lowChar))]
	buf[2] = digits[rand.Intn(len(digits))]

	i := 3
	if hasSpecial {
		specials := "=+%/!@#$"
		all = all + specials
		buf[i] = specials[rand.Intn(len(specials))]
		i++
	}

	for ; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf)
}
