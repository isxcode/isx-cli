from config import get_current_project_path
import os


def clean():
    command = 'cd ' + get_current_project_path() + ' && ./gradlew clean'
    os.system(command)