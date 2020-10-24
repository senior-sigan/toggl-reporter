import argparse
import datetime

from toggl_reporter.logic import build_report, get_workspaces, print_report
from toggl_reporter.toggl_api import TogglAPI


def parse_args():
    parser = argparse.ArgumentParser(description='Toggl Reporter!')
    parser.add_argument(
        '--token',
        required=True,
        help='API Token. Copy it from https://track.toggl.com/profile',
    )
    parser.add_argument(
        '--date',
        required=False,
        help='Which date make report for. Format 2020-10-27',
    )
    parser.add_argument(
        '--workspace',
        required=False,
        help='Workspace ID. Left it empty if you have only one workspace',
    )
    return parser.parse_args()


def select_ws(args, api):
    if args.workspace is not None:
        print('You need at least one workspace. Create workspace in toggl.')
        return args.workspace

    ws = get_workspaces(api.get_me())
    if len(ws) > 1:
        print(
            'You have many workspaces.',
            'Please set --workspace argument with an ID you need.',
        )
        for workspace in ws:
            print(f'{workspace["name"]}: {workspace["id"]}')
        return None
    if len(ws) == 1:
        return ws[0]['id']


def select_date(args):
    if args.date is not None:
        return args.date

    now = datetime.datetime.now()
    return '{:%Y-%m-%d}'.format(now)  # noqa: WPS323, P101


def run_with_args(args):
    api = TogglAPI(args.token)

    date = select_date(args)
    workspace_id = select_ws(args, api)

    report = api.get_report(workspace_id, date)
    report = build_report(report)
    print_report(report)


def run():
    run_with_args(parse_args())
