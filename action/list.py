from config import get_config


def list_project():
    isx_config = get_config()
    projects = isx_config['projects']
    develop_project = isx_config['develop-project']
    for index, project_name in enumerate(projects):
        project_info = '[' + str(index) + '] ' + project_name.ljust(12)
        if projects[project_name]['dir'] != '':
            home_info = '【 已下载 ' + projects[project_name]['dir'] + ' 】'
            project_info = project_info + home_info.ljust(60)
        else:
            project_info = project_info + '【 待下载'.ljust(60) + '】'
        if develop_project == project_name:
            project_info = project_info + ' 【 Coding 】'
        else:
            project_info = project_info + ' 【 Wait   】'
        print(project_info)
