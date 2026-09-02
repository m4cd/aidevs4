#!/usr/bin/env python3
"""
Windpower - pelna sekwencja zadania.

Prognoza jest losowana przy kazdym 'start', wiec unlockCode musi byc
liczony po odebraniu weather, w tej samej sesji.

Uzycie:  python3 windpower.py https://adres-endpointa/ APIkey
"""

import hashlib
import json
import sys
import time
import urllib.request

TASK = "windpower"
CUTOFF = 14.0   # safety.cutoffWindMs - powyzej: chorągiewka
MIN_OP = 4.0    # safety.minOperationalWindMs

URL = sys.argv[1] if len(sys.argv) > 2 else None
APIKEY = sys.argv[2] if len(sys.argv) > 2 else None

if not URL:
    sys.exit("Podaj URL endpointa jako argument.")

if not APIKEY:
    sys.exit("Podaj klucz API endpointa jako argument.")

def post(answer):
    payload = json.dumps({"apikey": APIKEY, "task": TASK, "answer": answer}).encode()
    req = urllib.request.Request(
        URL, data=payload, headers={"Content-Type": "application/json"}
    )
    with urllib.request.urlopen(req, timeout=15) as r:
        return json.loads(r.read())


def unlock(date, hour, wind_ms, pitch):
    """Format potwierdzony przez signedParams: obie liczby jako float."""
    s = f"{date}|{hour}|{float(wind_ms)}|{float(pitch)}"
    return hashlib.md5(s.encode()).hexdigest()


t0 = time.time()

# 1. start
print(post({"action": "start"}))

# 2-4. zakolejkowanie raportow, bez czekania na wyniki
for param in ("weather", "turbinecheck", "powerplantcheck"):
    post({"action": "get", "param": param})

# 5-7. odbior - getResult zwraca po jednym, rozroznienie po sourceFunction
results = {}
while len(results) < 3:
    r = post({"action": "getResult"})
    src = r.get("sourceFunction")
    if src:
        results[src] = r
    else:
        time.sleep(0.3)

weather = results["weather"]
print("turbinecheck:", json.dumps(results["turbinecheck"], ensure_ascii=False))
print("powerplantcheck:", json.dumps(results["powerplantcheck"], ensure_ascii=False))

# 8. budowa konfiguracji z ZYWEJ prognozy
configs = []
best = None

for row in weather["forecast"]:
    date, hour = row["timestamp"].split(" ")
    wind = row["windMs"]

    if wind >= CUTOFF:
        # sztorm - zabezpieczenie turbiny
        configs.append({
            "startDate": date,
            "startHour": hour,
            "pitchAngle": 90,
            "turbineMode": "idle",
            "unlockCode": unlock(date, hour, wind, 90),
        })
    elif wind >= MIN_OP and (best is None or wind > best["windMs"]):
        best = {"date": date, "hour": hour, "windMs": wind}

if best is None:
    sys.exit("Brak wpisu w zakresie operacyjnym (4 - 14 m/s).")

# punkt produkcyjny - pitch 0 = 100% uzysku
configs.append({
    "startDate": best["date"],
    "startHour": best["hour"],
    "pitchAngle": 0,
    "turbineMode": "production",
    "unlockCode": unlock(best["date"], best["hour"], best["windMs"], 0),
})

print(f"sztormy: {len(configs) - 1}, produkcja: {best}")
print(post({"action": "config", "configs": configs}))

# 9. done
print(post({"action": "done"}))
print(f"czas: {time.time() - t0:.1f}s")
