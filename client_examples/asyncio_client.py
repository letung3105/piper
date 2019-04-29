import asyncio
import pathlib
import ssl
import websockets
import json
import time
import urllib.request
import urllib.parse

# SSL certificates locations
cacert = "../keys/certs/pub/cacert.pem"
clientcert = "../keys/certs/pub/clientcert.pem"
clientkey = "../keys/certs/priv/clientkey.pem"

# SSL context
sslContext = ssl.create_default_context(
    purpose=ssl.Purpose.SERVER_AUTH,
    cafile=cacert,
)
sslContext.load_cert_chain(clientcert, clientkey)

HOST = "127.0.0.1"
PORT = 4433

# Login url
url = "https://%s:%s/subscribe" % (HOST, PORT)


async def readPipe(ws):
    while True:
        packet = await ws.recv()
        data = json.loads(packet)
        print('Now:', time.time_ns())
        print('Sent:', data['time'])
        print('Elapsed:', "%fms" % ((time.time_ns() - data['time']) * 1e-6))
        time.sleep(0.05)

        await ws.send(json.dumps({ "finished": True }).encode('utf-8'))

async def getData(uri):
    # User login credentials
    creds = {
        "username": "user",
        "password": "password"
    }
    credsJSON = json.dumps(creds)
    credsJSONBytes = credsJSON.encode('utf-8')   # needs to be bytes

    # Make request
    req = urllib.request.Request(url)
    req.add_header('Content-Type', 'application/json; charset=utf-8')
    req.add_header('Content-Length', len(credsJSONBytes))

    # Send request
    resp = urllib.request.urlopen(req, credsJSONBytes, context=sslContext)
    respBody = resp.read()
    respBodyJSON = json.loads(respBody)

    headers = {
        "Authorization": 'Bearer %s' % (respBodyJSON['token'])
    }

    async with websockets.connect(uri, ssl=sslContext, extra_headers=headers) as websocket:
        await readPipe(websocket)

asyncio.get_event_loop().run_until_complete(
        getData("wss://%s:%s/data" % (HOST, PORT))
)
asyncio.get_event_loop().run_forever()