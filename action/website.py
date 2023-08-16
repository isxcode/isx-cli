from config import get_current_project_path
import os


def website():
    command = 'cd ' + get_current_project_path() + ' && ./gradlew website'
    os.system(command)