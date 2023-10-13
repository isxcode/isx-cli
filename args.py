import sys


def get_action():
    return sys.argv[1]


def get_git_command():
    return ' '.join(sys.argv[1:])


def get_pr_num():
    return sys.argv[2]


def get_branch_num():
    if len(sys.argv) < 3:
        print("请输入正确的数字")
        exit(0)
    branch_num = sys.argv[2]
    if branch_num.isdigit() and int(branch_num) > 0:
        return branch_num
    else:
        print("请输入正确的数字")
        exit(0)


def get_pr_message():
    return sys.argv[3]
