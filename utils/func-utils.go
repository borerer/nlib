package utils

func Must(err error) {
	if err != nil {
		println(err.Error())
		panic(err)
	}
}
