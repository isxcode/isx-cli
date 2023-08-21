import subprocess


def exec_command(command):
    print("执行命令：" + command)
    completed_process = subprocess.run(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    if completed_process.returncode == 0:
        print(completed_process.stderr + "执行成功")
    else:
        print(completed_process.stderr + "执行失败")
