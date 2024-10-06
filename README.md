# Notes
## Git
In .git/hooks/pre-commit, add the following to change tabs to spaces:
```
#!/bin/sh
python3 go-space-format.py

SCRIPT_EXIT_STATUS=$?

if [ $SCRIPT_EXIT_STATUS -ne 0 ]; then
    echo "Python script failed, aborting commit."
    exit 1
fi

git add .

exit 0
```

Also don't forget to run `chmod +x .git/hooks/pre-commit` to make the script executable.
