import sys
import os
import time
from .cli import get_parser, Environment
from .config import load_and_validate
from .system import print_system_status
from .engine import run_sequence

def main():
    """
    The core logic driver.
    """
    # 1. Initialize Environment
    env = Environment()

    # 2. Parse Arguments
    parser = get_parser(env)
    args = parser.parse_args()

    # 3. System Check (Skip if dry-run to reduce noise, or keep it if you prefer)
    if not args.dry_run:
        try:
            print_system_status()
        except Exception:
            # Don't crash if nvidia-smi fails, just ignore
            pass

    # 4. Load Config
    try:
        config_data, config_abs_path = load_and_validate(args.path)
    except Exception as e:
        env.log_error(f"Configuration Failed: {e}")
        return 1

    # 5. Extract Global Settings
    conf_global = config_data.get('config', {})
    # Default to current python interpreter if not specified
    python_cmd = conf_global.get('python_cmd', sys.executable)
    log_dir = conf_global.get('log_dir', 'logs')

    # 6. Run Engine
    # We pass the parsed arguments to the engine
    try:
        stats = run_sequence(
            experiments=config_data['experiments'],
            config_path=config_abs_path,
            python_cmd=python_cmd,
            log_root=log_dir,
            fail_fast=args.fail_fast,
            dry_run=args.dry_run
        )
    except KeyboardInterrupt:
        env.log_error("Execution interrupted by user.")
        return 130
    except Exception as e:
        env.log_error(f"Engine Crash: {e}")
        if args.debug:
            import traceback
            traceback.print_exc()
        return 1

    # 7. Final Summary
    # The engine handles per-run printing, we just summarize the totals.
    if not args.dry_run:
        print(f"\n[Final Summary]")
        if stats['failed'] == 0:
            print(f"\033[1;32m✓ All {stats['success']} runs completed successfully.\033[0m")
            return 0
        else:
            print(f"\033[1;31m✗ Completed: {stats['success']} | Failed: {stats['failed']}\033[0m")
            return 1
    
    return 0
