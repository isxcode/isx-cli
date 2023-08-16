from config import get_current_project_path
from config import check_current_project
import os


def vscode():
    check_current_project()
    command = 'cd ' + get_current_project_path() + ' && code ./'
    os.system(command)
