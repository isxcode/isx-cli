from config import get_current_project
from config import get_token
from command import exec_command
from config import get_current_project_path
from config import get_current_project_vip_path
from config import get_account
from args import get_pr_num
from args import get_pr_message
import requests
import subprocess


def get_branch():
    try:
        # 先刷新upstream
        exec_command("cd " + get_current_project_path() + " && git fetch upstream && git fetch origin")
        exec_command("cd " + get_current_project_vip_path() + " && git fetch upstream && git fetch origin")
        result = subprocess.run(["git", "branch", "-r", "-l", "upstream/*-#" + get_pr_num()],
                                stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True,
                                check=True, cwd=get_current_project_path())
        branches = result.stdout.split('\n')
        for branch in branches:
            branch_sub = branch.strip().split('-')
            if len(branch_sub) > 1 and branch_sub[1] == "#"+get_pr_num():
                return branch.strip().replace("upstream/", "")
        print("未找到分支")
        exit(0)
    except subprocess.CalledProcessError as e:
        print("找寻分支异常")
        exit(0)


def pr():
    branch_name = get_branch()
    url = "https://api.github.com/repos/isxcode/" + get_current_project() + "/pulls"
    headers = {
        "Accept": "application/vnd.github+json",
        "Authorization": "Bearer " + get_token(),
        "X-GitHub-Api-Version": "2022-11-28"
    }
    data = {
        "title": "#" + get_pr_num() + " " + get_pr_message().replace("'", ""),
        "head": branch_name,
        "head_repo": get_account() + "/" + get_current_project(),
        "base": branch_name
    }
    response = requests.post(url, json=data, headers=headers)
    if response.status_code != 201:
        print("提交失败")
    else:
        print("提交成功")
    url = "https://api.github.com/repos/isxcode/" + get_current_project() + "-vip/pulls"
    headers = {
        "Accept": "application/vnd.github+json",
        "Authorization": "Bearer " + get_token(),
        "X-GitHub-Api-Version": "2022-11-28"
    }
    data = {
        "title": "#" + get_pr_num() + " " +get_pr_message(),
        "head": branch_name,
        "head_repo": get_account() + "/" + get_current_project() + "-vip",
        "base": branch_name
    }
    response = requests.post(url, json=data, headers=headers)
    if response.status_code != 201:
        print("提交失败")
    else:
        print("提交成功")