import os
import re
import sys

def convert_markdown_to_html_video(directory_path):
    # Regular expression to match the Markdown image syntax for videos
    pattern = re.compile(r'!\[.*?\]\((.*?)\)')
    
    # Iterate over all files in the directory
    for filename in os.listdir(directory_path):
        # Check if the file is a Markdown file
        if filename.endswith('.md'):
            file_path = os.path.join(directory_path, filename)
            # Open the Markdown file and read its content
            with open(file_path, 'r') as file:
                content = file.read()
            
            # Replace the Markdown image syntax with the HTML video tag
            modified_content = pattern.sub(lambda match: f'<video controls width="320" height="240"><source src="{match.group(1)}" type="video/mp4">Your browser does not support the video tag.</video>', content)
            
            # Write the modified content back to the file
            with open(file_path, 'w') as file:
                file.write(modified_content)
            print(f"Processed: {file_path}")

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python markdown_to_html_video.py <path_to_markdown_files>")
        sys.exit(1)
    
    directory_path = sys.argv[1]
    convert_markdown_to_html_video(directory_path)
