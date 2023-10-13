import subprocess


def exec_command(command):
    print("执行命令【 " + command + " 】")
    completed_process = subprocess.run(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    if completed_process.returncode == 0:
        if completed_process.stdout != '':
            print('stdout ==> ' + completed_process.stdout[:-1])
        if completed_process.stderr != '':
            print('stderr ==> ' + completed_process.stderr[:-1])
    else:
        print(completed_process.stderr + "执行失败")
