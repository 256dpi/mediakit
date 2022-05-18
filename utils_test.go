package mediakit

import "os"

func loadSample(name string) *os.File {
	f, err := os.Open("./samples/" + name)
	if err != nil {
		panic(err)
	}

	return f
}
