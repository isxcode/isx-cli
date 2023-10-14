import re

from args import get_git_command
from config import get_config
from command import exec_command


def git():
    isx_config = get_config()
    project_name = isx_config['develop-project']
    project_info = isx_config['projects'][project_name]
    sub_repository = list(project_info['sub-repository'])
    project_dir = project_info['dir']
    git_command = get_git_command()
    command = "cd " + project_dir + " && " + git_command
    exec_command(command)
    for repository in sub_repository:
        matcher = re.search(r"/([^/.]+)\.git", repository)
        project_name = matcher.group(1)
        command = "cd " + project_dir + "/" + project_name + " && " + git_command
        exec_command(command)
    command = "cd " + project_dir
    exec_command(command)
