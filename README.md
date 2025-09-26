# PlistWatch

## About
PlistWatch monitors real-time changes to plist files on your system.
It outputs a `defaults` command to recreate that change.

## Install
```
go install  github.com/catilac/plistwatch@latest
```

## Usage
Just run:
```
plistwatch 
```

Now make some changes, such as moving the Dock and moving it back by clicking the *Position of Screen* options. 
You should see the changes being reported. 
You may also see other events being reported.

And you should see output such as:
```
defaults write "com.apple.dock" "orientation" 'left'
```

The output can also be filtered:
```
Usage of plistwatch:
  -filter domains
    	a comma-separated list of domains. Prefix names with "!" to exclude them. Supports globbing.
```

Examples:
- Hide annoying settings domains
`plistwatch -filter "!com.apple.knowledge-agent,!ContextStoreAgent"`
- Only show changes to the dock
`plistwatch -filter "com.apple.dock"`
- Hide every Apple domain
`plistwatch -filter "!com.apple.*"`
