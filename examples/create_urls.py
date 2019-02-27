
import os

os.remove("./urls.txt")

with open("./urls.txt", "a+") as f:
    for i in range(2):
        # f.write("-u http://127.0.0.1:8080/add\n")
        f.write('-u http://127.0.0.1:8080/add -d \{"name":"yunson"\}\n')
