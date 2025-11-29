import platform
from .cuda_devices import detect_nvidia_gpus

def print_system_status():
    """Print current hardware status."""
    print("  • Platform:      " + platform.system() + " " + platform.machine())
    
    gpus = detect_nvidia_gpus()
    
    if gpus:
        print(f"  • Accelerator:   Detected {len(gpus)} NVIDIA GPU(s):")
        for gpu in gpus:
            print(f"    - [{gpu['id']}] {gpu['name']}")
    else:
        print("  • Accelerator:   None (Running on CPU)")
