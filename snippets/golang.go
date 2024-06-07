package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Pet struct {
	Name    string
	Species string
}

func (p Pet) String() string {
	return fmt.Sprintf("%s named %s", p.Species, p.Name)
}

type PetStore struct {
	Pets []Pet
}

func (store *PetStore) AddPet(name, species string) {
	pet := Pet{Name: name, Species: species}
	store.Pets = append(store.Pets, pet)
	fmt.Printf("Added a new %s named %s.\n", species, name)
}

func (store *PetStore) ScanAndAddPets(directory string) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Printf("Failed to list directory %s\n", directory)
		return
	}

	regex := regexp.MustCompile(`pet-name-(.+)\.yaml`)
	for _, file := range files {
		if match := regex.FindStringSubmatch(file.Name()); match != nil {
			store.AddPet(match[1], "Dog")
		}
	}
}

func (store *PetStore) ListPets() {
	for _, pet := range store.Pets {
		fmt.Println(pet)
	}
}

func main() {
	store := PetStore{}
	store.ScanAndAddPets("/path/to/pet/files") // Replace with actual path
	store.ListPets()
}
