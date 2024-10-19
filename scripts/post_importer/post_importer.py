# This script is used to import posts from a directory to the content/posts directory
# The first argument is the input directory
# The second argument is the output directory
# It supports markdown posts created with hugo

import os
import sys
import shutil
import re
from datetime import datetime
import pytz

def parse_date(content):
    date_match = re.search(r'^date:\s*"?(.+?)"?\s*$', content, re.MULTILINE)
    if not date_match:
        raise ValueError("Date not found in file")
    
    date_str = date_match.group(1).strip()
    
    # Try parsing with different formats
    formats = [
        "%Y-%m-%dT%H:%M:%S%z",  # Format 1: "2006-01-02T15:04:05Z07:00"
        "%B %d, %Y %H:%M",      # Format 2: "January 2, 2006 15:04"
        "%Y-%m-%d"              # Format 3: "2006-01-02"
    ]
    
    for fmt in formats:
        try:
            if fmt == "%Y-%m-%dT%H:%M:%S%z":
                # Special handling for ISO format with timezone
                date = datetime.strptime(date_str, fmt)
                return date.astimezone(pytz.UTC)
            elif fmt == "%B %d, %Y %H:%M":
                # Assume local timezone for this format
                date = datetime.strptime(date_str, fmt)
                local_tz = pytz.timezone("Australia/Adelaide")  # Change this to your local timezone
                return local_tz.localize(date).astimezone(pytz.UTC)
            else:
                # For date-only format, set time to midnight UTC
                date = datetime.strptime(date_str, fmt)
                return date.replace(tzinfo=pytz.UTC)
        except ValueError:
            continue
    
    raise ValueError(f"Invalid date format: {date_str}")

def process_files(input_dir, output_dir):
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    for filename in os.listdir(input_dir):
        if filename.endswith('.md'):
            input_path = os.path.join(input_dir, filename)
            
            with open(input_path, 'r', encoding='utf-8') as file:
                content = file.read()
            
            try:
                date = parse_date(content)
            except ValueError as e:
                print(f"Error processing {filename}: {str(e)}")
                sys.exit(1)
            
            year_month_day = date.strftime("%Y-%m-%d")
            output_subdir = os.path.join(output_dir, 'posts', year_month_day)
            
            if not os.path.exists(output_subdir):
                os.makedirs(output_subdir)
            
            output_path = os.path.join(output_subdir, filename)
            shutil.copy2(input_path, output_path)
            print(f"Copied {filename} to {output_path}")

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python post_importer.py <input_directory> <output_directory>")
        sys.exit(1)

    input_dir = sys.argv[1]
    output_dir = sys.argv[2]

    if not os.path.isdir(input_dir):
        print(f"Error: Input directory '{input_dir}' does not exist.")
        sys.exit(1)

    try:
        process_files(input_dir, output_dir)
        print("Processing completed successfully.")
    except Exception as e:
        print(f"An error occurred: {str(e)}")
        sys.exit(1)
