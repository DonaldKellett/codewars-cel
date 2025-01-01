package cel

import "io/ioutil"

func ReadUserSolution() ([]byte, error) {
	return ioutil.ReadFile("solution.txt")
}
