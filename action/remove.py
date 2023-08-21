from command import exec_command
from config import check_current_project
from config import get_config
from config import get_current_project
from config import get_current_project_path
from config import save_config


def remove():
    check_current_project()
    config = get_config()
    current_project = get_current_project()
    if config[current_project]['dir'] == '':
        print("当前项目已删除，请使用isx choose重新选择开发项目")
    exec_command("rm -rf " + get_current_project_path())
    config['current-project'] = ''
    config[current_project]['dir'] = ''
    save_config(config)
    print(current_project + "删除成功，请使用isx choose重新选择开发项目")