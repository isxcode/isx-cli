from config import get_config
from config import save_config


def logout():
    isx_config = get_config()
    isx_config['user']['account'] = ""
    isx_config['user']['token'] = ""
    save_config(isx_config)
    print('退出成功')
