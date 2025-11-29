import ctypes
from ctypes.util import find_library

# --- Types ---
CUresult = ctypes.c_uint32
CUdevice = ctypes.c_int32

# --- CUDA Error Codes ---
CUDA_SUCCESS = 0

# --- Library Loading ---
def load_cuda_driver():
    """Load CUDA driver library. Returns dll or None."""
    try:
        lib_path = find_library('cuda')
        if lib_path:
            return ctypes.CDLL(lib_path)
    except:
        pass
    return None

# --- API Function Setup ---
def setup_cuda_functions(dll):
    """Setup CUDA API function signatures. Returns dict of functions or None."""
    if not dll:
        return None
    
    funcs = {}
    
    try:
        # cuInit
        funcs['cuInit'] = dll.cuInit
        funcs['cuInit'].restype = CUresult
        funcs['cuInit'].argtypes = [ctypes.c_uint32]
        
        # cuDeviceGetCount
        funcs['cuDeviceGetCount'] = dll.cuDeviceGetCount
        funcs['cuDeviceGetCount'].restype = CUresult
        funcs['cuDeviceGetCount'].argtypes = [ctypes.POINTER(ctypes.c_int32)]
        
        # cuDeviceGet
        funcs['cuDeviceGet'] = dll.cuDeviceGet
        funcs['cuDeviceGet'].restype = CUresult
        funcs['cuDeviceGet'].argtypes = [ctypes.POINTER(CUdevice), ctypes.c_int32]
        
        # cuDeviceGetName
        funcs['cuDeviceGetName'] = dll.cuDeviceGetName
        funcs['cuDeviceGetName'].restype = CUresult
        funcs['cuDeviceGetName'].argtypes = [ctypes.POINTER(ctypes.c_char), ctypes.c_int32, CUdevice]
        
        return funcs
    except AttributeError:
        return None

# --- Detection Functions ---
def detect_nvidia_gpus():
    """
    Detect NVIDIA GPUs using CUDA driver API.
    Returns list of dicts with GPU info, or empty list if none found.
    
    Returns:
        [{'id': 0, 'name': 'RTX 4090', 'success': True}, ...]
    """
    dll = load_cuda_driver()
    if not dll:
        return []
    
    funcs = setup_cuda_functions(dll)
    if not funcs:
        return []
    
    # Initialize CUDA
    result = funcs['cuInit'](0)
    if result != CUDA_SUCCESS:
        return []
    
    # Get device count
    count = ctypes.c_int32()
    result = funcs['cuDeviceGetCount'](ctypes.byref(count))
    if result != CUDA_SUCCESS or count.value == 0:
        return []
    
    # Get info for each GPU
    gpus = []
    for i in range(count.value):
        device = CUdevice()
        result = funcs['cuDeviceGet'](ctypes.byref(device), i)
        if result != CUDA_SUCCESS:
            continue
        
        name_buffer = ctypes.create_string_buffer(256)
        result = funcs['cuDeviceGetName'](name_buffer, 256, device)
        
        gpu_info = {
            'id': i,
            'name': name_buffer.value.decode() if result == CUDA_SUCCESS else 'Unknown',
            'success': result == CUDA_SUCCESS
        }
        gpus.append(gpu_info)
    
    return gpus

def get_gpu_count():
    """Quick check: how many GPUs? Returns 0 if none or error."""
    gpus = detect_nvidia_gpus()
    return len(gpus)

def has_gpus():
    """Quick check: any GPUs available? Returns bool."""
    return get_gpu_count() > 0

# --- Usage Example ---
if __name__ == "__main__":
    gpus = detect_nvidia_gpus()
    
    if gpus:
        print(f"Found {len(gpus)} NVIDIA GPU(s):")
        for gpu in gpus:
            print(f"  [{gpu['id']}] {gpu['name']}")
    else:
        print("No NVIDIA GPUs detected")
