import subprocess

from args import get_branch_num
from command import exec_command
from config import check_current_project
from config import get_current_project_path


def get_branch(num, origin):
    git_command = "git branch -r -l " + origin + "/*-#" + num
    result = subprocess.run(git_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True,
                            cwd=get_current_project_path())
    branches = result.stdout.split('\n')
    for branch in branches:
        branch_sub = branch.strip().split('-')
        if len(branch_sub) > 1 and branch_sub[1] == "#" + num:
            return branch.strip().replace(origin + "/", "")
    return ""


def get_local_branch(num):
    git_command = "git branch -l " + "*-#" + num
    result = subprocess.run(git_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True,
                            cwd=get_current_project_path())
    branches = result.stdout.split('\n')
    for branch in branches:
        branch_sub = branch.strip().split('-')
        if len(branch_sub) > 1 and branch_sub[1] == "#" + num:
            return branch.replace("*",'').strip()
    return ""


def get():
    check_current_project()
    exec_command("cd " + get_current_project_path() + " && git fetch upstream && git fetch origin")
    branch_name = get_branch(get_branch_num(), "origin")
    if branch_name != '':
        print("分支名称：" + branch_name)
    else:
        branch_name = get_branch(get_branch_num(), "upstream")
        if branch_name != '':
            print("分支名称：" + branch_name)
        else:
            branch_name = get_local_branch(get_branch_num())
            if branch_name != '':
                print("分支名称：" + branch_name)
            else:
                print("分支不存在，请使用isx branch创建分支")
