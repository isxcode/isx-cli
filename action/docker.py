from config import get_current_project_path
import os


def docker():
    command = 'cd ' + get_current_project_path() + ' && ./gradlew docker'
    os.system(command)