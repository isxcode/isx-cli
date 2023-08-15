import sys
import os
import yaml

config_path = os.path.expanduser("~/.isx/config.yml")


# 检查配置文件并创建
def check_config():
    directory = os.path.dirname(config_path)
    if not os.path.exists(directory):
        os.makedirs(directory)
    if not os.path.exists(config_path):
        with open(config_path, "x") as f:
            config = {
                "current-project": "",
                "user": {
                    "account": "",
                    "token": ""
                },
                "spark-yun": {
                    "upstream": "https://github.com/isxcode/spark-yun.git",
                    "vip-upstream": "https://github.com/isxcode/spark-yun-vip.git",
                    "dir": ""
                },
                "flink-yun": {
                    "upstream": "https://github.com/isxcode/flink-yun.git",
                    "vip-upstream": "https://github.com/isxcode/flink-yun-vip.git",
                    "dir": ""
                },
                "isx-cli": {
                    "upstream": "https://github.com/isxcode/isx-cli.git",
                    "vip-upstream": "https://github.com/isxcode/isx-cli-vip.git",
                    "dir": ""
                },
                "isx-app": {
                    "upstream": "https://github.com/isxcode/isx-app.git",
                    "vip-upstream": "https://github.com/isxcode/isx-app-vip.git",
                    "dir": ""
                }
            }
            yaml.dump(config, f)


def get_config():
    config_file = open(config_path, mode='r', encoding='utf-8')
    return yaml.safe_load(config_file)


def save_config(config_data):
    with open(config_path, 'w') as f:
        yaml.dump(config_data, f)


def action():
    return sys.argv[1]


def get_project():
    return sys.argv[2]


def check_login():
    config = get_config()
    if config['user']['account'] == '' and config['user']['token'] == '':
        print('请先登录')
        login()
        exit(0)


# 用户登录
def login():
    config = get_config()
    if config['user']['account'] != '' and config['user']['token'] != '':
        print('请先退出登录')
        exit(0)
    if config['user']['account'] == '':
        account = input("请输入Github的Account：")
        config['user']['account'] = account
        save_config(config)
    if config['user']['token'] == '':
        print("快捷链接：https://github.com/settings/tokens")
        token = input("请输入Github的Token：")
        config['user']['token'] = token
        save_config(config)
    print('登录成功')


# 用户退出
def logout():
    config = get_config()
    config['user'] = {'account': '', 'token': ''}
    save_config(config)
    print('退出登录')





def clone():
    check_login()
    config = get_config()
    print("当前可选项目: 【spark-yun】【flink-yun】【isx-app】【isx-cli】")
    project = input("请输入项目：")
    while project != 'spark-yun' and project != 'flink-yun' and project != 'isx-app' and project != 'isx-cli':
        project = input("请输入正确项目名称：")
    project_path = input("请输入下载目录(全路径)：")
    while not os.path.exists(project_path):
        if not os.path.exists(project_path):
            project_path = input("目录不存在，请重新输入：")
    config[project]['dir'] = project_path + '/' + project
    config['current-project'] = project
    save_config(config)
    if not os.path.exists(project_path + '/' + project):
        upstream_git = config[project]['upstream'].replace('isxcode', config['user']['account'])
        command = 'cd ' + project_path + '&& git clone ' + upstream_git
        exit_status = os.system(command)
        if exit_status == 0:
            vip_upstream_git = config[project]['vip-upstream'].replace('isxcode', config['user']['account'])
            command = 'cd ' + project_path + '/' + project + ' && git clone ' + vip_upstream_git
            exit_status = os.system(command)
            if exit_status == 0:
                print('下载成功')
            else:
                print('下载失败')
        else:
            print('下载失败')
    else:
        print('该文件夹已下载')


def code():
    config = get_config()
    command = 'cd ' + config[config['current-project']]['dir'] + ' && code ./'
    os.system(command)


def idea():
    config = get_config()
    command = 'cd ' + config[config['current-project']]['dir'] + ' && idea ./'
    os.system(command)


# 启动项目
def start():
    config = get_config()
    command = 'cd ' + config[config['current-project']]['dir'] + ' && ./gradlew start'
    os.system(command)


def package():
    config = get_config()
    command = 'cd ' + config[config['current-project']]['dir'] + ' && ./gradlew package'
    os.system(command)


def website():
    config = get_config()
    command = 'cd ' + config[config['current-project']]['dir'] + ' && ./gradlew website'
    os.system(command)


def deploy():
    config = get_config()
    command = 'cd ' + config[config['current-project']]['dir'] + ' && ./gradlew deploy'
    os.system(command)


def clean():
    config = get_config()
    command = 'cd ' + config[config['current-project']]['dir'] + ' && ./gradlew clean'
    os.system(command)


def docker():
    config = get_config()
    command = 'cd ' + config[config['current-project']]['dir'] + ' && ./gradlew docker'
    os.system(command)


def show():
    config = get_config()
    print('当前项目:' + config['current-project'])
    print('项目目录:' + config[config['current-project']]['dir'])


def change():
    config = get_config()
    try:
        project = get_project()
    except IndexError:
        print('请输入已安装项目')
        exit(0)
    if project != 'spark-yun' and project != 'flink-yun' and project != 'isx-app' and project != 'isx-cli':
        print('请输入已安装项目')
        exit(0)
    config['current-project'] = project
    save_config(config)
    print("切换项目成功：" + project)


# 检测token是否合法
# def check_token()


# def logout():

# 一个一个函数和行为

#

# 将数据同步到 .isx-cli/config.yml 中
#

# def login(name):
#     print(f'Hi, {name}')  # Press ⌘F8 to toggle the breakpoint.

#
# def logout():
#     # 判断是否存在登录
#     print(f'Hi, {name}')  # Press ⌘F8 to toggle the breakpoint.
#

# Press the green button in the gutter to run the script.
if __name__ == '__main__':
    check_config()
    switch_action = {
        "login": login,
        "logout": logout,
        "clone": clone,
        "start": start,
        "code": code,
        "idea": idea,
        "show": show,
        "change": change,
        "package": package,
        "clean": clean,
        "docker": docker,
        "deploy": deploy,
        "website": website
    }
    do = switch_action.get(action())
    do()
    # 获取命令参数，判断参数合法性
    # 如果有文件则，实时读取文件，写入文件
    # 判断文件里是否
    # 判断当前系统是否有git命令

# See PyCharm help at https://www.jetbrains.com/help/pycharm/
