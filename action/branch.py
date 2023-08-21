import subprocess

import requests
import json

from command import exec_command
from config import get_current_project_path
from config import get_current_project_vip_path
from config import check_current_project
from args import get_branch_num
from config import get_token
import os


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


def exec_git_command(git_command):
    print("执行git命令：" + git_command)
    command = "cd " + get_current_project_path() + " && " + git_command
    os.system(command)
    command = "cd " + get_current_project_vip_path() + " && " + git_command
    os.system(command)
    command = "cd " + get_current_project_path()
    os.system(command)


def branch():
    check_current_project()
    exec_command("cd " + get_current_project_path() + " && git fetch upstream && git fetch origin")
    exec_command("cd " + get_current_project_vip_path() + " && git fetch upstream && git fetch origin")
    branch_name = get_branch(get_branch_num(), "origin")
    if branch_name != '':
        print("切换分支：" + branch_name)
        exec_git_command("git checkout --track origin/" + branch_name)
    else:
        branch_name = get_branch(get_branch_num(), "upstream")
        if branch_name != '':
            print("切换分支：" + branch_name)
            exec_git_command("git checkout --track upstream/" + branch_name)
        else:
            # 如果自己仓库和upstream仓库都没有分支，则需要创建分支
            url = "https://api.github.com/repos/isxcode/spark-yun/releases/latest"
            headers = {
                "Accept": "application/vnd.github+json",
                "Authorization": "Bearer " + get_token(),
                "X-GitHub-Api-Version": "2022-11-28"
            }
            response = requests.get(url, headers=headers)
            if response.status_code == 401:
                print("github token权限不足，请重新登录")
            else:
                release_name = json.loads(response.text)['name'].replace("v", '')
                branch_name = release_name + "-#" + get_branch_num()
                print("创建分支：" + branch_name)
                exec_git_command("git checkout -b " + branch_name + " upstream/latest")