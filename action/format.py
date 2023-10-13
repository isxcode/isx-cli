from config import get_config
import os


def format():
    isx_config = get_config()
    project_info = isx_config['projects'][isx_config['develop-project']]
    project_dir = project_info['dir']
    command = 'cd ' + project_dir + ' && ./gradlew format'
    os.system(command)
