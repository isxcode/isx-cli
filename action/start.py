from config import get_current_project_path
from config import check_current_project
import os


def start():
    check_current_project()
    command = 'cd ' + get_current_project_path() + ' && ./gradlew start'
    os.system(command)