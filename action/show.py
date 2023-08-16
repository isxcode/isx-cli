from config import get_config
from config import check_current_project


def show():
    check_current_project()
    config = get_config()
    print('当前项目:' + config['current-project'])
    if config['current-project'] == '':
        print("请使用isx develop命令，选择开发项目")
        exit(0)
    print('项目目录:' + config[config['current-project']]['dir']+"/"+config['current-project'])