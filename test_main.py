import requests
import json


def test_save_data():
    url = 'http://localhost:5000/save_data/'

    JSON = "application/json"
    header = {'Content-Type': JSON}
    data = {
        'id_user': '88888888',
        'email': 'test@mail/ru',
        'amount': 1500,
        'currency': 'USD ($)'
    }
    resp = requests.post(url, headers=header, params=data)

    assert resp.status_code == 200
    resp_body = resp.json()
    print(resp)


def test_show_transaction_by_id():
    url = 'http://localhost:5000/show_transaction/'

    JSON = "application/json"
    header = {'Content-Type': JSON}
    data = {
        'id_user': '88888888'
    }
    resp = requests.post(url, headers=header, params=data)

    assert resp.status_code == 200
    resp_body = resp.json()
    print(resp)


def test_show_transaction_by_email():
    url = 'http://localhost:5000/show_transaction/'

    JSON = "application/json"
    header = {'Content-Type': JSON}
    data = {
        'email': 'test@mail.ru'
    }
    resp = requests.post(url, headers=header, params=data)

    assert resp.status_code == 200
    resp_body = resp.json()
    print(resp)


def test_show_last_check():
    url = 'http://localhost:5000/show_last_check/'

    JSON = "application/json"
    header = {'Content-Type': JSON}
    data = {
        'id_user': '88888888'
    }
    resp = requests.post(url, headers=header, params=data)

    assert resp.status_code == 200
    resp_body = resp.json()
    print(resp)


def test_cancel_completed():
    url = 'http://localhost:5000/cancel_completed/'

    JSON = "application/json"
    header = {'Content-Type': JSON}
    data = {
        'id_user': '88888888'
    }
    resp = requests.post(url, headers=header, params=data)

    assert resp.status_code == 200
    resp_body = resp.json()
    print(resp)

def test_all():
    test_save_data()
    test_show_transaction_by_id()
    test_show_transaction_by_email()
    test_show_last_check()
    test_cancel_completed()
    return 


test_all()
# def test_status_checker():
#     url = 'http://localhost:5000/status_checker/'
#
#     data = {
#         'status': 'УСПЕХ',
#         'token': 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTU0NjgzNDd9.UcT_Swc_DL36kbo_G4rV_F8Oj2_nNamvuJRkKQhmAUg',
#     }
#     resp = requests.post(url, params=data)
#
#     assert resp.status_code == 200
#     print(resp.text)
#
#
# test_status_checker()