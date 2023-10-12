import json
import re
import subprocess

import requests

from args import get_branch_num
from command import exec_command
from config import check_current_project
from config import get_config
from config import get_token


def get_branch(branch_num, project_dir, prefix):
    if prefix != 'local':
        prefix = prefix + '/'
    search_git_command = "git branch -l " + prefix + "*-#" + branch_num
    result = subprocess.run(search_git_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True,
                            cwd=project_dir)
    branches = result.stdout.split('\n')
    for branch_meta in branches:
        branch_sub = branch_meta.strip().split('-')
        if len(branch_sub) > 1 and branch_sub[1] == "#" + branch_num:
            return branch_meta.replace("*", '').strip()
    return ""


def checkout_branch(project_dir, sub_repository, branch_name):
    exec_command("cd " + project_dir + " && git checkout " + branch_name)
    for repository_meta in sub_repository:
        matcher = re.search(r"/([^/.]+)\.git", repository_meta)
        project_name = matcher.group(1)
        exec_command("cd " + project_dir + '/' + project_name + " && git checkout " + branch_name)


def fetch_branch(project_dir, sub_repository):
    exec_command("cd " + project_dir + " && git fetch upstream && git fetch origin")
    for repository_meta in sub_repository:
        matcher = re.search(r"/([^/.]+)\.git", repository_meta)
        project_name = matcher.group(1)
        exec_command("cd " + project_dir + '/' + project_name + " && git fetch upstream && git fetch origin")


def push_branch(project_dir, sub_repository, branch_name, prefix):
    exec_command("cd " + project_dir + " && git push origin " + branch_name)
    for repository_meta in sub_repository:
        matcher = re.search(r"/([^/.]+)\.git", repository_meta)
        project_name = matcher.group(1)
        exec_command("cd " + project_dir + '/' + project_name + " && git push " + prefix + " " + branch_name)


def get_latest_branch():
    url = "https://api.github.com/repos/isxcode/spark-yun/releases/latest"
    headers = {
        "Accept": "application/vnd.github+json",
        "Authorization": "Bearer " + get_token(),
        "X-GitHub-Api-Version": "2022-11-28"
    }
    response = requests.get(url, headers=headers)
    if response.status_code == 401:
        print("github token权限不足，请重新登录")
        exit(0)
    else:
        release_name = json.loads(response.text)['name'].replace("v", '')
        return release_name + "-#" + get_branch_num()


def branch():
    isx_config = get_config()
    check_current_project(isx_config)
    project_info = isx_config['projects'][isx_config['develop-project']]
    project_dir = project_info['dir']
    sub_repository = list(project_info['sub-repository'])
    fetch_branch(project_dir, sub_repository)
    branch_num = get_branch_num()

    # 从本地拉分支
    branch_name = get_branch(branch_num, project_dir, 'local')
    if branch_name != '':
        checkout_branch(project_dir, sub_repository, branch_name)
        print("分支名:" + branch_name)
        exit(0)

    # 从个人仓库拉分支
    branch_name = get_branch(branch_num, project_dir, 'origin')
    if branch_name != '':
        branch_name_command = " -b " + branch_name + " origin/" + branch_name
        checkout_branch(project_dir, sub_repository, branch_name_command)
        print("分支名:" + branch_name)
        exit(0)

    # 从upstream远程仓库拉分支
    branch_name = get_branch(branch_num, project_dir, 'upstream')
    if branch_name != '':
        branch_name_command = " -b " + branch_name + " upstream/" + branch_name
        checkout_branch(project_dir, sub_repository, branch_name_command)
        push_branch(branch_name, 'origin')
        print("分支名:" + branch_name)
    else:
        branch_name_command = " -b " + get_latest_branch() + " upstream/latest"
        checkout_branch(project_dir, sub_repository, branch_name_command)
        push_branch(branch_name, 'origin')
        push_branch(branch_name, 'upstream')
        print("分支名:" + branch_name)
