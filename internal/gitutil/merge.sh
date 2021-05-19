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
if [ "" = "$4" ] ;then
  git merge $3
else
  git merge $3 -m "$4"
fi
git push $1 $2
