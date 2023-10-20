from config import init_config
from action.login import check_login
from args import get_action
from action.reset import reset
from action.login import login
from action.logout import logout
from action.clone import clone
from action.list import list_project
from action.choose import choose
from action.show import show
from action.idea import idea
from action.version import version
from action.install import install
from action.remove import remove
from action.format import format
from action.vscode import vscode
from action.clean import clean
from action.start import start
from action.package import package
from action.docker import docker
from action.deploy import deploy
from action.website import website
from action.web import web
from action.home import home
from action.branch import branch
from action.git import git
from action.get import get
from action.pr import pr
from action.frontend import frontend
from action.backend import backend
from args import check_min_args

if __name__ == '__main__':
    check_min_args()
    switch_action = {
        "reset": reset,
        "login": login,
        "logout": logout,
        "clone": clone,
        "list": list_project,
        "choose": choose,
        "show": show,
        "idea": idea,
        "vscode": vscode,
        "clean": clean,
        "start": start,
        "package": package,
        "docker": docker,
        "deploy": deploy,
        "website": website,
        "git": git,
        "get": get,
        "pr": pr,
        "branch": branch,
        "frontend": frontend,
        "backend": backend,
        "web": web,
        "home": home,
        "install": install,
        "remove": remove,
        "format": format,
        "version": version
    }
    if get_action() != 'login':
        check_login()
    do_action = switch_action.get(get_action())
    do_action()
