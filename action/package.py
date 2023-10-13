from config import get_config
import os


def package():
    isx_config = get_config()
    project_info = isx_config['projects'][isx_config['develop-project']]
    project_dir = project_info['dir']
    command = 'cd ' + project_dir + ' && ./gradlew package'
    os.system(command)