package utils

// usage: i, err := stringToInt(s); if err != nil ...
func StringToInt(s string) (int, error) {
  i, err := strconv.ParseInt(s, 10, 0)
  if err != nil {
    return 0, err
  }
  return int(i), nil
}

func IntToString(i int) string {
	s := strconv.Itoa(i)
	return s
}

