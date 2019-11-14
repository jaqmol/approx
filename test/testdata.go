package test

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Person ...
type Person struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// Equals ...
func (a *Person) Equals(b *Person) bool {
	return a.ID == b.ID &&
		a.FirstName == b.FirstName &&
		a.LastName == b.LastName &&
		a.Email == b.Email
}

// MarshallPeople ...
func MarshallPeople(testData []Person) [][]byte {
	acc := make([][]byte, 0)
	for _, p := range testData {
		b, err := json.Marshal(&p)
		check(err)
		acc = append(acc, b)
	}
	return acc
}

func unMarshallPeople(data [][]byte) []Person {
	acc := make([]Person, 0)
	for i, b := range data {
		p, err := unmarshallPerson(b)
		if err != nil {
			log.Fatalf("Couldn't unmarshall person #%v: \"%v\"\n", i+1, string(b))
		}
		acc = append(acc, *p)
	}
	return acc
}

func unmarshallPerson(data []byte) (*Person, error) {
	var p Person
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func unmarshallError(data []byte) (*Person, error) {
	var p Person
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// MakePersonForIDMap ...
func MakePersonForIDMap(people []Person) map[string]Person {
	acc := make(map[string]Person)
	for _, p := range people {
		acc[p.ID] = p
	}
	return acc
}

// LoadTestData ...
func LoadTestData() []Person {
	dat, err := ioutil.ReadFile("../test/test-data.json")
	var result []Person
	err = json.Unmarshal(dat, &result)
	check(err)
	return result
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
