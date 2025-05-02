import subprocess

from logger import get_logger

logger = get_logger()


def execute_command(cmd, dry_run=False):
    cmd_str = " ".join(cmd)

    if dry_run:
        logger.info(f"[DRY RUN] {cmd_str}")
        return 0, ""

    try:
        output = subprocess.check_output(
            cmd, stderr=subprocess.STDOUT, universal_newlines=True)
        return 0, output
    except subprocess.CalledProcessError as e:
        logger.error(f"Command failed with exit code {e.returncode}")
        logger.error(f"Output: {e.output}")
        return e.returncode, e.output
