import os
import yaml

config_path = os.path.expanduser("~/.isx/config.yml")
config_temp = {
    "current-project": "",
    "user": {
        "account": "",
        "token": ""
    },
    "spark-yun": {
        "repository": "https://github.com/isxcode/spark-yun.git",
        "dir": "",
        "has_private": True
    },
    "flink-yun": {
        "repository": "https://github.com/isxcode/flink-yun.git",
        "dir": "",
        "has_private": True
    },
    "isx-cli": {
        "repository": "https://github.com/isxcode/isx-cli.git",
        "dir": "",
        "has_private": False
    },
    "isx-app": {
        "repository": "https://github.com/isxcode/isx-app.git",
        "dir": "",
        "has_private": False
    }
}


def print_config():
    print("可选择开发项目：")
    print("- spark-yun")
    print("- flink-yun")
    print("- isx-cli")
    print("- isx-app")


def init_config():
    directory = os.path.dirname(config_path)
    if not os.path.exists(directory):
        os.makedirs(directory)
    if not os.path.exists(config_path):
        with open(config_path, "x") as f:
            yaml.dump(config_temp, f)


def get_config():
    config_file = open(config_path, mode='r', encoding='utf-8')
    return yaml.safe_load(config_file)


def save_config(config_data):
    with open(config_path, 'w') as f:
        yaml.dump(config_data, f)


def clear_config():
    os.remove(config_path)


def get_current_project():
    config = get_config()
    return config['current-project']


def get_account():
    config = get_config()
    return config['user']['account']


def get_token():
    config = get_config()
    return config['user']['token']


def get_current_project_dir():
    config = get_config()
    return config[config['current-project']]['dir']


def get_current_project_path():
    config = get_config()
    return config[config['current-project']]['dir'] + "/" + config['current-project']


def get_current_project_vip_path():
    config = get_config()
    return config[config['current-project']]['dir'] + "/" + config['current-project'] + "/" + config[
        'current-project'] + "-vip"


def check_current_project():
    config = get_config()
    if config['current-project'] == '':
        print("请使用 isx choose 命令，选择开发项目")
        exit(0)
    if config[config['current-project']]['dir'] == '':
        print("请使用 isx clone 命令，下载开发项目")
        exit(0)
