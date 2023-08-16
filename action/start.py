from config import get_current_project_path
import os


def start():
    command = 'cd ' + get_current_project_path() + ' && ./gradlew start'
    os.system(command)