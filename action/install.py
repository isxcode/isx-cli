import platform

from command import exec_command
from config import get_config
from config import check_current_project


def install():
    isx_config = get_config()
    check_current_project(isx_config)
    project_dir = isx_config['projects'][isx_config['develop-project']]['dir']
    system_name = platform.system()
    if system_name == "Darwin":
        exec_command("cd " + project_dir + " && bash install.sh")
    elif system_name == "Windows":
        exec_command("cd " + project_dir + " && install.bat")
    else:
        print("本地系统暂不支持")