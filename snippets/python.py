
a = {
    "food":{"popcorn","soda"},
    "beer":{"corona","dos equis"}
}

print(a)

##### Sets 

x = "ABCDE"
for i in range(1,5):
    print(i)
     
x = set()
x.clear()
x.add("a")

print(max(10,1))

##### Maps and numbers 

x = 0
x +=1

m = {}

if "a" in m:
    print("not em")
    
from collections import Counter
l = ['rose','tulips','sunflowers','tulips','rose']
my_count = Counter(l)
print(my_count, "list")

### Sorting a list

l = [99999,1,2,3,4,0,10]
l.sort()
print(l)


l.extend(["a", "b"])
print(l)


###### Classes ######

from typing import List

class Car:
    # drive returns thea list of stops a car made while driving 'm' miles....
    def drive(miles int) -> List[int]:
        return 4

print("a")
x = Car()




