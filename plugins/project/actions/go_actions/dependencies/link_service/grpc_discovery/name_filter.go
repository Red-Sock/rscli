package grpc_discovery

func FilterPackageName(packageName string) string {
	packageNameB := []rune(packageName)
	out := make([]rune, 0, len(packageNameB))

	for idx := range packageNameB {
		if packageNameB[idx] >= 'A' && packageNameB[idx] <= 'Z' {
			packageNameB[idx] = packageNameB[idx] + 32
			out = append(out, '!')
		}

		out = append(out, packageNameB[idx])
	}

	return string(out)
}
