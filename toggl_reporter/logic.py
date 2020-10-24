from toggl_reporter.dates import format_dur, round_to_minutes


def get_workspaces(user):
    return [
        {
            'id': workspace['id'],
            'name': workspace['name'],
        } for workspace in user['data']['workspaces']
    ]


def extract_entry(entry):
    return {
        'description': entry['description'],
        'start': entry['start'],
        'end': entry['end'],
        'dur': round_to_minutes(entry['dur']),
        'tags': entry['tags'],
    }


def build_report(toggl_report):
    total_dur = 0
    groups = {}
    for el in toggl_report['data']:
        project = el['project']
        entry = extract_entry(el)
        total_dur += entry['dur']
        if groups.get(project) is None:
            groups[project] = {
                'name': project,
                'total_dur': entry['dur'],
                'entries': [entry],
            }
        else:
            groups[project]['entries'].append(entry)
            groups[project]['total_dur'] += entry['dur']

    return {
        'date': toggl_report['req']['date'],
        'total_dur': total_dur,
        'groups': groups,
    }


def group_entries_by_description_and_sum_dur(entries):
    groups = {}
    for entry in entries:
        desk = entry['description']
        if groups.get(desk) is None:
            groups[desk] = entry
        else:
            groups[desk]['dur'] += entry['dur']
    return sorted(groups.values(), key=lambda el: el['start'])


def print_project(project):
    dur = format_dur(project['total_dur'])
    print('{0} {1}'.format(project['name'], dur))
    entries = group_entries_by_description_and_sum_dur(project['entries'])
    for entry in entries:
        duration = format_dur(entry['dur'])
        desc = entry['description']
        print('- {0} - {1}'.format(desc, duration))


def print_report(report):
    total_time = format_dur(report['total_dur'])
    print('REPORT for {0}'.format(report['date']))
    print('total time {0}'.format(total_time))
    print('')
    for _, project in report['groups'].items():
        print_project(project)
        print('')
