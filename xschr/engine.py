import sys
import os
import time
import subprocess
from datetime import datetime
from .config import resolve_script_path

def run_sequence(experiments, config_path, python_cmd, log_root, fail_fast=False, dry_run=False):
    """
    The main execution loop. Iterates through experiments and runs, managing subprocesses and logs.
    """
    
    # 1. Setup Logging Directory
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    run_dir = os.path.join(log_root, f"run_{timestamp}")
    
    # 2. Calculate Stats for the Plan
    total_experiments = len(experiments)
    total_runs = sum(len(exp.get('runs', [])) for exp in experiments)
    
    print(f"\n[Plan]")
    print(f"  • Config:  {config_path}")
    print(f"  • Output:  {run_dir}")
    print(f"  • Task:    Running {total_runs} jobs across {total_experiments} experiments")

    # 3. Safety Confirmation
    if not dry_run:
        try:
            # Flush stdout to ensure prompt appears before input
            sys.stdout.write("\nPress ENTER to launch 🚀 (or Ctrl+C to abort)...")
            sys.stdout.flush()
            input()
            
            # Create directory only after confirmation
            os.makedirs(run_dir, exist_ok=True)
        except KeyboardInterrupt:
            print("\nAborted.")
            sys.exit(0)

    stats = {'success': 0, 'failed': 0}
    
    # 4. The Loop
    for exp_idx, exp in enumerate(experiments):
        exp_name = exp.get('name', f"exp_{exp_idx}")
        script_rel = exp['script']
        
        # Resolve script path relative to the config file location
        script_path = resolve_script_path(config_path, script_rel)
        
        # Verify script exists (unless dry run)
        if not os.path.exists(script_path) and not dry_run:
            print(f"\n\033[91m[Error]\033[0m Script not found: {script_path}")
            stats['failed'] += 1
            if fail_fast: return stats
            continue

        print(f"\n>> Experiment: {exp_name}")

        for i, run in enumerate(exp['runs']):
            run_id = i + 1
            args = run['args']
            
            # Construct command
            # Handle args as string ("--lr 0.1") or list (["--lr", "0.1"])
            arg_list = args.split() if isinstance(args, str) else args
            cmd = [python_cmd, script_path] + arg_list
            
            # Visual indicator
            print(f"   [{run_id}/{len(exp['runs'])}] {script_rel} {args}")

            if dry_run:
                continue

            # Log file setup
            safe_exp_name = exp_name.replace(" ", "_").replace("/", "-")
            log_filename = f"{safe_exp_name}_{run_id}.log"
            log_path = os.path.join(run_dir, log_filename)
            
            # Execute
            success = _execute_subprocess(cmd, log_path)
            
            if success:
                print("     \033[92m✓ Success\033[0m")
                stats['success'] += 1
            else:
                print("     \033[91m✗ Failed\033[0m")
                stats['failed'] += 1
                
                if fail_fast:
                    print("\n\033[93m[!] Fail-fast triggered. Stopping queue.\033[0m")
                    return stats

    return stats

def _execute_subprocess(cmd, log_path):
    """
    Handles the low-level subprocess creation, output streaming, and logging.
    Returns True if exit code is 0, False otherwise.
    """
    # Force unbuffered output so we see print statements immediately
    env = os.environ.copy()
    env['PYTHONUNBUFFERED'] = '1'

    try:
        with open(log_path, 'w') as f:
            # Write Header
            f.write(f"Cmd: {' '.join(cmd)}\n")
            f.write(f"Start: {datetime.now()}\n")
            f.write("-" * 40 + "\n")
            f.flush()

            # Start Process
            # stderr=subprocess.STDOUT merges errors into the main output stream
            process = subprocess.Popen(
                cmd,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                text=True,
                env=env,
                bufsize=1 # Line buffered
            )

            # Stream Output
            while True:
                line = process.stdout.readline()
                if not line and process.poll() is not None:
                    break
                if line:
                    # Print to console (indented for visual hierarchy)
                    sys.stdout.write(f"     | {line}")
                    sys.stdout.flush()
                    
                    # Write to log
                    f.write(line)
                    f.flush()

            # Write Footer
            return_code = process.poll()
            f.write("\n" + "-" * 40 + "\n")
            f.write(f"End: {datetime.now()}\n")
            f.write(f"Exit Code: {return_code}\n")
            
            return return_code == 0

    except Exception as e:
        print(f"     \033[91m[System Error] {e}\033[0m")
        return False
