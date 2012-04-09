package main

var uniqs = make(chan uint64)

func init() {
	i := uint64(0)
	for {
		uniqs <- i
		i++
	}
}

func uniq() uint64 {
	return <-uniqs
}
