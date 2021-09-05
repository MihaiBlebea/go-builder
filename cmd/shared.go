package cmd

import "fmt"

// Used by the mute flag in commands
var mute bool

func askConfirm() (bool, error) {
	fmt.Println("Please confirm [yes/no]:")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false, err
	}

	okResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokResponses := []string{"n", "N", "no", "No", "NO"}

	if contains(okResponses, response) {
		return true, nil
	} else if contains(nokResponses, response) {
		return false, nil
	} else {
		return false, nil
	}
}

func contains(list []string, needle string) bool {
	for _, el := range list {
		if el == needle {
			return true
		}
	}

	return false
}
