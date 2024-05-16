#!/bin/bash

PROGRAM_NAME="serverstatus"
ORIGINAL_PATH="$(pwd)/$PROGRAM_NAME"
APPLICATIONS_PATH="/Applications/$PROGRAM_NAME"

# Move the program to the Applications folder
mv $ORIGINAL_PATH $APPLICATION_PATH

cat > ~/Library/LaunchAgents/$PROGRAM_NAME.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>$PROGRAM_NAME</string>
    <key>ProgramArguments</key>
    <array>
        <string>$APPLICATIONS_PATH</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
EOF

launchctl load ~/Library/LaunchAgents/$PROGRAM_NAME.plist