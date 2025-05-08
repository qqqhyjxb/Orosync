for port in 40000 40001 40002 40003 40004; do
    lsof -t -i :$port | xargs kill -9

done

echo "kill finish。"