package tokeniser

var keywords = []string{"use", "import", "when"}

func tryKeyword(id string) *Token {
	var (
		token      Token
		tokenFound = false
	)

	for _, keyword := range keywords {
		if id == keyword {
			token = Token{
				Type: keyword,
			}

			tokenFound = true
			break
		}
	}

	if tokenFound {
		return &token
	} else {
		return nil
	}
}
