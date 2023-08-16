from args import get_git_command
from config import check_current_project
from config import get_current_project_vip_path
from config import get_current_project_path
import os


def git():
    check_current_project()
    git_command = get_git_command()
    command = "cd " + get_current_project_path() + " && " + git_command
    os.system(command)
    command = "cd " + get_current_project_vip_path() + " && " + git_command
    os.system(command)
    command = "cd " + get_current_project_path()
    os.system(command)