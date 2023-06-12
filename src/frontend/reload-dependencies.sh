#! /bin/bash

rm -r node_modules 
rm package-lock.json
npm install 
quasar dev
