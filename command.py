import subprocess


def exec_command(command):
    print("执行命令：" + command)
    completed_process = subprocess.run(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    if completed_process.returncode == 0:
        if completed_process.stdout != '':
            print('命令输出：' + completed_process.stdout)
        if completed_process.stderr != '':
            print('命令输出：' + completed_process.stderr)
        print('执行命令成功')
    else:
        print(completed_process.stderr + "执行失败")
