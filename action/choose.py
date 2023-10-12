from config import get_config
from config import save_config
from action.list import list_project


def input_project_name(isx_config):
    list_project()
    project_number = input("请输入开发项目编号：")
    if project_number != '' and 0 <= int(project_number) < len(isx_config["projects"]) and project_number.isdigit():
        project_name = list(isx_config["projects"])[int(project_number)]
        if isx_config["projects"][project_name]['dir'] != '':
            return project_name
        else:
            print("项目未下载，请先执行【isx clone】命令")
            exit(0)
    else:
        print("选择编号异常，请重新执行【isx choose】命令")
        exit(0)


def choose():
    isx_config = get_config()
    project_name = input_project_name(isx_config)
    isx_config['develop-project'] = project_name

    save_config(isx_config)
    print(project_name + " 切换成功")
