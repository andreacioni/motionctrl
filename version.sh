#!/bin/bash

echo "Creating new version"

FILE=version/version.go
CURRENT_VERSION=$(cat $FILE | grep Number | awk '{print $3}' | cut -d '"' -f 2)

echo "Current version is: $CURRENT_VERSION. Enter new version:"

read NEW_VERSION

echo "New version is: $NEW_VERSION"

sed -i "s/$CURRENT_VERSION/$NEW_VERSION/g" $FILE
git add $FILE

git commit -m "Releasing v$NEW_VERSION"

git push

git tag -a v$NEW_VERSION -m "Release v$NEW_VERSION"

git push origin v$NEW_VERSION