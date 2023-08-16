from config import check_current_project
from config import get_current_project_path
import os


def idea():
    check_current_project()
    command = 'cd ' + get_current_project_path() + ' && idea ./'
    os.system(command)
