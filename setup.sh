echo -e "Carrpigeo initialization..."

# 1. Directories
mkdir -p storage

# 2. Config file
if [ ! -f config.yaml ]; then
    echo -e "Waiting for config file..."
    if [ -f config.example.yaml ]; then
        cp config.example.yaml config.yaml
    else
        echo -e "There is no example file. Check GitHub repositrory: https://github.com/dmi3midd/carrpigeo"
    fi
fi

echo -e "Initialization is completed."
