import requests
import unittest
import json
import time
import subprocess
from pdb import set_trace as stop

HOST_CONFIG = "http://develop.valenoq.com"
PORT_CONFIG = "10002"
CLIENTS = ["clientA", "clientB", "clientC", "clientD"]
HOST_DB = "http://develop.valenoq.com:10003"
TOPICS = ["ChangeColor", "SendMessage"]

class TestClient(unittest.TestCase):

    def setUp(self):
        reply = requests.get("{host}:{port}".format(host=HOST_CONFIG,
                                                    port=PORT_CONFIG))
        self.config = json.loads(reply.text)

        # cleaning DB:
        for client in CLIENTS:
            url = "{host}/{client}".format(host=HOST_DB, client=client)
            res = requests.delete(url)
            self.assertTrue(json.loads(res.text)["success"])

        # start clients
        proc = subprocess.check_output("./start.sh", shell=True)

    def tearDown(self):
        # stop client
        proc = subprocess.check_output("./stop.sh", shell=True)

    def test_initial_state(self):
        print("TestClient::test_initial_state")
        for client in CLIENTS:
            host = self.config[client]["host"]
            port = self.config[client]["port"]
            res = requests.get("{host}:{port}".format(host=host, port=port))
            self.assertTrue(res.ok)
            res = json.loads(res.text)
            self.assertEqual(len(res), 2)
            self.assertTrue({'Topic': 'ChangeColor', 'Data': ['#5e72e4']} in res)
            self.assertTrue({'Topic': 'SendMessage', 'Data': ['#SayHi']} in res)
            print("\tclient %s is OK" %client)

    def test_current_state(self):
        print("TestClient::test_current_state")
        for client in CLIENTS:
            host = self.config[client]["host"]
            port = self.config[client]["port"]
            for topic in TOPICS:
                ctx = dict()
                ctx["host"] = host
                ctx["port"] = port
                ctx["topic"] = topic
                res = requests.get("{host}:{port}/{topic}".format(**ctx))
                self.assertTrue(res.ok)
                self.assertTrue(len(res.text) > 5)
            print("\tclient %s is OK" %client)

    def test_update(self):
        print("TestClient::test_update")
        params = dict()
        client = "clientB"
        host = self.config[client]["host"]
        port = self.config[client]["port"]

        # update property
        params["topic"] = "ChangeColor"
        params["value"] = "#xxx"
        res = requests.post("{host}:{port}/update".format(host=host, port=port), params)
        self.assertTrue(res.ok)

        # check after the update
        res = requests.get("{host}:{port}/ChangeColor".format(host=host, port=port))
        self.assertEqual(res.text, "#xxx")

        # update property
        params["topic"] = "SendMessage"
        params["value"] = "xxx"
        res = requests.post("{host}:{port}/update".format(host=host, port=port), params)
        self.assertTrue(res.ok)

        params["topic"] = "SendMessage"
        params["value"] = "yyy mm"
        res = requests.post("{host}:{port}/update".format(host=host, port=port), params)
        self.assertTrue(res.ok)

        # check after the updates
        res = requests.get("{host}:{port}/SendMessage".format(host=host, port=port))
        self.assertEqual(res.text, "yyy mm")
        self.assertTrue(res.ok)

        # update property (insert existing -> the duplicate should not be created)
        # the newly inserted property must be a tail element of an array
        params["topic"] = "SendMessage"
        params["value"] = "#SayHi"
        res = requests.post("{host}:{port}/update".format(host=host, port=port), params)
        self.assertTrue(res.ok)

        res = requests.get("{host}:{port}".format(host=host, port=port))
        res = json.loads(res.text)
        self.assertEqual(len(res), 2)
        self.assertTrue({"Topic":"ChangeColor","Data":["#5e72e4","#xxx"]} in res)
        self.assertTrue({"Topic":"SendMessage","Data":["yyy mm","xxx","#SayHi"]} in res)


if __name__ == "__main__":
    unittest.main()
