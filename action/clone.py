import re

import requests

from action.login import check_login
from config import get_config
from config import save_config
import os
import subprocess


def clone():
    check_login()
    isx_config = get_config()
    project_name = input_project_number(isx_config)
    if isx_config['projects'][project_name]['dir'] != '':
        print("该项目已下载，请重新执行【isx clone】选择")
        exit(0)
    project_dir = input_project_dir(project_name)
    clone_github_code(isx_config, project_name, project_dir)
    isx_config["projects"][project_name]['dir'] = project_dir + '/' + project_name
    save_config(isx_config)


def input_project_number(isx_config):
    for index, project_name in enumerate(isx_config["projects"]):
        project_meta = isx_config["projects"][project_name]
        print('[' + str(index) + '] ' + project_name + " : " + project_meta['describe'])
    project_number = input("请输入下载项目编号：")
    if project_number != '' and 0 <= int(project_number) < len(isx_config["projects"]) and project_number.isdigit():
        return list(isx_config["projects"])[int(project_number)]
    else:
        print("选择编号异常，请重新执行【isx clone】命令")
        exit(0)


def input_project_dir(project_name):
    project_dir = input("请输入下载目录(全路径)：")
    if not os.path.exists(project_dir):
        print("路径输入异常，请重新执行【isx clone】命令")
        exit(0)
    if project_dir.endswith("/"):
        project_dir = project_dir[:-1]
    if os.path.exists(project_dir + '/' + project_name):
        print("目录已存在，请删除重试")
        exit(0)
    return project_dir


def clone_github_code(isx_config, project_name, project_dir):
    account = isx_config['user']['account']
    token = isx_config['user']['token']
    project_info = isx_config['projects'][project_name]
    all_repository = list(project_info['sub-repository'])
    all_repository.append(project_info['repository'])
    for repository in all_repository:
        check_project_forked(repository, account, token)
    clone_code(project_info['repository'], project_dir, account, token)
    for repository in list(project_info['sub-repository']):
        clone_code(repository, project_dir + '/' + project_name, account, token)
    print("代码下载完成")


def check_project_forked(repository, account, token):
    matcher = re.search(r"/([^/.]+)\.git", repository)
    project_name = matcher.group(1)
    headers = {
        "Accept": "application/vnd.github+json",
        "Authorization": "Bearer " + token,
        "X-GitHub-Api-Version": "2022-11-28"
    }
    response = requests.get("https://api.github.com/repos/" + account + "/" + project_name, headers=headers)
    if response.status_code == 404:
        print("请先fork项目【 https://github.com/isxcode/" + project_name + "/fork 】，再重新执行【isx clone】命令")
        exit(0)
    elif response.status_code != 200:
        print("检测项目是否被fork异常，请重新执行【isx clone】命令")
        exit(0)


def clone_code(repository, project_dir, account, token):
    matcher = re.search(r"/([^/.]+)\.git", repository)
    project_name = matcher.group(1)
    self_repository = repository.replace('https://', 'https://' + token + '@').replace('isxcode', account)
    command = 'cd ' + project_dir + ' && git clone ' + self_repository
    completed_process = subprocess.run(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    os.system('cd ' + project_dir + '/' + project_name + ' && git remote add upstream ' + repository)
    if completed_process.returncode == 0:
        print(completed_process.stdout + '【' + project_name + '】下载成功')
    else:
        print(completed_process.stderr + '【' + project_name + '】下载失败')