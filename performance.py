import json
import random
import sys
import requests

available_words = ["word", "desktop", "random", "People", "TrUtH",
                   "Desk??", "MoNitor", "headphones", "glass", "performance",
                   "exErcise", "yay", "NAY"]
d = {}
sentence = ""
sentences = []

# build 5000 10 word sentences and keep track of their appearances
for i in range(5000):
    for k in range(10):
        idx = random.Random().randint(0, len(available_words)-1)
        word = available_words[idx]
        sentence += " " + word
        if word not in d:
            d[word] = 0
        d[word] += 1
    sentences.append(sentence)
    sentence = ""

for i in range(len(available_words)):
    available_words[i] = available_words[i].lower()
available_words[5] = "desk"

r = random.Random()

# send the 5000 sentences with GETs between POSTs
with requests.Session() as s:
    for i in range(5000):
        text = sentences[i]
        data = {
            "Text": text
        }
        body = json.dumps(data)
        resp = s.post("http://127.0.0.1:7000/words", data=body)
        if resp.status_code != requests.codes.ok and text != "":
            print(f"fail, {resp}, body={body}")
            sys.exit(1)

        rd = r.randint(10, 50)
        print(f"Sending {rd} GET requests")
        for k in range(r.randint(10, 50)):
            limit = k % len(available_words)
            if limit == 0:
                limit = 1
            data = {
                "Words": available_words[:limit]
            }
            body = json.dumps(data)
            resp = s.get("http://127.0.0.1:7000/words", data=body)
            if resp.status_code != requests.codes.ok and len(available_words[:limit]) > 0:
                print(f"fail, {resp}")
                sys.exit(1)

# final correctness check
for word in d:
    data = {
        "Words": [word]
    }
    body = json.dumps(data)
    resp = s.get("http://127.0.0.1:7000/words", data=body)
    if resp.json().get(word) != d.get(word):
        print(f"fail, expected {d.get(word)} appearances for word {word}, got {resp.json().get(word)}")
        sys.exit(1)

print("Test finished.")
