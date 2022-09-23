package p

type (
	test struct {
		Num int
	}
)

func alternatingParams(a string, b bool, c string) {}

func combinedParams1(a, b string, c bool) {}

func combinedParams2(a bool, b, c string) {}

func combinedParams3(a, b string, c bool, d, e string) {}

func combinedParams4(a bool, b, c string, d bool, e string) {}

func returnAString(a, b string) string { return "" }

func returnAString2(a bool) string { return "" }

func (a *test) withStruct(b bool, c, d string, e bool, f string) (string, int) {
	return "", 0
}
