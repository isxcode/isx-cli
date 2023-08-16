from action.login import check_login
from config import get_config
from config import save_config
import os


def input_project_name(config):
    print("可选项目: 【spark-yun】【flink-yun】【isx-app】【isx-cli】")
    project = input("请输入项目：")
    try:
        config[project]
    except BaseException:
        print("项目不存在")
        exit(0)
    return project


def input_project_dir(config, project):
    project_dir = input("请输入下载目录(全路径)：")
    if not os.path.exists(project_dir):
        print("目录不存在")
        exit(0)
    config[project]['dir'] = project_dir
    save_config(config)
    return project_dir


def clone_github_code(repository, account, project_dir):
    command = 'cd ' + project_dir + '&& git clone ' + repository.replace("isxcode", account)
    return os.system(command)


def clone_code(config, project, project_dir):
    # 判断项目是否存在
    if os.path.exists(project_dir + '/' + project):
        print("项目已存在")
        exit(0)
    # 拉取开源代码
    clone_github_code(config[project]["repository"], config["user"]["account"], project_dir)
    # 拉取闭源代码
    clone_github_code(config[project]["repository"].replace(project, project + "-vip"), config["user"]["account"],
                      project_dir + "/" + project)


def clone():
    # 检查用户是否登录
    check_login()
    # 获取配置文件
    config = get_config()
    # 输入项目
    project = input_project_name(config)
    # 输入下载目录
    project_dir = input_project_dir(config, project)
    # 拉取代码
    clone_code(config, project, project_dir)
