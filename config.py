import os
import yaml

config_path = os.path.expanduser("~/.isx/isx-config.yml")
config_temp = {
    "user": {
        "account": "",
        "token": ""
    },
    "develop-project": "",
    "projects": {
        "spark-yun": {
            "dir": "",
            "describe": "基于Spark开发批处理大数据平台",
            "repository": "https://github.com/isxcode/spark-yun.git",
            "sub-repository": [
                "https://github.com/isxcode/spark-yun-vip.git",
            ]
        },
        "flink-yun": {
            "dir": "",
            "describe": "基于Flink开发流处理大数据平台",
            "repository": "https://github.com/isxcode/flink-yun.git",
            "sub-repository": [
                "https://github.com/isxcode/flink-yun-vip.git",
            ]
        },
        "isx-cli": {
            "dir": "",
            "describe": "isxcode组织代码开发cli脚手架",
            "repository": "https://github.com/isxcode/isx-cli.git",
            "sub-repository": [
            ]
        },
        "isx-base": {
            "dir": "",
            "describe": "isxcode组织代码开发模版",
            "repository": "https://github.com/isxcode/isx-base.git",
            "sub-repository": [
                "https://github.com/isxcode/isx-base-vip.git",
            ]
        }
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


# 直接获取配置文件
def get_config():
    init_config()
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


def check_current_project(isx_config):
    if isx_config['develop-project'] == '':
        print("请先执行 【isx choose】命令选择开发项目")
        exit(0)
