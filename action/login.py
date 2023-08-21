from config import get_config
from config import save_config


def check_login():
    config = get_config()
    if config['user']['account'] == '' or config['user']['token'] == '':
        print('请先登录')
        login()
        exit(0)


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
    print('登录成功,欢迎使用isx开发工具')