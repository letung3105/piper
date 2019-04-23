import random
import time
import json
import sys


COLORS = ['k', 'r', 'g', 'b', 'y']
FPS = 20
names = []

def randData():
    frame = []
    for name in names:
        position = [random.randint(0, 1500), random.randint(0, 1500)]
        dimension = [random.randint(-90, 90), random.randint(-90, 90)]
        frame.append({
            'name': name,
            'position': position,
            'dimension': dimension
        })

    return frame


def main():
    for i in range(2):
        names.append("manual_%d" % (i + 1))
        names.append("auto_%d" % (i + 1))

    for i in range(24):
        names.append("object_%d" % (i + 1))

    while True:
        print(json.dumps(randData()))
        sys.stdout.flush()
        time.sleep(1 / FPS)


if __name__ == '__main__':
    main()
