import math
from datetime import datetime, timedelta


def round_to_minutes(milliseconds):
    to_minutes = 60000
    return math.ceil(milliseconds / to_minutes) * to_minutes


def format_dur(milliseconds):
    return str(timedelta(milliseconds=milliseconds))


def format_date(date):
    date = datetime.fromisoformat(date)
    return '{:%Y-%m-%d %A}'.format(date)  # noqa: P101, WPS323
