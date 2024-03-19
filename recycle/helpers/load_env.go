package helpers

/*
func LoadEnv(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			fmt.Printf("Invalid line in .env file: %s\n", line)
			continue
		}
		key, value := parts[0], parts[1]
		os.Setenv(key, value)
	}

	return scanner.Err()
}
*/
