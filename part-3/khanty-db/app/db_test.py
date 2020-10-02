import requests
import unittest
from pdb import set_trace as stop

HOST = "http://develop.valenoq.com:10003"


class TestPOST(unittest.TestCase):
    
    def test_good_weather(self):
        print("TestPOST::test_good_weather")
        params = dict()
        params["topic"] = "topic_x"
        params["value"] = "value_x"
        url = "{host}/client_x".format(host=HOST)
        res = requests.post(url, params)
        self.assertEqual('about to store an entry in database: client_x:topic_x:value_x', res.text)
        self.assertTrue(res.ok)

    def test_empty_params(self):
        print("TestPOST::test_empty_params")
        url = "{host}/client_x".format(host=HOST)
        res = requests.post(url, dict())
        self.assertTrue(res.ok)


class TestGET(unittest.TestCase):

    def test_data_per_client_gw(self):
        print("TestGET::test_data_per_client_gw")
        url = "{host}/client_x".format(host=HOST)
        res = requests.get(url, dict())
        self.assertEqual('returning everything I have in DB for client = client_x', res.text)
        self.assertTrue(res.ok)

    def test_data_per_client_and_topic_gw(self):
        print("TestGET::test_data_per_client_and_topic_gw")
        url = "{host}/client_x/topic_y".format(host=HOST)
        res = requests.get(url, dict())
        self.assertTrue(res.ok)
        self.assertEqual('returning data for client = client_x and topic = topic_y', res.text)


class TestDELETE(unittest.TestCase):

    def test_data_per_client_gw(self):
        print("TestDELETE::test_data_per_client_gw")
        url = "{host}/client_x".format(host=HOST)
        res = requests.delete(url)
        self.assertEqual('deleting all entries for client = client_x', res.text)
        self.assertTrue(res.ok)

    def test_data_per_client_and_topic_gw(self):
        print("TestDELETE::test_data_per_client_and_topic_gw")
        url = "{host}/client_x/topic_y".format(host=HOST)
        res = requests.delete(url)
        self.assertTrue(res.ok)
        self.assertEqual('deleting all entries for client = client_x and topic = topic_y', res.text)


if __name__ == "__main__":
    unittest.main()
