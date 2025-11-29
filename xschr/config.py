import os
import sys
import json
import yaml 

def load_and_validate(config_path):
    """
    Loads a YAML configuration file and validates the schema.
    Returns: (config_data: dict, absolute_path: str)
    """
    abs_path = os.path.abspath(config_path)
    if not os.path.exists(abs_path):
        raise FileNotFoundError(f"Config file not found at {abs_path}")

    # Detect format by extension
    ext = os.path.splitext(abs_path)[1].lower()
    
    try:
        with open(abs_path, 'r') as f:
            if ext in ['.yaml', '.yml']:
                # Safe load prevents arbitrary code execution in YAML
                data = yaml.safe_load(f)
            else:
                # Let's try YAML, then fail.
                try:
                    f.seek(0)
                    data = yaml.safe_load(f)
                except:
                    raise ValueError(f"Unsupported file format: {ext}. Please use .yaml file")
                    
    except Exception as e:
        raise ValueError(f"Error parsing {ext} file: {e}")

    # Basic Validation
    if not data:
        raise ValueError("Config file is empty.")
        
    if 'experiments' not in data:
        raise ValueError("Config missing required 'experiments' list.")

    if not isinstance(data['experiments'], list):
        raise ValueError("'experiments' must be a list.")

    return data, abs_path

def resolve_script_path(config_path_abs, script_relative):
    """
    Resolves a script path relative to the location of the config file.
    This ensures xschr works regardless of where you run it from.
    """
    config_dir = os.path.dirname(config_path_abs)
    return os.path.normpath(os.path.join(config_dir, script_relative))
