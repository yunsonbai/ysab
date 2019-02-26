
import os

os.remove("./urls.txt")

f = open("./urls.txt", "a+")
for i in range(2):
    f.write("http://127.0.0.1:8080/add\n")
