for d in `ls boxes/`; do
    (cd "boxes/$d"; vagrant destroy -f)
done
rm -rf boxes
rm config.json
rm join_command
