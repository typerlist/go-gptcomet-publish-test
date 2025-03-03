import os
import shlex
import subprocess
import sys
import platform

def find_gptcomet_binary():
    platform_name = sys.platform
    
    machine = platform.machine().lower()
    if machine in ('arm64', 'aarch64', 'arm'):
        arch = 'arm64'
    elif machine in ('x86_64', 'amd64', 'x64', 'i386', 'x86'):
        arch = 'amd64'
    else:
        raise OSError(f"Unsupported architecture: {machine}")
    
    if platform_name == "win32":
        binary_name = f"gptcomet_{arch}.exe"
    elif platform_name == "darwin":
        binary_name = f"gptcomet_{arch}_mac"
    elif platform_name == "linux":
        binary_name = f"gptcomet_{arch}_linux"
    else:
        raise OSError(f"Unsupported platform: {platform_name}")
    
    binary_path = os.path.join(os.path.dirname(__file__), "bin", binary_name)
    if not os.path.isfile(binary_path):
        msg = f"gptcomet binary not found for {platform_name}-{arch}, please open an issue on github, thanks."
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
