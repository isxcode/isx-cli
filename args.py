import sys


def get_action():
    return sys.argv[1]


def get_git_command():
    return ' '.join(sys.argv[1:])


def get_pr_num():
    return sys.argv[2]


def get_branch_num():
    if len(sys.argv) < 3:
        print("请输入正确的数字")
        exit(0)
    branch_num = sys.argv[2]
    if branch_num.isdigit() and int(branch_num) > 0:
        return branch_num
    else:
        print("请输入正确的数字")
        exit(0)


def get_pr_message():
    return sys.argv[3]


def check_min_args():
    if len(sys.argv) < 2:
        command_length = 40
        print("请输入完整命令，参考如下：")
        print("isx reset".ljust(command_length) + " -- 重置配置")
        print("isx login".ljust(command_length) + " -- 用户登录")
        print("isx logout".ljust(command_length) + " -- 用户退出")
        print("isx clone".ljust(command_length) + " -- 下载项目")
        print("isx list".ljust(command_length) + " -- 查看本地项目列表")
        print("isx choose".ljust(command_length) + " -- 选择开发项目")
        print("isx branch <issue_number>".ljust(command_length) + " -- 切开发分支")
        print("isx install".ljust(command_length) + " -- 安装项目依赖")
        print("isx start".ljust(command_length) + " -- 启动项目")
        print("isx package".ljust(command_length) + " -- 打包项目")
        print("isx format".ljust(command_length) + " -- 格式化代码")
        print("isx get <issue_number>".ljust(command_length) + " -- 获取需求分支名称")
        print("isx <git_command>".ljust(command_length) + " -- 在所有模块中执行git命令")
        print("isx pr <issue_number> '<pr_command>'".ljust(command_length) + " -- 提交所有模块的pr")
        exit(0)
