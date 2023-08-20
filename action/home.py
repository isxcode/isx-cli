import os

from config import check_current_project
from config import get_current_project_path


def home():
    check_current_project()
    command = "cd " + get_current_project_path()
    os.system(command)
