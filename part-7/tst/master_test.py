import requests
import unittest
import json
import subprocess
import os.path
import time
from pdb import set_trace as stop

HOST_MASTER = "http://develop.valenoq.com:10004"
HOST_DB = "http://develop.valenoq.com:10003"


class TestMaster(unittest.TestCase):

    def setUp(self):
        url = "{host}/master".format(host=HOST_DB)
        res = requests.delete(url)
        self.assertTrue(json.loads(res.text)["success"])

    def test_check_subscription(self):
        print("TestMaster::test_check_subscription")
        params = dict()

        # subscribe
        params["topic"] = "topic_x"
        params["portno"] = "2222"
        url = "{host}/subscribe".format(host=HOST_MASTER)
        res = requests.post(url, params)

        # check suscription
        url = "{host}/master".format(host=HOST_DB)
        res = requests.get(url)
        self.assertTrue(res.ok)
        self.assertEqual([{"Topic":"topic_x","Data":["2222"]}], json.loads(res.text))

        # subscribe: same topic but another portno
        params["topic"] = "topic_x"
        params["portno"] = "2223"
        url = "{host}/subscribe".format(host=HOST_MASTER)
        res = requests.post(url, params)

        # check suscription
        url = "{host}/master".format(host=HOST_DB)
        res = requests.get(url)
        self.assertTrue(res.ok)
        self.assertEqual([{"Topic":"topic_x","Data":["2222", "2223"]}], json.loads(res.text))

    def test_check_notification(self):
        print("TestMaster::test_check_subscription")
        params = dict()
        expectation = dict()

        # check that nothing is in DB
        url = "{host}/master".format(host=HOST_DB)
        res = requests.get(url)
        self.assertTrue(res.ok)
        self.assertEqual(None, json.loads(res.text))

        # subscribe client 7777
        params["topic"] = "topic_x"
        params["portno"] = "7777"
        url = "{host}/subscribe".format(host=HOST_MASTER)
        res = requests.post(url, params)
        self.assertTrue(res.ok)

        # subscribe client 8888
        params["topic"] = "topic_x"
        params["portno"] = "8888"
        url = "{host}/subscribe".format(host=HOST_MASTER)
        res = requests.post(url, params)
        self.assertTrue(res.ok)

        # subscribe client 9999
        params["topic"] = "topic_x"
        params["portno"] = "9999"
        url = "{host}/subscribe".format(host=HOST_MASTER)
        res = requests.post(url, params)
        self.assertTrue(res.ok)

        # notify
        params["topic"] = "topic_x"
        params["value"] = "some_value"
        url = "{host}/notify".format(host=HOST_MASTER)
        res = requests.post(url, params)
        expectation["success"] = True
        expectation["clients"] = ["7777", "8888", "9999"]
        expectation["message"] = "clients notified OK"
        self.assertTrue(res.ok)
        self.assertEqual(json.loads(res.text), expectation)

        # unsubscribe client 9999
        url = "{host}/unsubscribe/topic_x/9999".format(host=HOST_MASTER)
        res = requests.delete(url)
        self.assertTrue(res.ok)

        # notify: client 9999 must not be notified
        params["topic"] = "topic_x"
        params["value"] = "xxx"
        url = "{host}/notify".format(host=HOST_MASTER)
        res = requests.post(url, params)
        expectation["success"] = True
        expectation["clients"] = ["7777", "8888"]
        expectation["message"] = "clients notified OK"
        self.assertTrue(res.ok)
        self.assertEqual(json.loads(res.text), expectation)

        # unsubscribe non-existing client
        url = "{host}/unsubscribe/topic_x/9999".format(host=HOST_MASTER)
        res = requests.delete(url)
        self.assertTrue(res.ok)

        # notify
        params["topic"] = "topic_x"
        params["value"] = "xxx"
        url = "{host}/notify".format(host=HOST_MASTER)
        res = requests.post(url, params)
        expectation["success"] = True
        expectation["clients"] = ["7777", "8888"]
        expectation["message"] = "clients notified OK"
        self.assertTrue(res.ok)
        self.assertEqual(json.loads(res.text), expectation)

        # unsubscribe client 8888
        url = "{host}/unsubscribe/topic_x/8888".format(host=HOST_MASTER)
        res = requests.delete(url)
        self.assertTrue(res.ok)

        # notify: client 8888 must not be notified
        params["topic"] = "topic_x"
        params["value"] = "yyyy"
        url = "{host}/notify".format(host=HOST_MASTER)
        res = requests.post(url, params)
        expectation["success"] = True
        expectation["clients"] = ["7777"]
        expectation["message"] = "clients notified OK"
        self.assertTrue(res.ok)
        self.assertEqual(json.loads(res.text), expectation)

    def test_notify_if_no_clients_subscribed(self):
        print("TestMaster::test_notify_no_clients_subscribed")
        params = dict()
        expectation = dict()

        # notify: must not lead to error
        params["topic"] = "topic_x"
        params["value"] = "yyyy"
        url = "{host}/notify".format(host=HOST_MASTER)
        res = requests.post(url, params)
        self.assertTrue(res.ok)
        expectation["success"] = True
        expectation["clients"] = None
        expectation["message"] = "clients notified OK"
        self.assertEqual(json.loads(res.text), expectation)


if __name__ == "__main__":
    unittest.main()
