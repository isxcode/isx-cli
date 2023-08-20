from config import get_config
from config import save_config
from action.clone import projects


def input_project_name(config):
    for project_meta in projects:
        project_dir = config[project_meta['project_name'].strip()]['dir']
        if project_dir != '':
            print(project_meta['project_index'] + '.' + project_meta['project_name'] + " : 已下载 ok")
        else:
            print(project_meta['project_index'] + '.' + project_meta['project_name'] + " : 未下载 x")
    project_number = input("请输入项目编号：")
    try:
        project_name = projects[int(project_number) - 1]['project_name'].strip()
        if config[project_name]['dir'] == '':
            print("该项目未下载，请选择其他项目")
            exit(0)
        return project_name
    except IndexError:
        print("请输入合法项目编号")
        exit(0)


def code():
    config = get_config()
    print("可选开发项目:")
    project = input_project_name(config)
    config['current-project'] = project
    save_config(config)
    print(project + " 切换成功")
