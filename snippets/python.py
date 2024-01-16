##### Example of using Class and nonlocal variables...
Class Solution:
    ### CLASS  Variables, accessed via Solution.Vowels
    VOWELS = ("a","e","i","o","u")
    def maxVowels(self, s: str, k: int) -> int:
        curr_max = 0

        def update_talley(num_vowels: int):
            ### LOCAL variables, accessed via nonlocal
            nonlocal curr_max
            curr_max = max(curr_max, num_vowels)

##### Pet Store Simulator
from collections import Counter
from typing import List
import math

##### Map of Maps for pet categories and their items
inventory = {
    "food": {"dog food", "cat food"},
    "toys": {"ball", "mouse toy"}
}

##### Classes for Pets and Store
class Pet:
    def __init__(self, name: str, species: str):
        self.name = name
        self.species = species

    def __str__(self):
        return f"{self.species} named {self.name}"

class PetStore:
    def __init__(self):
        self.pets = []
        self.sales = Counter()

    # Function Definitions with arguments and return type
    def add_pet(self, pet: Pet) -> None:
        self.pets.append(pet)

    # Lambda Function for sorting pets by name
    def list_pets(self) -> List[str]:
        return sorted(self.pets, key=lambda pet: pet.name)

    # File Operations for sales record
    def record_sale(self, item: str) -> None:
        self.sales[item] += 1
        with open("sales.txt", "a") as file:
            file.write(f"Sold {item}\n")

    def read_sales(self) -> None:
        try:
            with open("sales.txt", "r") as file:
                print(file.read())
        except FileNotFoundError:
            print("Sales file not found.")

# Creating instance of PetStore
store = PetStore()

##### Adding Pets (Using Class)
store.add_pet(Pet("Buddy", "Dog"))
store.add_pet(Pet("Whiskers", "Cat"))

##### List Comprehensions for inventory check
available_food = [item for item in inventory["food"]]
print(f"Available food: {available_food}")

##### Sorting a List of pets
sorted_pets = store.list_pets()
print("Pets in store:", sorted_pets)

##### Adding More Elements to List and Map of Maps
inventory["accessories"] = {"leash", "collar"}

##### Exception Handling in File Operations
store.record_sale("dog food")
store.read_sales()

##### Modules and Import Statements - Using math
print(f"Square root of 16 (using math module): {math.sqrt(16)}")

##### Generator Expressions for pet names
pet_names = (pet.name for pet in store.pets)
for name in pet_names:
    print(f"Pet name: {name}")

##### Dictionary Comprehensions for pet species count
species_count = {pet.species: store.pets.count(pet) for pet in store.pets}
print("Pet species count:", species_count)
