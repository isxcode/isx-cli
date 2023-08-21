from action.login import check_login
from config import get_config
from config import save_config
import os
import subprocess

projects = [
    {
        "project_index": "1",
        "project_name": "spark-yun",
        "project_describe": "基于Spark开发批处理大数据平台"
    }, {
        "project_index": "2",
        "project_name": "flink-yun",
        "project_describe": "基于Flink开发流处理大数据平台"
    }, {
        "project_index": "3",
        "project_name": "isx-app  ",
        "project_describe": "isxcode开源组织项目模版"
    }, {
        "project_index": "4",
        "project_name": "isx-cli  ",
        "project_describe": "isxcode开源组织项目脚手架"
    }
]


def input_project_name():
    for project_meta in projects:
        print(project_meta['project_index'] + '.' + project_meta['project_name'] + " : " + project_meta[
            'project_describe'])
    project_number = input("请输入项目编号：")
    try:
        return projects[int(project_number) - 1]['project_name'].strip()
    except IndexError:
        print("请输入合法项目编号")
        exit(0)


def input_project_dir(config, project):
    project_dir = input("请输入下载目录(全路径)：")
    if project_dir.endswith("/"):
        project_dir = project_dir[:-1]
    if not os.path.exists(project_dir):
        print("请输入合法目录")
        exit(0)
    config[project]['dir'] = project_dir
    save_config(config)
    return project_dir


def clone_github_code(repository, account, project_dir, project_name):
    command = 'cd ' + project_dir + ' && git clone ' + repository.replace("isxcode", account)
    print("执行下载命令：" + command)
    completed_process = subprocess.run(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    os.system('cd ' + project_dir + '/' + project_name + ' && git remote add upstream ' + repository)
    if completed_process.returncode == 0:
        print(completed_process.stdout + "下载成功")
    else:
        print(completed_process.stderr + "下载失败")


def clone_code(config, project, project_dir):
    # 判断项目是否存在
    if os.path.exists(project_dir + '/' + project):
        print("目录已存在，请删除重试")
        exit(0)
    # 拉取开源代码
    clone_github_code(config[project]["repository"], config["user"]["account"], project_dir, project)
    # 拉取闭源代码
    if config[project]["has_private"] is True:
        clone_github_code(config[project]["repository"].replace(project, project + "-vip"), config["user"]["account"],
                          project_dir + "/" + project, project + "-vip")


def clone():
    # 检查用户是否登录
    check_login()
    # 获取配置文件
    config = get_config()
    # 输入项目
    print("可选下载项目: ")
    project = input_project_name()
    # 输入下载目录
    project_dir = input_project_dir(config, project)
    # 拉取代码
    clone_code(config, project, project_dir)
