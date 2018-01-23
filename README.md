## govendor-cleanup
govendor seems to sometime add transative dependencies refering to another
packages vendor folder.

This tool tries to clean that up by removing and fetching the dependency
directly using version rules or revision.

### Install
```
go install github.com/dorant/govendor-cleaner
rehash
```

### Usage
First make sure all current dependencies are needed
```
govendor list +u
govendor remove +u
```

Make sure all current dependencies are vendored in:
```
govendor add +e
```
Now you might have dependencies in vendor/vendor.json that points to other packages vendor folder.

Following command goes through the vendor.json and tries to solve theses paths to make them point directly to original package:

```
cd <pkg>
govendor-cleaner vendor/vendor.json
```

Warning! No backup is created, make sure you have your starting files pushed to git

