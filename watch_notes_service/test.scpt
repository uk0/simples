on run argv
set folder_name to item 1 of argv
tell application "Finder"
    set save_location to (POSIX path of (path to temporary items folder) as text) & "taskbridge/notesync/" & folder_name
    do shell script "mkdir -p " & quoted form of save_location
end tell
tell application "Notes"
    set myFolder to first folder whose name = folder_name
    set myNotes to notes of myFolder
    repeat with theNote in myNotes
        set nId to id of theNote
        set nName to name of theNote
        set nBody to body of theNote
        set nCreation to creation date of theNote
        set nModified to modification date of theNote
        set attachmentList to "~~START_ATTACHMENTS~~\n"
        repeat with theAttachment in attachments of theNote
          set attachmentList to attachmentList & name of theAttachment & "~~" & url of theAttachment & "\n"
        end repeat
        set attachmentList to attachmentList & "~~END_ATTACHMENTS~~"
        set stagedContent to nId & "~~" & nName & "~~" & nCreation & "~~" & nModified
        set stagedContent to stagedContent & "\n" & attachmentList & "\n" & nBody
        tell application "Finder"
            set aPath to "taskbridge:notesync:" & folder_name & ":" & nName & ".staged"
            set accessRef to (open for access file ((path to temporary items folder as text) & aPath) with write permission)
            try
                set eof accessRef to 0
                write stagedContent to accessRef as «class utf8»
                close access accessRef
            on error errMsg
                close access accessRef
                log errMsg
            end try
        end tell
    end repeat
end tell
return save_location
end run