def logout():
    config = get_config()
    config['user'] = {'account': '', 'token': ''}
    save_config(config)
    print('退出成功!!!')