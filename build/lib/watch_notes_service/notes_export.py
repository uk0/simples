import glob
import os
from pathlib import Path

from sqlalchemy.util import md5_hex

from watch_notes_service import helpers
from watch_notes_service import notescript
from watch_notes_service.note import Note

def main():
    HOME_BASE = "/Users/firshme/Desktop/work/simples/pre/"
    staged = glob.glob(HOME_BASE + '*.md')
    for s in staged:
        os.remove(s)
    print("clear old md files")
    staged = glob.glob(HOME_BASE + '.attachments/*')
    for s in staged:
        os.remove(s)
    print("clear old attachments files")

    get_notes_script = notescript.get_notes_script
    # folder_name
    return_code, stdout, stderr = helpers.run_applescript(notescript.get_notes_script, "blog")
    local_notes = []
    if return_code != 0:
        print("AppleScript failed with return code {}".format(return_code))
        os.waitstatus_to_exitcode(1)

    staging_folder_path = stdout.strip()
    for filename in os.listdir(staging_folder_path):
        f_name, f_ext = os.path.splitext(filename)
        if not f_ext == ".staged":
            continue
        staged_file = os.path.join(staging_folder_path, filename)
        with open(staged_file) as fp:
            staged_content = fp.read()
            local_notes.append(
                Note.create_from_local(staged_content, Path(HOME_BASE)))
    print("Found {} local notes".format(len(local_notes)))

    for note in local_notes:
        with open(f"{HOME_BASE}{md5_hex(note.name)}.md","w+") as wr:
            wr.write(note.body_markdown)
            print("Note export : {}".format(note.name))
    staged = glob.glob(staging_folder_path + '/*.staged')
    for s in staged:
        os.remove(s)

    print("sync notes to server")
    os.system("bash /Users/firshme/Desktop/work/simples/watch_notes_service/sync_note.sh")
