## govendor-cleanup
govendor seems to sometime add transative dependencies refering to another
packages vendor folder.

This tool tries to clean that up by removing and fetching the dependency
directly using version rules or revision.

Install:
```
go install github.com/dorant/govendor-cleaner 
rehash
```

Usage:
Following command will go through the vendor.json and tries to solve paths pointing to other packages vendor folder.

```
cd <pkg>
govendor-cleaner vendor/vendor.json
```

Warning! No backup is created, make sure you have your starting files pushed to git

