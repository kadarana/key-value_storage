package benchpress

func SliceFunc(s []int) {
	var sum int

	for i, e := range s[:len(s)-1] {
		sum += e * s[i+1]
	}
}

func FindSum(list []int) int {
	sum := 0

	for _, number := range list {
		sum += number
	}

	return sum
}

func FindSumSlow(list []*int) int {
	sum := 0

	for _, number := range list {
		sum += *number
	}

	return sum
}
