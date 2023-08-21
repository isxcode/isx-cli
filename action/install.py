import os

from command import exec_command
from config import check_current_project
from config import get_current_project
from config import get_current_project_path


def install():
    check_current_project()
    project = get_current_project()
    project_path = get_current_project_path()
    download_path = "~/Downloads"
    if project == 'spark-yun':
        spark_file_name = "spark-3.4.0-bin-hadoop3.tgz"
        spark_download_url = "https://archive.apache.org/dist/spark/spark-3.4.0/spark-3.4.0-bin-hadoop3.tgz"
        if not os.path.exists(os.path.expanduser(download_path + "/" + spark_file_name)):
            exec_command("cd " + download_path + " && wget " + spark_download_url)
        if not os.path.exists(project_path + "/spark-yun-dist/spark-min"):
            spark_min_dir = project_path + "/spark-yun-dist/spark-min"
            exec_command(
                "mkdir -p " + spark_min_dir + "  && tar vzxf " + download_path + "/" + spark_file_name + " --strip-components=1 -C "+spark_min_dir)
    print("依赖安装完成")