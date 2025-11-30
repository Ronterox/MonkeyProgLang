let loop = fn(len, i, inc, iter) {
	if (i > len - 1) {
		return i
	}
	iter(i)
	return loop(len, i+inc, inc)
}

loop(0, 10, 1, fn() { 10 })
