import re

import requests

from action.get import get_branch_name
from args import get_pr_num, get_pr_message
from config import get_account
from config import get_config
from config import get_token


def pr():
    isx_config = get_config()
    branch_name = get_branch_name()
    project_name = isx_config['develop-project']
    message = get_pr_message().replace("'", "")
    project_info = isx_config['projects'][project_name]
    all_repository = list(project_info['sub-repository'])
    all_repository.append(project_info['repository'])
    for repository in all_repository:
        matcher = re.search(r"/([^/.]+)\.git", repository)
        project_name_meta = matcher.group(1)
        url = repository.replace('.git', '/pulls').replace('https://github.com', 'https://api.github.com/repos')
        headers = {
            "Accept": "application/vnd.github+json",
            "Authorization": "Bearer " + get_token(),
            "X-GitHub-Api-Version": "2022-11-28"
        }
        data = {
            "title": "#" + get_pr_num() + " " + message,
            "head": branch_name,
            "head_repo": get_account() + "/" + project_name_meta,
            "base": branch_name
        }
        response = requests.post(url, json=data, headers=headers)
        if response.status_code != 201:
            print(project_name_meta.ljust(18) + branch_name.ljust(18) + "【pr未提交】")
        else:
            print(project_name_meta.ljust(18) + branch_name.ljust(18) + "【pr已提交】")
