let loop = fn(len, i, inc) {
	if (i > len - 1) {
		return i
	}
	return loop(len, i+inc, inc)
}

loop(0, 10, 1)

((a * ([1, 2, 3, 4][(b * c)])) * d)
((a * [1, 2, 3, 4][(b * c)]) * d)
