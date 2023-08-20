from args import get_git_command
from config import check_current_project
from config import get_current_project_vip_path
from config import get_current_project_path
from config import get_current_project

import os


def install():
    check_current_project()
    project = get_current_project()
    project_path = get_current_project_path()
    # 判断code 是否为spark-yun
    # 使用命令下载spark二进制文件
    # 并解压到项目中
