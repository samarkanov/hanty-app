import requests
import unittest
import json
from pdb import set_trace as stop

HOST = "http://develop.valenoq.com:10003"


class TestDB(unittest.TestCase):

    def setUp(self):
        url = "{host}/client_x".format(host=HOST)
        res = requests.delete(url)
        self.assertTrue(json.loads(res.text)["success"])

        url = "{host}/client_y".format(host=HOST)
        res = requests.delete(url)
        self.assertTrue(json.loads(res.text)["success"])

    def store_in_db(self):
        params = dict()
        params["topic"] = "topic_x"
        params["value"] = "value_x"
        url = "{host}/client_x".format(host=HOST)
        res = requests.post(url, params)
        self.assertTrue(res.ok)
        self.assertTrue(json.loads(res.text)["success"])

        params["topic"] = "topic_x"
        params["value"] = "value_y"
        url = "{host}/client_y".format(host=HOST)
        res = requests.post(url, params)
        self.assertTrue(json.loads(res.text)["success"])

        params["topic"] = "topic_m"
        params["value"] = "value_m"
        url = "{host}/client_y".format(host=HOST)
        res = requests.post(url, params)
        self.assertTrue(json.loads(res.text)["success"])

        params["topic"] = "topic_m"
        params["value"] = "value_b"
        url = "{host}/client_y".format(host=HOST)
        res = requests.post(url, params)
        self.assertTrue(json.loads(res.text)["success"])

    def test_gw(self):
        print("TestDB::test_gw")

        # storing values with POST
        self.store_in_db()

        # getting values with GET
        url = "{host}/client_y".format(host=HOST)
        res = requests.get(url)
        res = json.loads(res.text)
        self.assertTrue({'Topic': 'topic_x', 'Data': ['value_y']} in res)
        self.assertTrue({'Topic': 'topic_m', 'Data': ['value_m', 'value_b']} in res)
        self.assertEqual(len(res), 2)

        url = "{host}/client_x".format(host=HOST)
        res = requests.get(url)
        res = json.loads(res.text)
        self.assertTrue({'Topic': 'topic_x', 'Data': ['value_x']} in res)
        self.assertEqual(len(res), 1)

        # deleting values with DELETE
        url = "{host}/client_y/topic_m".format(host=HOST)
        requests.delete(url)

        url = "{host}/client_x".format(host=HOST)
        requests.delete(url)

        # check after delete
        url = "{host}/client_y".format(host=HOST)
        res = requests.get(url)
        res = json.loads(res.text)
        self.assertTrue({'Topic': 'topic_x', 'Data': ['value_y']} in res)
        self.assertEqual(len(res), 1)

        url = "{host}/client_x".format(host=HOST)
        res = requests.get(url)
        res = json.loads(res.text)
        self.assertTrue(res is None)


if __name__ == "__main__":
    unittest.main()
