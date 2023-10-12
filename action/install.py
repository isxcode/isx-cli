import platform

from command import exec_command
from config import get_config


def install():
    isx_config = get_config()
    if isx_config['develop-project'] == '':
        print("请先执行 【isx choose】命令选择开发项目，在执行 【isx install】")
        exit(0)
    project_dir = isx_config['projects'][isx_config['develop-project']]['dir']
    system_name = platform.system()
    if system_name == "Darwin":
        exec_command("cd " + project_dir + " && bash install.sh")
    elif system_name == "Windows":
        exec_command("cd " + project_dir + " && install.bat")
    else:
        print("本地系统暂不支持")