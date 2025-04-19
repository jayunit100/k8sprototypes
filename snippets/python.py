import subprocess
import re
from collections import deque, Counter
from functools import reduce

class Pet:
    def __init__(self, name: str, species: str):
        self.name = name
        self.species = species

    def __str__(self):
        return f"{self.species} named {self.name}"

class PetStore:
    def __init__(self):
        self.pets = deque()  # Using deque for efficient appends and pops
    def total_pets(self) -> int:
        return len(self.pets)

    def average_pet_name_length(self) -> float:
        if not self.pets:
            return 0.0
        total_length = sum(map(lambda p: len(p.name), self.pets))
        return total_length / len(self.pets)

    def species_percentages(self) -> dict:
        counts = self.count_pet_species()
        total = len(self.pets)
        return {species: (count / total) * 100 for species, count in counts.items()}

    def add_pet(self, name: str, species: str) -> None:
        pet = Pet(name, species)
        self.pets.append(pet)
        print(f"Added a new {species} named {name}.")

    def scan_and_add_pets(self, directory: str) -> None:
        try:
            #### This shows an example of how to use python subprocess.  It would be more effective to use os.list of course.
            # Choose command based on OS
            command = ["ls", directory] if os.name == 'posix' else ["dir", directory, "/B", "/A:-D"]
            # Execute the command and decode the output
            files = subprocess.check_output(command, shell=(os.name != 'posix')).decode('utf-8')

            # Using map and lambda to filter and extract names
            pet_names = map(lambda f: re.match(r'pet-name-(.+)\.yaml', f), files.splitlines())
            pet_names = filter(None, pet_names)  # Remove None entries
            for match in pet_names:
                self.add_pet(match.group(1), "Dog")
        except subprocess.CalledProcessError:
            print(f"Failed to list directory {directory}")

    def list_pets(self) -> None:
        for pet in self.pets:
            print(pet)

    def count_pet_species(self) -> Counter:
        # Using map to extract species and Counter to count occurrences
        species = map(lambda pet: pet.species, self.pets)
        return Counter(species)

    def oldest_pet(self) -> str:
        # Reduce function to find the oldest pet added (first in the deque)
        return reduce(lambda a, b: a if self.pets.index(a) < self.pets.index(b) else b, self.pets).__str__()

# Usage
store = PetStore()
store.scan_and_add_pets("/path/to/pet/files")  # Replace with actual path
store.list_pets()

# Print the count of each species
print("Pet species count:", store.count_pet_species())

# Print the oldest pet
print("Oldest pet in the store:", store.oldest_pet())
print("Total pets:", store.total_pets())
print("Average pet name length:", store.average_pet_name_length())
print("Species percentages:", store.species_percentages())

