import shutil
import subprocess
import sys
import platform

def get_gpu_info():
    """
    Returns a list of detected GPUs or None.
    """
    # 1. Check for NVIDIA (Linux/Windows)
    if shutil.which("nvidia-smi"):
        try:
            cmd = ["nvidia-smi", "--query-gpu=name,memory.total", "--format=csv,noheader"]
            result = subprocess.check_output(cmd, encoding='utf-8')
            return [line.strip() for line in result.strip().split('\n') if line.strip()]
        except Exception:
            pass
            
    # 2. Check for Mac (MPS/Metal)
    if platform.system() == "Darwin" and platform.machine() == "arm64":
        return ["Apple Silicon (Metal/MPS)"]

    return None

def print_system_status():
    """Prints current hardware status."""
    gpus = get_gpu_info()
    
    # Force flush to ensure it prints before the next section
    print("  • Platform:      " + platform.system() + " " + platform.machine())
    
    if gpus:
        print(f"  • Accelerator:   Detected {len(gpus)} device(s):")
        for idx, gpu in enumerate(gpus):
            print(f"    - [{idx}] {gpu}")
    else:
        print("  • Accelerator:   None (Running on CPU)")
    
    sys.stdout.flush()
