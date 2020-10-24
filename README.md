# Toggl reporter

[![PyPI version](https://badge.fury.io/py/toggl-reporter.svg)](https://badge.fury.io/py/toggl-reporter)

I need to write work report every evening. To make this routine simpler I created this CLI reporter.

## How to use

```shell
toggl_reporter --token=YOUR_TOKEN --date=2020-10-20
```

Script will ask you details if needed. By default date is today if you not specify it.

### Report

```text
REPORT for 2020-10-20
total time 7:00:00

my-cool-project 4:00:00
- create github repo
- write code
- test code

study-python 3:00:00
- do my homework
- watch lection
```

## Hot to install

You need a Python 3.

```shell
pip install toggl-reporter
```
