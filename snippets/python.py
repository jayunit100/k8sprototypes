import os
import re
from collections import Counter

class Pet:
    def __init__(self, name: str, species: str):
        self.name = name
        self.species = species

    def __str__(self):
        return f"{self.species} named {self.name}"

class PetStore:
    def __init__(self):
        self.pets = []

    def add_pet(self, name: str, species: str) -> None:
        pet = Pet(name, species)
        self.pets.append(pet)
        print(f"Added a new {species} named {name}.")

    def scan_and_add_pets(self, directory: str) -> None:
        try:
            # List all files in the given directory
            for filename in os.listdir(directory):
                # Using regular expression to match 'pet-name-<name>.yaml'
                match = re.match(r'pet-name-(.+)\.yaml', filename)
                if match:
                    pet_name = match.group(1)
                    # Assuming a fixed species for simplicity. Modify as needed.
                    self.add_pet(pet_name, "Dog")
        except FileNotFoundError:
            print(f"Directory {directory} not found.")

    def list_pets(self) -> None:
        for pet in self.pets:
            print(pet)

# Usage
store = PetStore()
store.scan_and_add_pets("/path/to/pet/files")  # Replace with actual path
store.list_pets()
