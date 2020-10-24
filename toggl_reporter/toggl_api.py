from urllib.parse import urljoin

import requests
from requests.auth import HTTPBasicAuth


class TogglAPI(object):
    def __init__(self, token: str):
        self.token = token
        self.auth = HTTPBasicAuth(token, 'api_token')
        self.base_url = 'https://api.track.toggl.com'
        self.app_name = 'senior_sigan_reporter'

    def get_report(self, workspace_id, date):
        path = '/reports/api/v2/details'
        query = {
            'user_agent': self.app_name,
            'workspace_id': workspace_id,
            'since': date,
            'until': date,
        }
        url = urljoin(self.base_url, path)
        response = requests.request('GET', url, auth=self.auth, params=query)
        response.raise_for_status()
        report_data = response.json()
        report_data['req'] = {
            'workspace_id': workspace_id,
            'date': date,
        }
        return report_data

    def get_me(self):
        path = '/api/v8/me'
        url = urljoin(self.base_url, path)
        response = requests.request('GET', url, auth=self.auth)
        response.raise_for_status()
        return response.json()
