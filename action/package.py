from config import get_current_project_path
import os


def package():
    command = 'cd ' + get_current_project_path() + ' && ./gradlew package'
    os.system(command)