#!/bin/bash

# merge.sh /root/yamaHub/yamaHub-repositories/interpidtjuniversity/miniselfop.git master xxx_dev
# git clone /root/yamaHub/yamaHub-repositories/interpidtjuniversity/miniselfop.git
# git checkout master
# git merge xxx_dev
# git push /root/yamaHub/yamaHub-repositories/interpidtjuniversity/miniselfop.git master

git clone $1
cd $4
git checkout $3
git checkout $2
git merge $3
git push $1 $2
