import os
import shlex
import subprocess
import sys


def find_gptcomet_binary():
    binary_name = "gptcomet.exe" if sys.platform == "win32" else "gptcomet"
    binary_path = os.path.join(os.path.dirname(__file__), "bin", binary_name)
    if not os.path.isfile(binary_path):
        msg = "gptcomet binary not found, this version of gptcomet may be incomplete, please open an issue on github, thanks."
        raise FileNotFoundError(msg)
    return binary_path


def main():
    binary = find_gptcomet_binary()
    args = [binary] + sys.argv[1:]
    if sys.platform == "win32":
        # no need for windows
        subprocess.run(args, check=True)  # noqa: S602
    else:
        command = shlex.join(args)
        subprocess.run(command, shell=True, check=True)  # noqa: S602


if __name__ == "__main__":
    main()
