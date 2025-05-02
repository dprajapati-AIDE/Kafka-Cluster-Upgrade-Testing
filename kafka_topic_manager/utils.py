import subprocess

from logger import get_logger

logger = get_logger()


def get_broker_connection(zk, port):
    return f"{zk}:{port}"


def execute_command(cmd, dry_run=False):
    cmd_str = " ".join(cmd)

    if dry_run:
        logger.info(f"[DRY RUN] {cmd_str}")
        return

    try:
        output = subprocess.check_output(
            cmd, stderr=subprocess.STDOUT, text=True)
        return output
    except subprocess.CalledProcessError as e:
        logger.error(f"Command failed with exit code {e.returncode}")
        logger.error(f"Output: {e.output}")
        return e.returncode, e.output
