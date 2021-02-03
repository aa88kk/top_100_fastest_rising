#!/bin/sh

# copy from 'https://github.com/EvanLi/Github-Ranking/blob/master/auto_run.sh'

echo -e "\n----------Run Time:----------"
date
cd /root/top_100_fastest_rising
git pull

/root/top_100_fastest_rising/top_100_fastest_rising

git add .
today=`date +"%Y-%m-%d"`
git commit -m "updated at  $today."
git push 
