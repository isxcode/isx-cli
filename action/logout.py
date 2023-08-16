from config import get_config
from config import save_config


def logout():
    config = get_config()
    config['user'] = {'account': '', 'token': ''}
    save_config(config)
    print('退出成功!!!')