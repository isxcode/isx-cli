import os

from config import get_current_project
from config import get_token
from config import get_current_project_path
from config import get_account
from args import get_pr_num
from args import get_pr_message
import requests
import os
import subprocess


def pr():
    command = ["git", "branch", "--show-current"]
    completed_process = subprocess.run(command)
    print("<" + completed_process + ">")


    # url = "https://api.github.com/repos/isxcode/" + get_current_project() + "/pulls"
    # headers = {
    #     "Accept": "application/vnd.github+json",
    #     "Authorization": " Bearer " + get_token(),
    #     "X-GitHub-Api-Version": "2022-11-28"
    # }
    # data = {
    #     "title": get_pr_message(),
    #     "head": get_pr_num(),
    #     "head_repo": get_account() + "/" + get_current_project(),
    #     "base": get_pr_num()
    # }
    # # vip模块
    # response = requests.post(url, headers=headers, json=data)
    # url = "https://api.github.com/repos/isxcode/" + get_current_project() + "-vip/pulls"
    # headers = {
    #     "Accept": "application/vnd.github+json",
    #     "Authorization": " Bearer " + get_token(),
    #     "X-GitHub-Api-Version": "2022-11-28"
    # }
    # data = {
    #     "title": get_pr_message(),
    #     "head": get_pr_num(),
    #     "head_repo": get_account() + "/" + get_current_project() + "-vip",
    #     "base": get_pr_num()
    # }
    # response = requests.post(url, headers=headers, json=data)
