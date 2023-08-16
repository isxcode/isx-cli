import sys


def get_action():
    return sys.argv[1]


def get_project():
    return sys.argv[2]


def get_git_command():
    return ' '.join(sys.argv[1:])


def get_pr_num():
    return sys.argv[3]


def get_pr_message():
    return sys.argv[4]