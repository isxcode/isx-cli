from config import get_current_project_path
import os


def deploy():
    command = 'cd ' + get_current_project_path() + ' && ./gradlew deploy'
    os.system(command)