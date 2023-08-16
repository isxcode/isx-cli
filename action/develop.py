from args import get_project
from config import get_config
from config import save_config


def develop():
    config = get_config()
    project = get_project()
    try:
        config[project]
    except Exception:
        print('请输入已安装项目')
        exit(0)
    try:
        config['current-project'] = project
    except Exception:
        print('请输入已安装项目')
        exit(0)
    save_config(config)
    print("切换项目成功：" + project)