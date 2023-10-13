from action.branch import get_remote_branch
from action.branch import get_local_branch
from args import get_branch_num
from config import get_config


def get():
    isx_config = get_config()
    project_info = isx_config['projects'][isx_config['develop-project']]
    project_dir = project_info['dir']
    branch_num = get_branch_num()
    branch_name = get_local_branch(branch_num, project_dir)
    if branch_name != '':
        print("====================================")
        print("分支名称 ==> " + branch_name)
        print("====================================")
        exit(0)

    # 从个人仓库拉分支
    branch_name = get_remote_branch(branch_num, project_dir, 'origin')
    if branch_name != '':
        print("====================================")
        print("当前开发分支 ==> " + branch_name)
        print("====================================")
        exit(0)

    # 从upstream远程仓库拉分支
    branch_name = get_remote_branch(branch_num, project_dir, 'upstream')
    if branch_name != '':
        print("====================================")
        print("当前开发分支 ==> " + branch_name)
        print("====================================")
        exit(0)
    else:
        print("====================================")
        print("开发分支不存在")
        print("====================================")
        exit(0)
