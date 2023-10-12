import requests

from config import get_config
from config import save_config


# 用户token登录
def login():
    isx_config = get_config()
    if isx_config['user']['account'] != '' and isx_config['user']['token'] != '':
        print('用户已登录，若要登录新账号或者配置新的token，请先使用【isx logout】退出登录')
        exit(0)
    else:
        isx_config['user']['account'] = input("请输入Github的Account：")
        print("快捷链接：https://github.com/settings/tokens")
        isx_config['user']['token'] = input("请输入Github的Token：")
        check_github_token(isx_config)
        save_config(isx_config)


# 检查github的token是否合法
def check_github_token(isx_config):
    headers = {
        "Accept": "application/vnd.github+json",
        "Authorization": "Bearer " + isx_config['user']['token'],
        "X-GitHub-Api-Version": "2022-11-28"
    }
    response = requests.get("https://api.github.com/octocat", headers=headers)
    if response.status_code == 200:
        print("登录成功,欢迎使用isx开发工具")
    else:
        print("无法验证token合法性，登录失败")
        exit(0)


# 检查用户是否登录
def check_login():
    config = get_config()
    if config['user']['account'] == '' or config['user']['token'] == '':
        print('请先登录')
        login()
