import datetime
from requests.auth import HTTPBasicAuth
import requests
from urllib.parse import urljoin
import argparse

class TogglAPI:
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
            'until': date
        }
        url = urljoin(self.base_url, path)
        response = requests.request("GET", url, auth=self.auth, params=query)
        response.raise_for_status()
        data = response.json()
        data['req'] = {
            'workspace_id': workspace_id,
            'date': date,
        }
        return data

    def get_me(self):
        path = '/api/v8/me'
        url = urljoin(self.base_url, path)
        response = requests.request("GET", url, auth=self.auth)
        response.raise_for_status()
        data = response.json()
        return data

def get_workspaces(user):
    return [{
        'id': w['id'],
        'name': w['name']
    } for w in user['data']['workspaces']]
    
def extract_entry(entry):
    return {
        'description': entry['description'],
        'start': entry['start'],
        'end': entry['end'],
        'dur': entry['dur'],
        'tags': entry['tags']
    }

def format_dur(milliseconds):
    return str(datetime.timedelta(milliseconds=milliseconds))

def format_date(date):
    d=datetime.datetime.fromisoformat(date)
    return d.strftime("%Y-%m-%d %A")

def build_report(data):
    total_dur = 0
    groups = {}
    for entry in data['data']:
        project = entry['project']
        entry_ = extract_entry(entry)
        date = entry_['start']
        total_dur += entry['dur']
        if groups.get(project) is None:
            groups[project] = {
                'project': project,
                'total_dur': entry['dur'],
                'entries': [entry_]
            }
        else:
            groups[project]['entries'].append(entry_)
            groups[project]['total_dur'] += entry['dur']

    return {
        'date': data['req']['date'],
        'total_dur': total_dur,
        'groups': groups
    }

def print_report(report):
    print(f"REPORT for {report['date']}")
    print(f"total time {format_dur(report['total_dur'])}")
    print("")
    for proj_name in report['groups']:
        group = report['groups'][proj_name]
        dur = format_dur(group['total_dur'])
        print(f"{proj_name} {dur}")
        group['entries'] = sorted(group['entries'], key=lambda el: el['start'])

        for entry in group['entries']:
            print(f"- {entry['description']}")
        print("")

def parse_args():
    parser = argparse.ArgumentParser(description='Toggl Reporter!')
    parser.add_argument('--token', required=True, help='API Token. Copy it from https://track.toggl.com/profile')
    parser.add_argument('--date', required=False, help='Which date make report for. Format 2020-10-27')
    parser.add_argument('--workspace', required=False, help='Workspace ID. You can left it empty if you have only one workspace')
    return parser.parse_args()

def select_ws(args, api):
    if args.workspace is not None and len(args.workspace) != 0:
        return args.workspace

    ws = get_workspaces(api.get_me())
    if len(ws) == 0:
        print("Create workspace in toggl")
        exit(0)
    if len(ws) > 1:
        print("You have many workspaces. Please set --workspace argument with an ID you need.")
        for w in ws:
            print(f"{w['name']}: {w['id']}")
        exit(0)
    if len(ws) == 1:
        return ws[0]['id']

def select_date(args):
    if args.date is not None and len(args.date) != 0:
        return args.date

    return datetime.datetime.now().strftime("%Y-%m-%d")

def main(args):
    api = TogglAPI(args.token)

    date = select_date(args)
    workspace_id = select_ws(args, api)

    data = api.get_report(workspace_id, date)
    report = build_report(data)
    print_report(report)

if __name__ == "__main__":
    main(parse_args())
